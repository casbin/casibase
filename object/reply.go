// Copyright 2020 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"sync"
	"time"

	"github.com/casbin/casbin-forum/authz"
)

type Reply struct {
	Id          int    `xorm:"int notnull pk autoincr" json:"id"`
	Author      string `xorm:"varchar(100)" json:"author"`
	TopicId     int    `xorm:"int" json:"topicId"`
	CreatedTime string `xorm:"varchar(40)" json:"createdTime"`
	Deleted     bool   `xorm:"bool" json:"-"`
	ThanksNum   int    `xorm:"int" json:"thanksNum"`

	Content string `xorm:"mediumtext" json:"content"`
}

// GetReplyCount returns all replies num so far, both deleted and not deleted.
func GetReplyCount() int {
	count, err := adapter.engine.Count(&Reply{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

// GetReplies returns more information about reply of a topic.
func GetReplies(topicId int, memberId string) []*ReplyWithAvatar {
	replies := []*Reply{}
	err := adapter.engine.Asc("created_time").And("deleted = ?", 0).Find(&replies, &Reply{TopicId: topicId})
	if err != nil {
		panic(err)
	}

	res := []*ReplyWithAvatar{}
	for _, v := range replies {
		temp := ReplyWithAvatar{
			Reply:        *v,
			Avatar:       GetMemberAvatar(v.Author),
			ThanksStatus: GetThanksStatus(memberId, v.Id, 5),
			Deletable:    ReplyDeletable(v.CreatedTime, memberId, v.Author),
			Editable:     GetReplyEditableStatus(memberId, v.Author, v.CreatedTime),
		}
		res = append(res, &temp)
	}

	return res
}

// GetReply returns a single reply.
func GetReply(id int) *Reply {
	reply := Reply{Id: id}
	existed, err := adapter.engine.Get(&reply)
	if err != nil {
		panic(err)
	}

	if existed {
		return &reply
	}
	return nil
}

// GetReplyWithDetails returns more information about reply, including avatar, thanks status, deletable and editable.
func GetReplyWithDetails(memberId string, id int) *ReplyWithAvatar {
	reply := Reply{Id: id}
	existed, err := adapter.engine.Get(&reply)
	if err != nil {
		panic(err)
	}

	if existed {
		res := ReplyWithAvatar{
			Reply:        reply,
			Avatar:       GetMemberAvatar(reply.Author),
			ThanksStatus: GetThanksStatus(memberId, reply.Id, 5),
			Deletable:    ReplyDeletable(reply.CreatedTime, memberId, reply.Author),
			Editable:     GetReplyEditableStatus(memberId, reply.Author, reply.CreatedTime),
		}
		return &res
	}
	return nil
}

/*
func GetReplyId() int {
	reply := new(Reply)
	_, err := adapter.engine.Desc("created_time").Omit("content").Limit(1).Get(reply)
	if err != nil {
		panic(err)
	}

	res := util.ParseInt(reply.Id) + 1

	return res
}
*/

// UpdateReply updates reply's all field.
func UpdateReply(id int, reply *Reply) bool {
	if GetReply(id) == nil {
		return false
	}

	_, err := adapter.engine.Id(id).AllCols().Update(reply)
	if err != nil {
		panic(err)
	}

	//return affected != 0
	return true
}

// UpdateReplyWithLimitCols updates reply's not null field.
func UpdateReplyWithLimitCols(id int, reply *Reply) bool {
	if GetReply(id) == nil {
		return false
	}

	_, err := adapter.engine.Id(id).Update(reply)
	if err != nil {
		panic(err)
	}

	//return affected != 0
	return true
}

// AddReply returns add reply result and reply id.
func AddReply(reply *Reply) (bool, int) {
	//reply.Content = strings.ReplaceAll(reply.Content, "\n", "<br/>")

	affected, err := adapter.engine.Insert(reply)
	if err != nil {
		panic(err)
	}

	return affected != 0, reply.Id
}

/*
func DeleteReply(id string) bool {
	affected, err := adapter.engine.Id(id).Delete(&Reply{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}
*/

// DeleteReply soft delete reply.
func DeleteReply(id int) bool {
	reply := new(Reply)
	reply.Deleted = true
	affected, err := adapter.engine.Id(id).Cols("deleted").Update(reply)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

// GetLatestReplies returns member's latest replies.
func GetLatestReplies(author string, limit int, offset int) []LatestReply {
	replys := []*Reply{}
	err := adapter.engine.Where("author = ?", author).And("deleted = ?", 0).Limit(limit, offset).Find(&replys)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	var lock sync.Mutex
	errChan := make(chan error, 10)
	var result []LatestReply
	for _, v := range replys {
		wg.Add(1)
		v := v
		go func() {
			defer wg.Done()
			topic := Topic{Id: v.TopicId}
			existed, err := adapter.engine.Select("id, author, node_id, node_name, title, author").Get(&topic)
			if err != nil {
				errChan <- err
				return
			}

			if existed {
				var temp = LatestReply{
					TopicId:      topic.Id,
					NodeId:       topic.NodeId,
					NodeName:     topic.NodeName,
					Author:       topic.Author,
					ReplyContent: v.Content,
					TopicTitle:   topic.Title,
					ReplyTime:    v.CreatedTime,
				}
				lock.Lock()
				result = append(result, temp)
				lock.Unlock()
			}
		}()
	}
	wg.Wait()
	close(errChan)
	if len(errChan) != 0 {
		for v := range errChan {
			panic(v)
		}
	}
	return result
}

// GetRepliesNum returns member's all replies num.
func GetRepliesNum(memberId string) int {
	var total int64
	var err error

	reply := new(Reply)
	total, err = adapter.engine.Where("author = ?", memberId).And("deleted = ?", 0).Count(reply)
	if err != nil {
		panic(err)
	}

	return int(total)
}

// GetReplyTopicTitle only returns reply's topic title.
func GetReplyTopicTitle(id int) string {
	topic := Topic{Id: id}
	existed, err := adapter.engine.Cols("title").Get(&topic)
	if err != nil {
		panic(err)
	}

	if existed {
		return topic.Title
	}
	return ""
}

// GetReplyAuthor only returns reply's topic author.
func GetReplyAuthor(id int) string {
	reply := Reply{Id: id}
	existed, err := adapter.engine.Cols("author").Get(&reply)
	if err != nil {
		panic(err)
	}

	if existed {
		return reply.Author
	}
	return ""
}

// AddReplyThanksNum updates reply's thanks num.
func AddReplyThanksNum(id int) bool {
	reply := GetReply(id)
	if reply == nil {
		return false
	}

	reply.ThanksNum++
	affected, err := adapter.engine.Id(id).Cols("thanks_num").Update(reply)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

// ReplyDeletable checks whether the reply can be deleted.
func ReplyDeletable(date, memberId, author string) bool {
	if authz.IsRootMod(memberId) {
		return true
	}

	if memberId != author {
		return false
	}

	t, err := time.Parse("2006-01-02T15:04:05+08:00", date)
	if err != nil {
		return false
	}
	h, _ := time.ParseDuration("-1h")
	t = t.Add(8 * h)

	now := time.Now()
	if now.Sub(t).Minutes() > ReplyDeletableTime {
		return false
	}

	return true
}

// GetReplyEditableStatus checks whether the reply can be edited.
func GetReplyEditableStatus(member, author, createdTime string) bool {
	if authz.IsRootMod(member) {
		return true
	}
	if member != author {
		return false
	}

	t, err := time.Parse("2006-01-02T15:04:05+08:00", createdTime)
	if err != nil {
		return false
	}
	h, _ := time.ParseDuration("-1h")
	t = t.Add(8 * h)

	now := time.Now()
	if now.Sub(t).Minutes() > ReplyEditableTime {
		return false
	}

	return true
}

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
	"time"

	"github.com/casbin/casbin-forum/authz"
	"github.com/casbin/casbin-forum/util"
)

type Topic struct {
	Id              int      `xorm:"int notnull pk autoincr" json:"id"`
	Author          string   `xorm:"varchar(100)" json:"author"`
	NodeId          string   `xorm:"varchar(100)" json:"nodeId"`
	NodeName        string   `xorm:"varchar(100)" json:"nodeName"`
	Title           string   `xorm:"varchar(100)" json:"title"`
	CreatedTime     string   `xorm:"varchar(40)" json:"createdTime"`
	Tags            []string `xorm:"varchar(200)" json:"tags"`
	LastReplyUser   string   `xorm:"varchar(100)" json:"lastReplyUser"`
	LastReplyTime   string   `xorm:"varchar(40)" json:"lastReplyTime"`
	ReplyCount      int      `json:"replyCount"`
	UpCount         int      `json:"upCount"`
	HitCount        int      `json:"hitCount"`
	Hot             int      `json:"hot"`
	FavoriteCount   int      `json:"favoriteCount"`
	HomePageTopTime string   `xorm:"varchar(40)" json:"homePageTopTime"`
	TabTopTime      string   `xorm:"varchar(40)" json:"tabTopTime"`
	NodeTopTime     string   `xorm:"varchar(40)" json:"nodeTopTime"`
	Deleted         bool     `xorm:"bool" json:"-"`

	Content string `xorm:"mediumtext" json:"content"`
}

func GetTopicCount() int {
	count, err := adapter.engine.Count(&Topic{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetCreatedTopicsNum(memberId string) int {
	topic := new(Topic)
	total, err := adapter.engine.Where("author = ?", memberId).And("deleted = ?", 0).Count(topic)
	if err != nil {
		panic(err)
	}

	return int(total)
}

func GetTopics(limit int, offset int) []*TopicWithAvatar {
	topics := []*Topic{}
	err := adapter.engine.Desc("home_page_top_time").Desc("last_reply_time").Desc("created_time").And("deleted = ?", 0).Omit("content").Limit(limit, offset).Find(&topics)
	if err != nil {
		panic(err)
	}

	res := []*TopicWithAvatar{}
	for _, v := range topics {
		temp := TopicWithAvatar{
			Topic:  *v,
			Avatar: GetMemberAvatar(v.Author),
		}
		res = append(res, &temp)
	}

	return res
}

func GetTopicWithAvatar(id int, memberId string) *TopicWithAvatar {
	topic := Topic{Id: id}
	existed, err := adapter.engine.Get(&topic)
	if err != nil {
		panic(err)
	}

	res := TopicWithAvatar{
		Topic:        topic,
		Avatar:       GetMemberAvatar(topic.Author),
		ThanksStatus: GetThanksStatus(memberId, id, 4),
		Editable:     GetTopicEditableStatus(memberId, topic.Author, topic.NodeId, topic.CreatedTime),
	}

	if existed {
		return &res
	} else {
		return nil
	}
}

func GetTopic(id int) *Topic {
	topic := Topic{Id: id}
	existed, err := adapter.engine.Get(&topic)
	if err != nil {
		panic(err)
	}

	if existed {
		return &topic
	} else {
		return nil
	}
}

func GetTopicTitle(id int) string {
	topic := Topic{Id: id}
	existed, err := adapter.engine.Cols("title").Get(&topic)
	if err != nil {
		panic(err)
	}

	if existed {
		return topic.Title
	} else {
		return ""
	}
}

func GetTopicAuthor(id int) string {
	topic := Topic{Id: id}
	existed, err := adapter.engine.Cols("author").Get(&topic)
	if err != nil {
		panic(err)
	}

	if existed {
		return topic.Author
	} else {
		return ""
	}
}

func GetTopicNodeId(id int) string {
	topic := Topic{Id: id}
	existed, err := adapter.engine.Cols("node_id").Get(&topic)
	if err != nil {
		panic(err)
	}

	if existed {
		return topic.NodeId
	} else {
		return ""
	}
}

func GetTopicsWithNode(nodeId string, limit int, offset int) []*NodeTopic {
	topics := []*Topic{}
	err := adapter.engine.Desc("node_top_time").Desc("last_reply_time").Desc("created_time").Where("node_id = ?", nodeId).And("deleted = ?", 0).Limit(limit, offset).Find(&topics)
	if err != nil {
		panic(err)
	}

	res := []*NodeTopic{}
	for _, v := range topics {
		temp := NodeTopic{
			Topic:         *v,
			Avatar:        GetMemberAvatar(v.Author),
			ContentLength: len(v.Content),
		}
		temp.Content = ""
		res = append(res, &temp)
	}

	return res
}

func UpdateTopic(id int, topic *Topic) bool {
	if GetTopic(id) == nil {
		return false
	}

	_, err := adapter.engine.Id(id).AllCols().Update(topic)
	if err != nil {
		panic(err)
	}

	//return affected != 0
	return true
}

func UpdateTopicWithLimitCols(id int, topic *Topic) bool {
	if GetTopic(id) == nil {
		return false
	}

	_, err := adapter.engine.Id(id).Update(topic)
	if err != nil {
		panic(err)
	}

	//return affected != 0
	return true
}

// AddTopic return add topic result and topic id
func AddTopic(topic *Topic) (bool, int) {
	affected, err := adapter.engine.Insert(topic)
	if err != nil {
		panic(err)
	}

	return affected != 0, topic.Id
}

/*
func DeleteTopic(id string) bool {
	affected, err := adapter.engine.Id(id).Delete(&Topic{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}
*/

func DeleteTopic(id int) bool {
	topic := new(Topic)
	topic.Deleted = true
	affected, err := adapter.engine.Id(id).Cols("deleted").Update(topic)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

/*
func GetTopicId() int {
	topic := new(Topic)
	_, err := adapter.engine.Desc("created_time").Omit("content").Limit(1).Get(topic)
	if err != nil {
		panic(err)
	}

	res := util.ParseInt(topic.Id) + 1

	return res
}
*/

func GetAllCreatedTopics(author string, tab string, limit int, offset int) []*Topic {
	topics := []*Topic{}
	err := adapter.engine.Desc("created_time").Where("author = ?", author).And("deleted = ?", 0).Omit("content").Limit(limit, offset).Find(&topics)
	if err != nil {
		panic(err)
	}

	return topics
}

func AddTopicHitCount(topicId int) bool {
	topic := GetTopic(topicId)
	if topic == nil {
		return false
	}

	topic.HitCount++
	affected, err := adapter.engine.Id(topicId).Cols("hit_count").Update(topic)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func ChangeTopicFavoriteCount(topicId int, num int) bool {
	topic := GetTopic(topicId)
	if topic == nil {
		return false
	}

	topic.FavoriteCount += num
	affected, err := adapter.engine.Id(topicId).Cols("favorite_count").Update(topic)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func ChangeTopicReplyCount(topicId int, num int) bool {
	topic := GetTopic(topicId)
	if topic == nil {
		return false
	}

	topic.ReplyCount += num
	affected, err := adapter.engine.Id(topicId).Cols("reply_count").Update(topic)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func ChangeTopicLastReplyUser(topicId int, memberId string, updateTime bool) bool {
	topic := GetTopic(topicId)
	if topic == nil {
		return false
	}

	topic.LastReplyUser = memberId
	if updateTime {
		topic.LastReplyTime = util.GetCurrentTime()
	}
	if len(memberId) == 0 {
		topic.LastReplyTime = ""
	}
	affected, err := adapter.engine.Id(topicId).Cols("last_reply_user, last_reply_time").Update(topic)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func GetTopicsWithTab(tab string, limit, offset int) []*TopicWithAvatar {
	topics := []*Topic{}
	res := []*TopicWithAvatar{}

	if tab == "all" {
		res = GetTopics(limit, offset)
	} else {
		err := adapter.engine.Table("topic").Join("INNER", "node", "topic.node_id = node.id").Where("node.tab_id = ?", tab).Where("deleted = ?", 0).Desc("tab_top_time").Desc("topic.last_reply_time").Omit("content").Limit(limit, offset).Find(&topics)
		if err != nil {
			panic(err)
		}
		for _, v := range topics {
			temp := TopicWithAvatar{
				Topic:  *v,
				Avatar: GetMemberAvatar(v.Author),
			}
			res = append(res, &temp)
		}
	}

	return res
}

func UpdateTopicHotInfo(topicId string, hot int) bool {
	topic := new(Topic)

	topic.Hot = hot
	affected, err := adapter.engine.Id(topicId).Cols("hot").Update(topic)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func GetHotTopic(limit int) []*TopicWithAvatar {
	topics := []*Topic{}
	err := adapter.engine.Desc("hot").And("deleted = ? ", 0).Limit(limit).Find(&topics)
	if err != nil {
		panic(err)
	}

	res := []*TopicWithAvatar{}
	for _, v := range topics {
		temp := TopicWithAvatar{
			Topic:  *v,
			Avatar: GetMemberAvatar(v.Author),
		}
		res = append(res, &temp)
	}

	return res
}

func GetTopicEditableStatus(member, author, nodeId, createdTime string) bool {
	if authz.IsRootMod(member) || authz.IsNodeMod(member, nodeId) {
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
	if now.Sub(t).Minutes() > TopicEditableTime {
		return false
	}

	return true
}

// ChangeTopicTopExpiredTime changes topic's top expired time.
// topType: tab, node or homePage.
func ChangeTopicTopExpiredTime(id int, date, topType string) bool {
	topic := GetTopic(id)
	if topic == nil {
		return false
	}

	switch topType {
	case "tab":
		topic.TabTopTime = date
	case "node":
		topic.NodeTopTime = date
	case "homePage":
		topic.HomePageTopTime = date
	}

	affected, err := adapter.engine.Id(id).Cols("tab_top_time, node_top_time, home_page_top_time").Update(topic)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

// ExpireTopTopic searches and expires expired top topic.
func ExpireTopTopic() int {
	topics := []*Topic{}
	err := adapter.engine.Where("tab_top_time != ?", "").Or("node_top_time != ?", "").Or("home_page_top_time != ?", "").Cols("id, tab_top_time, node_top_time, home_page_top_time").Find(&topics)
	if err != nil {
		panic(err)
	}

	var num int
	date := util.GetCurrentTime()
	for _, v := range topics {
		if v.TabTopTime <= date {
			res := ChangeTopicTopExpiredTime(v.Id, "", "tab")
			if res {
				num++
			}
		}
		if v.NodeTopTime <= date {
			res := ChangeTopicTopExpiredTime(v.Id, "", "node")
			if res {
				num++
			}
		}
		if v.HomePageTopTime <= date {
			res := ChangeTopicTopExpiredTime(v.Id, "", "homePage")
			if res {
				num++
			}
		}
	}

	return num
}

package object

import (
	"fmt"
	"strings"
	"time"

	"github.com/casbin/casnode/service"
	"github.com/casbin/casnode/util"
	crawler "github.com/casbin/google-groups-crawler"
	"github.com/gomarkdown/markdown"
)

func (n Node) AddTopicToMailingList(title, content, author string) {
	if len(n.MailingList) == 0 {
		return
	}
	content = string(markdown.ToHTML([]byte(content), nil, nil))
	_ = service.SendEmail(title, content, n.MailingList, author)
}

func (n Node) SyncFromGoogleGroup() {
	if !strings.Contains(n.MailingList, "@googlegroups.com") {
		return
	}
	topicTitles := n.GetAllTopicTitlesOfNode()

	isInTopicList := func(topicTitle string) bool {
		for _, title := range topicTitles {
			if title == topicTitle {
				return true
			}
		}
		return false
	}

	group := crawler.NewGoogleGroup(n.MailingList, n.GoogleGroupCookie)
	conversations := group.GetAllConversations(*HttpClient)
	for _, conv := range conversations {
		messages := conv.GetAllMessages(*HttpClient, true)
		if len(messages) < 1 {
			fmt.Printf("Google Groups Crawler: Getting messages from Google Group: %s for node: %s failed, please check your cookie.\n", group.GroupName, n.Id)
			break
		}
		var newTopic Topic
		AuthorMember := AddMemberByNameAndEmailIfNotExist(messages[0].Author, messages[0].AuthorEmail)
		if AuthorMember == nil {
			continue
		}
		if !isInTopicList(conv.Title) {
			newTopic = Topic{
				Author:        AuthorMember.Id,
				NodeId:        n.Id,
				NodeName:      n.Name,
				Title:         conv.Title,
				Content:       FilterUnsafeHTML(messages[0].Content),
				CreatedTime:   util.GetTimeFromTimestamp(int64(conv.Time)),
				LastReplyTime: util.GetTimeFromTimestamp(int64(conv.Time)),
				EditorType:    "richtext",
			}
			AddTopic(&newTopic)
		} else {
			var topics []Topic
			err := adapter.Engine.Where("title = ? and deleted = 0", conv.Title).Find(&topics)
			if err != nil {
				panic(err)
			}
			if len(topics) == 0 {
				continue
			}
			for _, t := range topics {
				if conv.Title == t.Title {
					newTopic = t
					break
				}
			}
		}

		replies := newTopic.GetAllRepliesOfTopic()
		isInReplies := func(replyStr string) bool {
			for _, c := range replies {
				if c == replyStr {
					return true
				}
			}
			return false
		}

		for _, msg := range messages[1:] {
			msg.Content = FilterUnsafeHTML(msg.Content)
			AuthorMember = AddMemberByNameAndEmailIfNotExist(msg.Author, msg.AuthorEmail)
			if AuthorMember == nil {
				continue
			}
			if isInReplies(msg.Content) {
				continue
			}
			newReply := Reply{
				Author:      AuthorMember.Id,
				TopicId:     newTopic.Id,
				EditorType:  "richtext",
				Content:     msg.Content,
				CreatedTime: util.GetTimeFromTimestamp(int64(msg.Time)),
			}
			AddReply(&newReply)
			newTopic.LastReplyTime = util.GetTimeFromTimestamp(int64(msg.Time))
			newTopic.LastReplyUser = AuthorMember.Id
		}
		UpdateTopic(newTopic.Id, &newTopic)
	}
}

func AutoSyncGoogleGroup() {
	if AutoSyncPeriodSecond < 30 {
		return
	}
	for {
		time.Sleep(time.Duration(AutoSyncPeriodSecond) * time.Second)
		SyncAllNodeFromGoogleGroup()
	}
}

func SyncAllNodeFromGoogleGroup() {
	if AutoSyncPeriodSecond < 30 {
		return
	}
	fmt.Println("Sync from google group started...")
	var nodes []Node
	err := adapter.Engine.Find(&nodes)
	if err != nil {
		panic(err)
	}
	for _, node := range nodes {
		node.SyncFromGoogleGroup()
	}
}

func (r Reply) AddReplyToMailingList() {
	targetTopic := GetTopic(r.TopicId)
	targetNode := GetNode(targetTopic.NodeId)
	if len(targetNode.MailingList) == 0 {
		return
	}
	if r.EditorType == "markdown" {
		r.Content = string(markdown.ToHTML([]byte(r.Content), nil, nil))
	}
	mailTitle := fmt.Sprintf("Re: %s", targetTopic.Title)
	_ = service.SendEmail(mailTitle, r.Content, targetNode.MailingList, r.Author)
}

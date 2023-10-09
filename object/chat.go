// Copyright 2023 The casbin Authors. All Rights Reserved.
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
	"fmt"

	"github.com/casibase/casibase/util"
	"xorm.io/core"
)

type Chat struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	UpdatedTime string `xorm:"varchar(100)" json:"updatedTime"`

	// Organization string   `xorm:"varchar(100)" json:"organization"`
	DisplayName  string   `xorm:"varchar(100)" json:"displayName"`
	Category     string   `xorm:"varchar(100)" json:"category"`
	Type         string   `xorm:"varchar(100)" json:"type"`
	User1        string   `xorm:"varchar(100)" json:"user1"`
	User2        string   `xorm:"varchar(100)" json:"user2"`
	Users        []string `xorm:"varchar(100)" json:"users"`
	MessageCount int      `json:"messageCount"`
}

func GetGlobalChats() ([]*Chat, error) {
	chats := []*Chat{}
	err := adapter.engine.Asc("owner").Desc("created_time").Find(&chats)
	if err != nil {
		return chats, err
	}

	return chats, nil
}

func GetChats(owner string) ([]*Chat, error) {
	chats := []*Chat{}
	err := adapter.engine.Desc("created_time").Find(&chats, &Chat{Owner: owner})
	if err != nil {
		return chats, err
	}

	return chats, nil
}

func getChat(owner, name string) (*Chat, error) {
	chat := Chat{Owner: "admin", Name: name}
	existed, err := adapter.engine.Get(&chat)
	if err != nil {
		return nil, err
	}

	if existed {
		return &chat, nil
	} else {
		return nil, nil
	}
}

func GetChat(id string) (*Chat, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getChat(owner, name)
}

func UpdateChat(id string, chat *Chat) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	_, err := getChat(owner, name)
	if err != nil {
		return false, err
	}
	if chat == nil {
		return false, nil
	}

	_, err = adapter.engine.ID(core.PK{owner, name}).AllCols().Update(chat)
	if err != nil {
		return false, err
	}

	// return affected != 0
	return true, nil
}

func AddChat(chat *Chat) (bool, error) {
	//if chat.Type == "AI" && chat.User2 == "" {
	//	provider, err := GetDefaultModelProvider()
	//	if err != nil {
	//		return false, err
	//	}
	//
	//	if provider != nil {
	//		chat.User2 = provider.Name
	//	}
	//}

	affected, err := adapter.engine.Insert(chat)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func DeleteChat(chat *Chat) (bool, error) {
	affected, err := adapter.engine.ID(core.PK{chat.Owner, chat.Name}).Delete(&Chat{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func (chat *Chat) GetId() string {
	return fmt.Sprintf("%s/%s", chat.Owner, chat.Name)
}

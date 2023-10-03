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

	"github.com/casbin/casibase/util"
	"xorm.io/core"
)

type VectorScore struct {
	Vector string  `xorm:"varchar(100)" json:"vector"`
	Score  float32 `json:"score"`
}

type Message struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	// Organization string `xorm:"varchar(100)" json:"organization"`
	Chat         string        `xorm:"varchar(100) index" json:"chat"`
	ReplyTo      string        `xorm:"varchar(100) index" json:"replyTo"`
	Author       string        `xorm:"varchar(100)" json:"author"`
	Text         string        `xorm:"mediumtext" json:"text"`
	VectorScores []VectorScore `xorm:"mediumtext" json:"vectorScores"`
}

func GetGlobalMessages() ([]*Message, error) {
	messages := []*Message{}
	err := adapter.engine.Asc("owner").Desc("created_time").Find(&messages)
	if err != nil {
		return messages, err
	}

	return messages, nil
}

func GetChatMessages(chat string) ([]*Message, error) {
	messages := []*Message{}
	err := adapter.engine.Asc("created_time").Find(&messages, &Message{Chat: chat})
	if err != nil {
		return messages, err
	}

	return messages, nil
}

func GetMessages(owner string) ([]*Message, error) {
	messages := []*Message{}
	err := adapter.engine.Desc("created_time").Find(&messages, &Message{Owner: owner})
	if err != nil {
		return messages, err
	}

	return messages, nil
}

func getMessage(owner, name string) (*Message, error) {
	message := Message{Owner: owner, Name: name}
	existed, err := adapter.engine.Get(&message)
	if err != nil {
		return &message, err
	}

	if existed {
		return &message, nil
	} else {
		return nil, nil
	}
}

func GetMessage(id string) (*Message, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getMessage(owner, name)
}

func UpdateMessage(id string, message *Message) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	_, err := getMessage(owner, name)
	if err != nil {
		return false, err
	}
	if message == nil {
		return false, nil
	}

	_, err = adapter.engine.ID(core.PK{owner, name}).AllCols().Update(message)
	if err != nil {
		return false, err
	}

	return true, nil
}

func AddMessage(message *Message) (bool, error) {
	affected, err := adapter.engine.Insert(message)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func DeleteMessage(message *Message) (bool, error) {
	affected, err := adapter.engine.ID(core.PK{message.Owner, message.Name}).Delete(&Message{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func (message *Message) GetId() string {
	return fmt.Sprintf("%s/%s", message.Owner, message.Name)
}

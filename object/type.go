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

type LatestReply struct {
	TopicId      string `json:"topicId"`
	NodeId       string `json:"nodeId"`
	NodeName     string `json:"nodeName"`
	Author       string `json:"author"`
	ReplyContent string `json:"replyContent"`
	TopicTitle   string `json:"topicTitle"`
	ReplyTime    string `json:"replyTime"`
}

type NodeFavoritesRes struct {
	NodeInfo *Node `json:"nodeInfo"`
	TopicNum int   `json:"topicNum"`
}

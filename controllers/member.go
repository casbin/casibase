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

package controllers

import (
	"encoding/json"

	"github.com/casbin/casbin-forum/object"
)

func (c *APIController) GetMembers() {
	c.Data["json"] = object.GetMembers()
	c.ServeJSON()
}

func (c *APIController) GetMember() {
	id := c.Input().Get("id")

	c.Data["json"] = object.GetMember(id)
	c.ServeJSON()
}

func (c *APIController) GetMemberAvatar() {
	id := c.Input().Get("id")

	c.Data["json"] = object.GetMemberAvatar(id)
	c.ServeJSON()
}

func (c *APIController) UpdateMember() {
	id := c.Input().Get("id")

	var member object.Member
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &member)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = object.UpdateMember(id, &member)
	c.ServeJSON()
}

func (c *APIController) UpdateMemberInfo() {
	id := c.Input().Get("id")
	memberId := c.GetSessionUser()

	var tempMember object.Member
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &tempMember)
	if err != nil {
		panic(err)
	}

	var resp Response
	if memberId != id {
		resp = Response{Status: "fail", Msg: "Unauthorized."}
	} else {
		var member = object.Member{
			Company: tempMember.Company,
			Bio:     tempMember.Bio,
			Website: tempMember.Website,
		}
		res := object.UpdateMemberInfo(id, &member)
		resp = Response{Status: "ok", Msg: "success", Data: res}
	}

	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *APIController) AddMember() {
	var member object.Member
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &member)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = object.AddMember(&member)
	c.ServeJSON()
}

func (c *APIController) DeleteMember() {
	id := c.Input().Get("id")

	c.Data["json"] = object.DeleteMember(id)
	c.ServeJSON()
}

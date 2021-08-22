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
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/casbin/casnode/object"
	"github.com/casbin/casnode/service"
)

var HttpClient *http.Client

func InitHttpClient() {
	HttpClient = object.GetProxyHttpClient()
}

func UploadAvatarToOSS(avatar, memberId string) string {
	if len(avatar) == 0 {
		data := []byte(memberId)
		has := md5.Sum(data)
		memberMd5 := fmt.Sprintf("%x", has)
		avatar = "https://www.gravatar.com/avatar/" + memberMd5 + "?d=retro"
	}

	response, err := HttpClient.Get(avatar)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	avatarInfo, err := ioutil.ReadAll(response.Body)

	avatarURL := service.UploadAvatarToOSS(avatarInfo, memberId)

	return avatarURL
}

// ResponseOk ...
func (c *ApiController) ResponseOk(data ...interface{}) {
	resp := Response{Status: "ok"}
	switch len(data) {
	case 2:
		resp.Data2 = data[1]
		fallthrough
	case 1:
		resp.Data = data[0]
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

// ResponseError ...
func (c *ApiController) ResponseError(error string, data ...interface{}) {
	resp := Response{Status: "error", Msg: error}
	switch len(data) {
	case 2:
		resp.Data2 = data[1]
		fallthrough
	case 1:
		resp.Data = data[0]
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *ApiController) RequireAdmin(memberId string) {
	resp := Response{Status: "error", Msg: "Unauthorized.", Data: memberId}
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *ApiController) RequireAdminRight() bool {
	user := c.GetSessionUser()
	if user == nil || !object.CheckIsAdmin(user) {
		c.ResponseError("This operation requires admin privilege")
		return false
	}

	return true
}

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

//go:build !skipCi
// +build !skipCi

package ai

import (
	"github.com/sashabaranov/go-openai"
	"testing"

	"github.com/casbin/casibase/object"
	"github.com/casbin/casibase/proxy"
)

func TestRun(t *testing.T) {
	object.InitConfig()
	proxy.InitHttpClient()

	text, err := queryAnswer("", "hi", 5)
	if err != nil {
		panic(err)
	}

	println(text)
}

func TestToken(t *testing.T) {
	println(getTokenSize(openai.GPT3TextDavinci003, ""))
}

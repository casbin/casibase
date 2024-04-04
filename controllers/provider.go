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

package controllers

import (
	"encoding/json"

	"github.com/casibase/casibase/object"
)

// GetGlobalProviders
// @Title GetGlobalProviders
// @Tag Provider API
// @Description get global providers
// @Success 200 {array} object.Provider The Response object
// @router /get-global-providers [get]
func (c *ApiController) GetGlobalProviders() {
	providers, err := object.GetGlobalProviders()
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(object.GetMaskedProviders(providers, true))
}

// GetProviders
// @Title GetProviders
// @Tag Provider API
// @Description get providers
// @Success 200 {array} object.Provider The Response object
// @router /get-providers [get]
func (c *ApiController) GetProviders() {
	owner := "admin"

	providers, err := object.GetProviders(owner)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(object.GetMaskedProviders(providers, true))
}

// GetProvider
// @Title GetProvider
// @Tag Provider API
// @Description get provider
// @Param id query string true "The id of provider"
// @Success 200 {object} object.Provider The Response object
// @router /get-provider [get]
func (c *ApiController) GetProvider() {
	id := c.Input().Get("id")

	provider, err := object.GetProvider(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(object.GetMaskedProvider(provider, true))
}

// UpdateProvider
// @Title UpdateProvider
// @Tag Provider API
// @Description update provider
// @Param id query string true "The id (owner/name) of the provider"
// @Param body body object.Provider true "The details of the provider"
// @Success 200 {object} controllers.Response The Response object
// @router /update-provider [post]
func (c *ApiController) UpdateProvider() {
	id := c.Input().Get("id")

	var provider object.Provider
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &provider)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	success, err := object.UpdateProvider(id, &provider)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(success)
}

// AddProvider
// @Title AddProvider
// @Tag Provider API
// @Description add provider
// @Param body body object.Provider true "The details of the provider"
// @Success 200 {object} controllers.Response The Response object
// @router /add-provider [post]
func (c *ApiController) AddProvider() {
	var provider object.Provider
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &provider)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	provider.Owner = "admin"
	success, err := object.AddProvider(&provider)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(success)
}

// DeleteProvider
// @Title DeleteProvider
// @Tag Provider API
// @Description delete provider
// @Param body body object.Provider true "The details of the provider"
// @Success 200 {object} controllers.Response The Response object
// @router /delete-provider [post]
func (c *ApiController) DeleteProvider() {
	var provider object.Provider
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &provider)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	success, err := object.DeleteProvider(&provider)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(success)
}

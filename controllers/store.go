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
	"fmt"
	"net"
	"strings"

	"github.com/casibase/casibase/object"
)

func (c *ApiController) GetGlobalStores() {
	stores, err := object.GetGlobalStores()
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(stores)
}

func (c *ApiController) GetStores() {
	owner := c.Input().Get("owner")

	stores, err := object.GetStores(owner)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(stores)
}

func isIpAddress(host string) bool {
	// Attempt to split the host and port, ignoring the error
	hostWithoutPort, _, err := net.SplitHostPort(host)
	if err != nil {
		// If an error occurs, it might be because there's no port
		// In that case, use the original host string
		hostWithoutPort = host
	}

	// Attempt to parse the host as an IP address (both IPv4 and IPv6)
	ip := net.ParseIP(hostWithoutPort)
	// if host is not nil is an IP address else is not an IP address
	return ip != nil
}

func getOriginFromHost(host string) string {
	protocol := "https://"
	if !strings.Contains(host, ".") {
		// "localhost:14000"
		protocol = "http://"
	} else if isIpAddress(host) {
		// "192.168.0.10"
		protocol = "http://"
	}

	return fmt.Sprintf("%s%s", protocol, host)
}

func (c *ApiController) GetStore() {
	id := c.Input().Get("id")

	var store *object.Store
	var err error
	if id == "admin/_casibase_default_store_" {
		store, err = object.GetDefaultStore("admin")
	} else {
		store, err = object.GetStore(id)
	}
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	if store == nil {
		c.ResponseError("store is empty")
		return
	}

	host := c.Ctx.Request.Host
	origin := getOriginFromHost(host)
	err = store.Populate(origin)
	if err != nil {
		// gentle error
		c.ResponseOk(store, err.Error())
		return
	}

	c.ResponseOk(store)
}

func (c *ApiController) UpdateStore() {
	id := c.Input().Get("id")

	var store object.Store
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &store)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	success, err := object.UpdateStore(id, &store)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(success)
}

func (c *ApiController) AddStore() {
	var store object.Store
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &store)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	sucess, err := object.AddStore(&store)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(sucess)
}

func (c *ApiController) DeleteStore() {
	var store object.Store
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &store)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	sucess, err := object.DeleteStore(&store)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(sucess)
}

func (c *ApiController) RefreshStoreVectors() {
	var store object.Store
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &store)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	ok, err := object.RefreshStoreVectors(&store)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(ok)
}

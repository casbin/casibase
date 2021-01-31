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

import * as Setting from "../Setting";

export function getAccount() {
  return fetch(`${Setting.ServerUrl}/api/get-account`, {
    method: "GET",
    credentials: "include",
  }).then((res) => res.json());
}

export function signup(values) {
  return fetch(`${Setting.ServerUrl}/api/signup`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(values),
  }).then((res) => res.json());
}

export function signin(values) {
  return fetch(`${Setting.ServerUrl}/api/signin`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(values),
  }).then((res) => res.json());
}

export function signout() {
  return fetch(`${Setting.ServerUrl}/api/signout`, {
    method: "POST",
    credentials: "include",
  }).then((res) => res.json());
}

export function getStsToken() {
  return fetch(`${Setting.ServerUrl}/api/get-member-sts-token`, {
    method: "GET",
    credentials: "include",
  }).then((res) => res.json());
}

export function resetPassword(step, values) {
  return fetch(`${Setting.ServerUrl}/api/reset-password?step=${step}`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(values),
  }).then((res) => res.json());
}

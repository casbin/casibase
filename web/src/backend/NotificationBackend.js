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

export function getNotifications(limit, page) {
    return fetch(
        `${Setting.ServerUrl}/api/get-notifications?limit=${limit}&page=${page}`,
        {
            method: "GET",
            credentials: "include",
        }
    ).then((res) => res.json());
}

export function deleteNotifications(id) {
    return fetch(`${Setting.ServerUrl}/api/delete-notifications?id=${id}`, {
        method: "POST",
        credentials: "include",
    }).then((res) => res.json());
}

export function getUnreadNotificationNum() {
    return fetch(`${Setting.ServerUrl}/api/get-unread-notification-num`, {
        method: "GET",
        credentials: "include",
    }).then((res) => res.json());
}

export function updateReadStatus() {
    return fetch(`${Setting.ServerUrl}/api/update-read-status`, {
        method: "POST",
        credentials: "include",
    }).then((res) => res.json());
}

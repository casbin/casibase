// Copyright 2021 The casbin Authors. All Rights Reserved.
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
import * as ConfBackend from "./backend/ConfBackend.js";

export let AuthConfig = {
  // serverUrl: "https://door.casbin.com",
  serverUrl: "",
  clientId: "",
  appName: "",
  organizationName: "",
};

export let FrontConfig = {
  forumName: "",
  logoImage: "",
  footerLogoImage: "",
  footerLogoUrl: "",
  signinBoxStrong: "",
  signinBoxSpan: "",
  footerDeclaration: "",
  footerAdvise: "",
};

export function getFrontConf(field, callback) {
  let config;
  if (field == "visualConf") {
    config = FrontConfig;
  } else if (field == "authConf") {
    config = AuthConfig;
  }
  var storage = window.localStorage;
  for (let conf in config) {
    if (storage[conf] !== undefined) {
      config[conf] = storage[conf];
    }
  }
  ConfBackend.getFrontConfByField(field).then((res) => {
    for (let key in res) {
      if (res[key].Value !== "") {
        config[res[key].Id] = res[key].Value;
      }
      storage[res[key].Id] = config[res[key].Id];
    }
    callback();
  });
}

export const ShowEmbedButtons = false;

export const ShowGithubCorner = true;
export const GithubRepo = "https://github.com/casbin/casnode";

export const Domain = "forum.casbin.com";

export const ForceLanguage = "";
export const DefaultLanguage = "en";

// Support: richtext | markdown
export const DefaultEditorType = "markdown";

//Default search engine
//Support: baidu(www.baidu.com) | google(www.google.com) | cn-bing(cn.bing.com)
export const DefaultSearchSite = "google";

export const EnableNotificationAutoUpdate = false;

export const NotificationAutoUpdatePeriod = 10; // second

export const DefaultTopicPageReplyNum = 100;

export const ReplyMaxDepth = 10;
export const ReplyMobileMaxDepth = 3;

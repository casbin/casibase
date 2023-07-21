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

import {Tag, Tooltip, message} from "antd";
import {SyncOutlined} from "@ant-design/icons";
import {isMobile as isMobileDevice} from "react-device-detect";
import i18next from "i18next";
import Sdk from "casdoor-js-sdk";
import FileSaver from "file-saver";
import XLSX from "xlsx";
import moment from "moment/moment";
import * as StoreBackend from "./backend/StoreBackend";

export let ServerUrl = "";
export let CasdoorSdk;

export function initServerUrl() {
  const hostname = window.location.hostname;
  if (hostname === "localhost") {
    ServerUrl = `http://${hostname}:14000`;
  }
}

export function initCasdoorSdk(config) {
  CasdoorSdk = new Sdk(config);
}

function getUrlWithLanguage(url) {
  if (url.includes("?")) {
    return `${url}&language=${getLanguage()}`;
  } else {
    return `${url}?language=${getLanguage()}`;
  }
}

export function getSignupUrl() {
  return getUrlWithLanguage(CasdoorSdk.getSignupUrl());
}

export function getSigninUrl() {
  return getUrlWithLanguage(CasdoorSdk.getSigninUrl());
}

export function getUserProfileUrl(userName, account) {
  return getUrlWithLanguage(CasdoorSdk.getUserProfileUrl(userName, account));
}

export function getMyProfileUrl(account) {
  return getUrlWithLanguage(CasdoorSdk.getMyProfileUrl(account));
}

export function signin() {
  return CasdoorSdk.signin(ServerUrl);
}

export function parseJson(s) {
  if (s === "") {
    return null;
  } else {
    return JSON.parse(s);
  }
}

export function myParseInt(i) {
  const res = parseInt(i);
  return isNaN(res) ? 0 : res;
}

export function openLink(link) {
  // this.props.history.push(link);
  const w = window.open("about:blank");
  w.location.href = link;
}

export function goToLink(link) {
  window.location.href = link;
}

export function goToLinkSoft(ths, link) {
  ths.props.history.push(link);
}

export function showMessage(type, text) {
  if (type === "") {
    return;
  } else if (type === "success") {
    message.success(text);
  } else if (type === "error") {
    message.error(text);
  }
}

export function isAdminUser(account) {
  if (account === undefined || account === null) {
    return false;
  }
  return account.owner === "built-in" || account.isGlobalAdmin === true;
}

export function isLocalAdminUser(account) {
  if (account === undefined || account === null) {
    return false;
  }
  return account.isAdmin === true || isAdminUser(account);
}

export function deepCopy(obj) {
  if (obj === null) {
    showMessage("error", "deepCopy obj is null.");
  }
  return Object.assign({}, obj);
}

export function insertRow(array, row, i) {
  return [...array.slice(0, i), row, ...array.slice(i)];
}

export function addRow(array, row) {
  return [...array, row];
}

export function prependRow(array, row) {
  return [row, ...array];
}

export function deleteRow(array, i) {
  // return array = array.slice(0, i).concat(array.slice(i + 1));
  return [...array.slice(0, i), ...array.slice(i + 1)];
}

export function swapRow(array, i, j) {
  return [...array.slice(0, i), array[j], ...array.slice(i + 1, j), array[i], ...array.slice(j + 1)];
}

export function trim(str, ch) {
  if (str === undefined) {
    return undefined;
  }

  let start = 0;
  let end = str.length;

  while (start < end && str[start] === ch) {++start;}

  while (end > start && str[end - 1] === ch) {--end;}

  return (start > 0 || end < str.length) ? str.substring(start, end) : str;
}

export function isMobile() {
  // return getIsMobileView();
  return isMobileDevice;
}

export function getFormattedDate(date) {
  if (date === undefined || date === null) {
    return null;
  }

  date = date.replace("T", " ");
  date = date.replace("+08:00", " ");
  return date;
}

export function getFormattedDateShort(date) {
  return date.slice(0, 10);
}

export function getShortName(s) {
  return s.split("/").slice(-1)[0];
}

export function getShortText(s, maxLength = 35) {
  if (s.length > maxLength) {
    return `${s.slice(0, maxLength)}...`;
  } else {
    return s;
  }
}

function getRandomInt(s) {
  let hash = 0;
  if (s.length !== 0) {
    for (let i = 0; i < s.length; i++) {
      const char = s.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash;
    }
  }

  return hash;
}

export function getAvatarColor(s) {
  const colorList = ["#f56a00", "#7265e6", "#ffbf00", "#00a2ae"];
  let random = getRandomInt(s);
  if (random < 0) {
    random = -random;
  }
  return colorList[random % 4];
}

export function getLanguage() {
  return i18next.language;
}

export function setLanguage(language) {
  localStorage.setItem("language", language);
  changeMomentLanguage(language);
  i18next.changeLanguage(language);
}

export function changeLanguage(language) {
  localStorage.setItem("language", language);
  changeMomentLanguage(language);
  i18next.changeLanguage(language);
  window.location.reload(true);
}

export function changeMomentLanguage(lng) {
  return;
  // if (lng === "zh") {
  //   moment.locale("zh", {
  //     relativeTime: {
  //       future: "%s内",
  //       past: "%s前",
  //       s: "几秒",
  //       ss: "%d秒",
  //       m: "1分钟",
  //       mm: "%d分钟",
  //       h: "1小时",
  //       hh: "%d小时",
  //       d: "1天",
  //       dd: "%d天",
  //       M: "1个月",
  //       MM: "%d个月",
  //       y: "1年",
  //       yy: "%d年",
  //     },
  //   });
  // }
}

export function getTag(text, type, state) {
  let icon = null;
  let style = {};
  if (state === "Pending") {
    icon = <SyncOutlined spin />;
    style = {borderStyle: "dashed", backgroundColor: "white"};
  }

  if (type === "Read") {
    return (
      <Tooltip placement="top" title={"Read"}>
        <Tag icon={icon} style={style} color={"success"}>
          {text}
        </Tag>
      </Tooltip>
    );
  } else if (type === "Write") {
    return (
      <Tooltip placement="top" title={"Write"}>
        <Tag icon={icon} style={style} color={"processing"}>
          {text}
        </Tag>
      </Tooltip>
    );
  } else if (type === "Admin") {
    return (
      <Tooltip placement="top" title={"Admin"}>
        <Tag icon={icon} style={style} color={"error"}>
          {text}
        </Tag>
      </Tooltip>
    );
  } else {
    return null;
  }
}

export function getTags(factors, type) {
  if (!factors) {
    return [];
  }

  if (type === "factors") {
    return getFactorTag(factors);
  } else if (type === "users") {
    return getUserTag(factors);
  }
}

function getFactorTag(factors) {
  const res = [];
  factors.forEach((factor, i) => {
    if (factor.data.length !== 0) {
      res.push(
        <Tooltip placement="top" title={getShortText(JSON.stringify(factor.data), 500)}>
          <Tag color={"success"}>
            {factor.name}
          </Tag>
        </Tooltip>
      );
    } else {
      res.push(
        <Tag color={"warning"}>
          {factor.name}
        </Tag>
      );
    }
  });
  return res;
}

function getUserTag(users) {
  const res = [];
  users.forEach((user, i) => {
    if (user.length !== 0) {
      res.push(
        <Tooltip placement="top" title={getShortText(JSON.stringify(user), 500)}>
          <Tag color={"success"}>
            {user}
          </Tag>
        </Tooltip>
      );
    } else {
      res.push(
        <Tag color={"warning"}>
          {user}
        </Tag>
      );
    }
  });
  return res;
}

export function getLabelTags(labels) {
  if (!labels) {
    return [];
  }

  const res = [];
  labels.forEach((label, i) => {
    res.push(
      <Tooltip placement="top" title={getShortText(JSON.stringify(label.text), 500)}>
        <Tag color={"processing"}>
          {`${label.startTime}: ${label.text !== "" ? label.text : "(Empty)"}`}
        </Tag>
      </Tooltip>
    );
  });
  return res;
}

export function getPercentage(f) {
  if (f === undefined) {
    return 0.0;
  }

  return (100 * f).toFixed(1);
}

function s2ab(s) {
  const buf = new ArrayBuffer(s.length);
  const view = new Uint8Array(buf);
  for (let i = 0; i !== s.length; i++) {
    view[i] = s.charCodeAt(i) & 0xFF;
  }
  return buf;
}

export function sheet2blob(sheet, sheetName) {
  const workbook = {
    SheetNames: [sheetName],
    Sheets: {},
  };
  workbook.Sheets[sheetName] = sheet;
  return workbook2blob(workbook);
}

export function workbook2blob(workbook) {
  const wopts = {
    bookType: "xlsx",
    bookSST: false,
    type: "binary",
  };
  const wbout = XLSX.write(workbook, wopts);
  return new Blob([s2ab(wbout)], {type: "application/octet-stream"});
}

export function downloadXlsx(wordset) {
  const data = [];
  wordset.factors.forEach((factor, i) => {
    const row = {};

    row[0] = factor.name;
    factor.data.forEach((dataItem, i) => {
      row[i + 1] = dataItem;
    });

    data.push(row);
  });

  const sheet = XLSX.utils.json_to_sheet(data, {skipHeader: true});
  // sheet["!cols"] = [
  //   { wch: 18 },
  //   { wch: 7 },
  // ];

  try {
    const blob = sheet2blob(sheet, "factors");
    const fileName = `factors-${wordset.name}.xlsx`;
    FileSaver.saveAs(blob, fileName);
  } catch (error) {
    showMessage("error", `failed to download: ${error.message}`);
  }
}

export function downloadLabels(table) {
  const data = [];
  table.forEach((label, i) => {
    const row = {};

    row[0] = label.startTime;
    row[1] = label.endTime;
    row[2] = label.text;
    data.push(row);
  });

  const sheet = XLSX.utils.json_to_sheet(data, {skipHeader: true});
  try {
    const blob = sheet2blob(sheet, "labels");
    const fileName = `labels-${this.props.video.name}-${table.length}.xlsx`;
    FileSaver.saveAs(blob, fileName);
  } catch (error) {
    showMessage("error", `failed to download: ${error.message}`);
  }
}

export const redirectCatchJsonError = async(url) => {
  try {
    const response = await fetch(url);
    const msg = await response.json();
    if (response.ok) {
      this.props.history.push(url);
    } else {
      showMessage("error", `error in redirect: ${msg}`);
    }
  } catch (error) {
    showMessage("error", `failed to redirect: ${error.message}`);
  }
};

export function toFixed(f, n) {
  return parseFloat(f.toFixed(n));
}

export function getRandomName() {
  return Math.random().toString(36).slice(-6);
}

export function getFriendlyFileSize(size) {
  if (size < 1024) {
    return size + " B";
  }

  const i = Math.floor(Math.log(size) / Math.log(1024));
  let num = (size / Math.pow(1024, i));
  const round = Math.round(num);
  num = round < 10 ? num.toFixed(2) : round < 100 ? num.toFixed(1) : round;
  return `${num} ${"KMGTPEZY"[i - 1]}B`;
}

export function getTreeWithParents(tree) {
  const res = deepCopy(tree);
  res.children = tree.children.map((file, index) => {
    file.parent = tree;
    return getTreeWithParents(file);
  });
  return res;
}

export function getTreeWithSearch(tree, s) {
  const res = deepCopy(tree);
  res.children = tree.children.map((file, index) => {
    if (file.children.length === 0) {
      if (file.title.includes(s)) {
        return file;
      } else {
        return null;
      }
    } else {
      const tmpTree = getTreeWithSearch(file, s);
      if (tmpTree.children.length !== 0) {
        return tmpTree;
      } else {
        if (file.title.includes(s)) {
          return file;
        } else {
          return null;
        }
      }
    }
  }).filter((file, index) => {
    return file !== null;
  });
  return res;
}

export function getExtFromPath(path) {
  const filename = path.split("/").pop();
  if (filename.includes(".")) {
    return filename.split(".").pop().toLowerCase();
  } else {
    return "";
  }
}

export function getExtFromFile(file) {
  const res = file.title.split(".")[1];
  if (res === undefined) {
    return "";
  } else {
    return res;
  }
}

export function getFileCategory(file) {
  if (file.isLeaf) {
    return i18next.t("store:File");
  } else {
    return i18next.t("store:Folder");
  }
}

export function getDistinctArray(arr) {
  return [...new Set(arr)];
}

export function getCollectedTime(filename) {
  // 20220827_210300_CH~Logo.png
  const tokens = filename.split("~");
  if (tokens.length < 2) {
    return null;
  }

  const time = tokens[0].slice(0, -3);
  const m = new moment(time, "YYYYMMDD_HH:mm:ss");
  return m.format();
}

export function getSubject(filename) {
  // 20220827_210300_CH~Logo.png
  const tokens = filename.split("~");
  if (tokens.length < 2) {
    return null;
  }

  const subject = tokens[0].slice(tokens[0].length - 2);
  if (subject === "MA") {
    return i18next.t("store:Math");
  } else if (subject === "CH") {
    return i18next.t("store:Chinese");
  } else if (subject === "NU") {
    return null;
  } else {
    return subject;
  }
}

export function submitStoreEdit(storeObj) {
  const store = deepCopy(storeObj);
  store.fileTree = undefined;
  StoreBackend.updateStore(storeObj.owner, storeObj.name, store)
    .then((res) => {
      if (res.status === "ok") {
        if (res.data) {
          showMessage("success", "Successfully saved");
        } else {
          showMessage("error", "failed to save: server side failure");
        }
      } else {
        showMessage("error", `failed to save: ${res.msg}`);
      }
    })
    .catch(error => {
      showMessage("error", `failed to save: ${error}`);
    });
}

export const StaticBaseUrl = "https://cdn.casbin.org";

export const Countries = [{label: "English", key: "en", country: "US", alt: "English"},
  {label: "中文", key: "zh", country: "CN", alt: "中文"},
  {label: "Español", key: "es", country: "ES", alt: "Español"},
  {label: "Français", key: "fr", country: "FR", alt: "Français"},
  {label: "Deutsch", key: "de", country: "DE", alt: "Deutsch"},
  {label: "Indonesia", key: "id", country: "ID", alt: "Indonesia"},
  {label: "日本語", key: "ja", country: "JP", alt: "日本語"},
  {label: "한국어", key: "ko", country: "KR", alt: "한국어"},
  {label: "Русский", key: "ru", country: "RU", alt: "Русский"},
  {label: "TiếngViệt", key: "vi", country: "VN", alt: "TiếngViệt"},
  {label: "Português", key: "pt", country: "BR", alt: "Português"},
];

export function getItem(label, key, icon, children, type) {
  return {
    key,
    icon,
    children,
    label,
    type,
  };
}

export function getOption(label, value) {
  return {
    label,
    value,
  };
}

export function scrollToDiv(divId) {
  if (divId) {
    const ele = document.getElementById(divId);
    if (ele) {
      ele.scrollIntoView({behavior: "smooth"});
    }
  }
}

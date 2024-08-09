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

import React from "react";
import {Link} from "react-router-dom";
import {Button, Popconfirm, Table, Upload} from "antd";
import {UploadOutlined} from "@ant-design/icons";
import moment from "moment";
import BaseListPage from "./BaseListPage";
import * as Setting from "./Setting";
import * as VideoBackend from "./backend/VideoBackend";
import i18next from "i18next";

class VideoListPage extends BaseListPage {
  constructor(props) {
    super(props);
  }

  UNSAFE_componentWillMount() {
    this.getVideos();
  }

  getVideos() {
    VideoBackend.getVideos(this.props.account.name)
      .then((res) => {
        if (res.status === "ok") {
          this.setState({
            data: res.data,
            pagination: {
              ...this.state.pagination,
              total: this.state.pagination.total + 1,
            },
          });
        } else {
          Setting.showMessage("error", `Failed to get videos: ${res.msg}`);
        }
      });
  }

  newVideo() {
    const randomName = Setting.getRandomName();
    return {
      owner: this.props.account.name,
      name: `video_${randomName}`,
      createdTime: moment().format(),
      displayName: `New Video - ${randomName}`,
      videoId: "",
      coverUrl: "",
      labels: [],
      dataUrls: [],
      dataUrl: "",
      playAuth: "",
    };
  }

  deleteVideo(record) {
    VideoBackend.deleteVideo(record)
      .then((res) => {
        if (res.status === "ok") {
          Setting.showMessage("success", "Video deleted successfully");
          this.setState({
            data: this.state.data.filter((item) => item.name !== record.name),
            pagination: {
              ...this.state.pagination,
              total: this.state.pagination.total - 1,
            },
          });
        } else {
          Setting.showMessage("error", `Video failed to delete: ${res.msg}`);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `Video failed to delete: ${error}`);
      });
  }

  uploadFile(info) {
    const {status, response: res} = info.file;
    // if (status !== "uploading") {
    //   console.log(info.file, info.fileList);
    // }
    if (status === "done") {
      if (res.status === "ok") {
        Setting.showMessage("success", "Video uploaded successfully");
        const videoName = res.data;
        this.props.history.push(`/videos/${videoName}`);
      } else {
        Setting.showMessage("error", `Video failed to upload: ${res.msg}`);
      }
    } else if (status === "error") {
      Setting.showMessage("success", "Video failed to upload");
    }
  }

  renderUpload() {
    const props = {
      name: "file",
      accept: ".mp4",
      method: "post",
      action: `${Setting.ServerUrl}/api/upload-video`,
      withCredentials: true,
      onChange: (info) => {
        this.uploadFile(info);
      },
    };

    return (
      <Upload {...props}>
        <Button type="primary" size="small">
          <UploadOutlined /> {i18next.t("video:Upload Video")} (.mp4)
        </Button>
      </Upload>
    );
  }

  renderTable(videos) {
    const columns = [
      {
        title: i18next.t("general:Name"),
        dataIndex: "name",
        key: "name",
        width: "180px",
        sorter: (a, b) => a.name.localeCompare(b.name),
        render: (text, record, index) => {
          text = text.replace(".MP4", ".mp4");
          return (
            <Link to={`/videos/${text}`}>
              {text}
            </Link>
          );
        },
      },
      // {
      //   title: i18next.t("general:Display name"),
      //   dataIndex: "displayName",
      //   key: "displayName",
      //   width: "200px",
      //   sorter: (a, b) => a.displayName.localeCompare(b.displayName),
      // },
      {
        title: i18next.t("video:Tag"),
        dataIndex: "tag",
        key: "tag",
        width: "120px",
        sorter: (a, b) => a.tag.localeCompare(b.tag),
      },
      {
        title: i18next.t("video:School"),
        dataIndex: "school",
        key: "school",
        width: "80px",
        sorter: (a, b) => a.school.localeCompare(b.school),
      },
      {
        title: i18next.t("video:Grade"),
        dataIndex: "grade",
        key: "grade",
        width: "80px",
        sorter: (a, b) => a.grade.localeCompare(b.grade),
      },
      {
        title: i18next.t("video:Subject"),
        dataIndex: "subject",
        key: "subject",
        width: "80px",
        sorter: (a, b) => a.subject.localeCompare(b.subject),
      },
      {
        title: i18next.t("video:Video ID"),
        dataIndex: "videoId",
        key: "videoId",
        width: "250px",
        sorter: (a, b) => a.videoId.localeCompare(b.videoId),
      },
      {
        title: i18next.t("video:Cover"),
        dataIndex: "coverUrl",
        key: "coverUrl",
        width: "170px",
        render: (text, record, index) => {
          return (
            <a target="_blank" rel="noreferrer" href={text}>
              <img src={text} alt={text} width={150} />
            </a>
          );
        },
      },
      // {
      //   title: i18next.t("video:Labels"),
      //   dataIndex: "labels",
      //   key: "labels",
      //   width: "120px",
      //   sorter: (a, b) => a.factors.localeCompare(b.factors),
      //   render: (text, record, index) => {
      //     return Setting.getLabelTags(text);
      //   },
      // },
      {
        title: i18next.t("video:Label count"),
        dataIndex: "labelCount",
        key: "labelCount",
        width: "90px",
        sorter: (a, b) => a.labelCount - b.labelCount,
      },
      {
        title: i18next.t("video:Segment count"),
        dataIndex: "segmentCount",
        key: "segmentCount",
        width: "110px",
        sorter: (a, b) => a.segmentCount - b.segmentCount,
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: "action",
        key: "action",
        width: "220px",
        render: (text, record, index) => {
          return (
            <div>
              <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}} onClick={() => Setting.openLink(`/videos/${record.name}`)}>{i18next.t("general:Open")}</Button>
              <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}} type="primary" onClick={() => this.props.history.push(`/videos/${record.name}`)}>{i18next.t("general:Edit")}</Button>
              <Popconfirm
                title={`${i18next.t("general:Sure to delete")}: ${record.name} ?`}
                onConfirm={() => this.deleteVideo(record)}
                okText={i18next.t("general:OK")}
                cancelText={i18next.t("general:Cancel")}
              >
                <Button style={{marginBottom: "10px"}} type="primary" danger>{i18next.t("general:Delete")}</Button>
              </Popconfirm>
            </div>
          );
        },
      },
    ];

    const paginationProps = {
      total: this.state.pagination.total,
      showQuickJumper: true,
      showSizeChanger: true,
      showTotal: () => i18next.t("general:{total} in total").replace("{total}", this.state.pagination.total),
    };

    return (
      <div>
        <Table scroll={{x: "max-content"}} columns={columns} dataSource={videos} rowKey="name" size="middle" bordered pagination={paginationProps}
          title={() => (
            <div>
              {i18next.t("general:Videos")}
              {/* &nbsp;&nbsp;&nbsp;&nbsp;*/}
              {/* <Button type="primary" size="small" onClick={this.addVideo.bind(this)}>{i18next.t("general:Add")}</Button>*/}
                   &nbsp;&nbsp;&nbsp;&nbsp;
              {
                this.renderUpload()
              }
            </div>
          )}
          loading={videos === null}
          onChange={this.handleTableChange}
        />
      </div>
    );
  }

  fetch = (params = {}) => {
    const field = params.searchedColumn, value = params.searchText;
    const sortField = params.sortField, sortOrder = params.sortOrder;
    this.setState({loading: true});
    VideoBackend.getVideos("admin", params.pagination.current, params.pagination.pageSize, field, value, sortField, sortOrder)
      .then((res) => {
        this.setState({
          loading: false,
        });
        if (res.status === "ok") {
          this.setState({
            data: res.data,
            pagination: {
              ...params.pagination,
              total: res.data2,
            },
            searchText: params.searchText,
            searchedColumn: params.searchedColumn,
          });
        } else {
          if (Setting.isResponseDenied(res)) {
            this.setState({
              isAuthorized: false,
            });
          } else {
            Setting.showMessage("error", res.msg);
          }
        }
      });
  };
}

export default VideoListPage;

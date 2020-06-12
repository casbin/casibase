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

import React from "react";
import * as Setting from "../Setting";
import * as TopicBackend from "../backend/TopicBackend";
import {withRouter} from "react-router-dom";
import Avatar from "../Avatar";

class AllCreatedTopicsBox extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            classes: props,
            memberId: props.match.params.memberId,
            tab: props.match.params.tab,
            topics: [],
            TAB_LIST: [
                {label: "Q&A", value: "qna"},
                {label: "Tech", value: "tech"},
                {label: "Play", value: "play"},
                {label: "Jobs", value: "jobs"},
                {label: "Deals", value: "deals"},
                {label: "City", value: "city"}
            ]
        };
    }

    componentDidMount() {
        this.getAllCreatedTopics();
    }

    getAllCreatedTopics() {
        TopicBackend.getAllCreatedTopics(this.state.memberId, this.state.tab, 10, 1)
            .then((res) => {
                this.setState({
                    topics: res,
                });
            });
    }

    renderTab(tab) {
        return (
                {
                    ...this.state.tab === tab.value ?
                        <a href={`/member/${this.state.memberId}/${tab.value}`} class="cell_tab_current"> {tab.label} </a> :
                        <a href={`/member/${this.state.memberId}/${tab.value}`} class="cell_tab"> {tab.label} </a>
                }
        )
    }

    renderTopic(topic) {
        const style = topic.nodeId !== "promotions" ? null : {
            backgroundImage: `url('${Setting.getStatic("/static/img/corner_star.png")}')`,
            backgroundRepeat: "no-repeat",
            backgroundSize: "20px 20px",
            backgroundPosition: "right top"
        };

        return (
            <div className="cell item" style={style}>
                <table cellPadding="0" cellSpacing="0" border="0" width="100%">
                    <tbody>
                    <tr>
                        <td width="auto" valign="middle">
                            <span className="item_title">
                                <a href={`/t/${topic.id}`} className="topic-link"> {topic.title} </a>
                            </span>
                            <div className="sep5" />
                            <span className="topic_info">
                                <div className="votes" />
                                <a className="node" href={`/go/${topic.nodeId}`}> {topic.nodeName} </a>
                                &nbsp;•&nbsp;
                                <strong><a href={`/member/${topic.author}`}> {topic.author} </a></strong>
                                &nbsp;•&nbsp;
                                {Setting.getPrettyDate(topic.createdTime)}
                                &nbsp;•&nbsp;
                                last reply from
                                <strong>
                                    <a href={`/member/${topic.lastReplyUser}`}> {topic.lastReplyUser} </a>
                                </strong>
                            </span>
                        </td>
                        <td width="70" align="right" valign="middle">
                            <a href={`/t/${topic.id}`} className="count_livid">6</a>
                        </td>
                    </tr>
                    </tbody>
                </table>
            </div>
        )
    }

    render() {
        {
            if (this.state.tab === "replies") {
                return (
                    <div />
                );
            }
        }
        return (
            <div className="box">
                        <div class="cell_tabs">
                            <div class="fl">
                                <Avatar username={this.state.memberId} size={"small"} />
                            </div>
                            {
                                this.state.tab === undefined ?
                                    <a href={`/member/${this.state.memberId}`} class="cell_tab_current"> {`${this.state.memberId}'s all topics`} </a> :
                                    <a href={`/member/${this.state.memberId}`} class="cell_tab"> {`${this.state.memberId}'s all topics`} </a>
                            }
                            {
                                this.state.TAB_LIST.map((tab) => {
                                    return this.renderTab(tab);
                                    }
                                )
                            }
                        </div>
                                {
                                    this.state.topics.map((topic) => {
                                        return this.renderTopic(topic);
                                        }
                                    )
                                }
                        {
                            this.state.tab === undefined ?
                                <div className="inner"><span className="chevron">»</span>
                                    <a href={`/member/${this.state.memberId}/topics`}> {`${this.state.memberId}'s more topics`} </a>
                                </div> :
                                <div />
                        }

                </div>
        );
    }
}

export default withRouter(AllCreatedTopicsBox);

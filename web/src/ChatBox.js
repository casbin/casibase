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
import {Avatar, ChatContainer, ConversationHeader, MainContainer, Message, MessageInput, MessageList} from "@chatscope/chat-ui-kit-react";
import "@chatscope/chat-ui-kit-styles/dist/default/styles.min.css";
import {marked} from "marked";
import DOMPurify from "dompurify";
import katex from "katex";
import "katex/dist/katex.min.css";
import hljs from "highlight.js";
import "highlight.js/styles/atom-one-dark-reasonable.css";
import * as Conf from "./Conf";
import * as Setting from "./Setting";
import i18next from "i18next";

marked.setOptions({
  renderer: new marked.Renderer(),
  gfm: true,
  tables: true,
  breaks: true,
  pedantic: false,
  sanitize: false,
  smartLists: true,
  smartypants: true,
});

class ChatBox extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      dots: ".",
      value: "",
    };
    this.timer = null;
  }

  componentDidMount() {
    this.timer = setInterval(() => {
      this.setState(prevState => {
        switch (prevState.dots) {
        case ".":
          return {dots: ".."};
        case "..":
          return {dots: "..."};
        case "...":
          return {dots: "."};
        default:
          return {dots: "."};
        }
      });
    }, 500);
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }

  handleSend = (innerHtml) => {
    if (this.state.value === "" || this.props.disableInput) {
      return;
    }
    const newValue = this.state.value.replace(/<img src="([^"]*)" alt="([^"]*)" width="(\d+)" height="(\d+)" data-original-width="(\d+)" data-original-height="(\d+)">/g, (match, src, alt, width, height, scaledWidth, scaledHeight) => {
      return `<img src="${src}" alt="${alt}" width="${scaledWidth}" height="${scaledHeight}">`;
    });
    let fileName = "";
    if (this.inputImage.files[0]) {
      fileName = this.inputImage.files[0].name;
    }
    this.props.sendMessage(newValue, fileName);
    this.setState({value: ""});
  };

  handleImageClick = () => {
    this.inputImage.click();
  };

  handleInputChange = async() => {
    const file = this.inputImage.files[0];
    const reader = new FileReader();
    reader.onload = (e) => {
      const img = new Image();
      img.onload = () => {
        const originalWidth = img.width;
        const originalHeight = img.height;
        const inputMaxWidth = 70;
        const chatMaxWidth = 600;
        let Ratio = 1;
        if (originalWidth > inputMaxWidth) {
          Ratio = inputMaxWidth / originalWidth;
        }
        const scaledWidth = Math.round(originalWidth * Ratio);
        const scaledHeight = Math.round(originalHeight * Ratio);
        if (originalWidth > chatMaxWidth) {
          Ratio = chatMaxWidth / originalWidth;
        }
        const chatScaledWidth = Math.round(originalWidth * Ratio);
        const chatScaledHeight = Math.round(originalHeight * Ratio);
        this.setState({
          value: this.state.value + `<img src="${img.src}" alt="${img.alt}" width="${scaledWidth}" height="${scaledHeight}" data-original-width="${chatScaledWidth}" data-original-height="${chatScaledHeight}">`,
        });
      };

      img.src = e.target.result;
    };
    reader.readAsDataURL(file);
  };

  renderMarkdown(text) {
    if (text === "") {
      text = this.state.dots;
    }
    const rawHtml = marked(text);
    let cleanHtml = DOMPurify.sanitize(rawHtml);
    /* replace <p></p> with <div></div>, reduce paragraph spacing. */
    cleanHtml = cleanHtml.replace(/<p>/g, "<div>").replace(/<\/p>/g, "</div>");
    /* h2 is larger than h1, h2 is the largest, so replace h1 with h2, and set margin as 20px. */
    cleanHtml = cleanHtml.replace(/<h1>/g, "<h2>").replace(/<(h[1-6])>/g, "<$1 style='margin-top: 20px; margin-bottom: 20px'>");
    /* adjust margin and internal gap for unordered list and ordered list. */
    cleanHtml = cleanHtml.replace(/<(ul)>/g, "<ul style='display: flex; flex-direction: column; gap: 10px; margin-top: 10px; margin-bottom: 10px'>").replace(/<(ol)>/g, "<ol style='display: flex; flex-direction: column; gap: 0px; margin-top: 20px; margin-bottom: 20px'>");
    /* adjust code block, for auto line feed. */
    cleanHtml = cleanHtml.replace(/<pre>/g, "<pre style='white-space: pre-wrap; white-space: -moz-pre-wrap; white-space: -pre-wrap; white-space: -o-pre-wrap; word-wrap: break-word;'>");
    return cleanHtml;
  }

  renderLatex(text) {
    const inlineLatexRegex = /\(\s*(([a-zA-Z])|(\\.+?)|([^)]*?[_^!].*?))\s*\)/g;
    const blockLatexRegex = /\[\s*(.+?)\s*\]/g;

    text = text.replace(blockLatexRegex, (match, formula) => {
      try {
        return katex.renderToString(formula, {throwOnError: false, displayMode: true});
      } catch (error) {
        return match;
      }
    });

    return text.replace(inlineLatexRegex, (match, formula) => {
      try {
        return katex.renderToString(formula, {throwOnError: false, displayMode: false});
      } catch (error) {
        return match;
      }
    });
  }

  renderCode(text) {
    const tempDiv = document.createElement("div");
    tempDiv.innerHTML = text;
    tempDiv.querySelectorAll("pre code").forEach((block) => {
      hljs.highlightBlock(block);
    });
    return tempDiv.innerHTML;
  }

  renderText(text) {
    let html;
    html = this.renderMarkdown(text);
    html = this.renderLatex(html);
    html = this.renderCode(html);
    return <div dangerouslySetInnerHTML={{__html: html}} style={{display: "flex", flexDirection: "column", gap: "0px"}} />;
  }

  render() {
    let title = Setting.getUrlParam("title");
    if (title === null) {
      title = Conf.AiName;
    }

    let messages = this.props.messages;
    if (messages === null) {
      messages = [];
    }
    return (
      <React.Fragment>
        <MainContainer style={{display: "flex", width: "100%", height: "100%", border: "1px solid rgb(242,242,242)", borderRadius: "6px"}} >
          <ChatContainer style={{display: "flex", width: "100%", height: "100%"}}>
            {
              (title === "") ? null : (
                <ConversationHeader style={{backgroundColor: "rgb(246,240,255)", height: "42px"}}>
                  <ConversationHeader.Content userName={title} />
                </ConversationHeader>
              )
            }
            <MessageList style={{marginTop: "10px"}}>
              {messages.filter(message => message.isHidden === false).map((message, index) => (
                <Message key={index} model={{
                  type: "custom",
                  sender: message.name,
                  direction: message.author === "AI" ? "incoming" : "outgoing",
                }} avatarPosition={message.author === "AI" ? "tl" : "tr"}>
                  <Avatar src={message.author === "AI" ? Conf.AiAvatar : (this.props.hideInput === true ? "https://cdn.casdoor.com/casdoor/resource/built-in/admin/casibase-user.png" : this.props.account.avatar)} name="GPT" />
                  <Message.CustomContent>
                    {this.renderText(message.text)}
                  </Message.CustomContent>
                </Message>
              ))}
            </MessageList>
            {
              this.props.hideInput === true ? null : (
                <MessageInput disabled={false}
                  sendDisabled={this.state.value === "" || this.props.disableInput}
                  placeholder={i18next.t("chat:Type message here")}
                  onSend={this.handleSend}
                  onChange={(val) => {
                    this.setState({value: val});
                  }}
                  value={this.state.value}
                  onPaste={(evt) => {
                    evt.preventDefault();
                    this.setState({value: this.state.value + evt.clipboardData.getData("text")});
                  }}
                  onAttachClick={() => {
                    this.handleImageClick();
                  }}
                />
              )
            }
          </ChatContainer>
        </MainContainer>
        <input
          ref={e => this.inputImage = e}
          type="file"
          accept="image/*"
          multiple={false}
          onChange={this.handleInputChange}
          style={{display: "none"}}
        />
      </React.Fragment>
    );
  }
}

export default ChatBox;

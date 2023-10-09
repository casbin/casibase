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

package object

import (
	"fmt"

	"github.com/casibase/casibase/util"
	"github.com/casibase/casibase/video"
	"xorm.io/core"
)

type Label struct {
	Id        string  `xorm:"varchar(100)" json:"id"`
	StartTime float64 `json:"startTime"`
	EndTime   float64 `json:"endTime"`
	Text      string  `xorm:"varchar(100)" json:"text"`
}

type Video struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	DisplayName string `xorm:"varchar(500)" json:"displayName"`

	VideoId    string   `xorm:"varchar(100)" json:"videoId"`
	CoverUrl   string   `xorm:"varchar(200)" json:"coverUrl"`
	Labels     []*Label `xorm:"mediumtext" json:"labels"`
	DataUrls   []string `xorm:"mediumtext" json:"dataUrls"`
	DataUrl    string   `xorm:"varchar(200)" json:"dataUrl"`
	TagOnPause bool     `json:"tagOnPause"`

	PlayAuth string `xorm:"-" json:"playAuth"`
}

func GetGlobalVideos() ([]*Video, error) {
	videos := []*Video{}
	err := adapter.engine.Asc("owner").Desc("created_time").Find(&videos)
	if err != nil {
		return videos, err
	}

	return videos, nil
}

func GetVideos(owner string) ([]*Video, error) {
	videos := []*Video{}
	err := adapter.engine.Desc("created_time").Find(&videos, &Video{Owner: owner})
	if err != nil {
		return videos, err
	}

	return videos, nil
}

func getVideo(owner string, name string) (*Video, error) {
	v := Video{Owner: owner, Name: name}
	existed, err := adapter.engine.Get(&v)
	if err != nil {
		return &v, err
	}

	if existed {
		if v.VideoId != "" {
			v.PlayAuth = video.GetVideoPlayAuth(v.VideoId)
		}
		return &v, nil
	} else {
		return nil, nil
	}
}

func GetVideo(id string) (*Video, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getVideo(owner, name)
}

func UpdateVideo(id string, video *Video) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	_, err := getVideo(owner, name)
	if err != nil {
		return false, err
	}
	if video == nil {
		return false, nil
	}

	_, err = adapter.engine.ID(core.PK{owner, name}).AllCols().Update(video)
	if err != nil {
		return false, err
	}

	// return affected != 0
	return true, nil
}

func AddVideo(video *Video) (bool, error) {
	affected, err := adapter.engine.Insert(video)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func DeleteVideo(video *Video) (bool, error) {
	affected, err := adapter.engine.ID(core.PK{video.Owner, video.Name}).Delete(&Video{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func (video *Video) GetId() string {
	return fmt.Sprintf("%s/%s", video.Owner, video.Name)
}

func (video *Video) Populate() error {
	store, err := GetDefaultStore("admin")
	if err != nil {
		return err
	}
	if store == nil {
		return nil
	}

	dataUrls, err := store.GetVideoData()
	if err != nil {
		return err
	}

	video.DataUrls = dataUrls
	return nil
}

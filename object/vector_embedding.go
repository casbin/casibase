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
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/casbin/casibase/ai"
	"github.com/casbin/casibase/storage"
	"github.com/casbin/casibase/util"
	"golang.org/x/time/rate"
)

func filterTextFiles(files []*storage.Object) []*storage.Object {
	extSet := map[string]bool{
		".txt":  true,
		".md":   true,
		".docx": true,
		".doc":  false,
		".pdf":  true,
	}

	var res []*storage.Object
	for _, file := range files {
		ext := filepath.Ext(file.Key)
		if extSet[ext] {
			res = append(res, file)
		}
	}

	return res
}

func getTextFiles(provider string, prefix string) ([]*storage.Object, error) {
	files, err := storage.ListObjects(provider, prefix)
	if err != nil {
		return nil, err
	}

	return filterTextFiles(files), nil
}

func getObjectReadCloser(object *storage.Object) (io.ReadCloser, error) {
	resp, err := http.Get(object.Url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}
	return resp.Body, nil
}

func addEmbeddedVector(authToken string, text string, storeName string, fileName string) (bool, error) {
	embedding, err := ai.GetEmbeddingSafe(authToken, text)
	if err != nil {
		return false, err
	}

	displayName := text
	if len(text) > 25 {
		displayName = text[:25]
	}

	vector := &Vector{
		Owner:       "admin",
		Name:        fmt.Sprintf("vector_%s", util.GetRandomName()),
		CreatedTime: util.GetCurrentTime(),
		DisplayName: displayName,
		Store:       storeName,
		File:        fileName,
		Text:        text,
		Data:        embedding,
	}
	return AddVector(vector)
}

func setTextObjectVector(authToken string, provider string, key string, storeName string) (bool, error) {
	lb := rate.NewLimiter(rate.Every(time.Minute), 3)

	textObjects, err := getTextFiles(provider, key)
	if err != nil {
		return false, err
	}
	if len(textObjects) == 0 {
		return false, nil
	}

	for _, textObject := range textObjects {
		readCloser, err := getObjectReadCloser(textObject)
		if err != nil {
			return false, err
		}
		defer readCloser.Close()

		splitTxts := ai.GetSplitTxt(readCloser, textObject.Key)
		for _, splitTxt := range splitTxts {
			if lb.Allow() {
				success, err := addEmbeddedVector(authToken, splitTxt, storeName, textObject.Key)
				if err != nil {
					return false, err
				}
				if !success {
					return false, nil
				}
			} else {
				err := lb.Wait(context.Background())
				if err != nil {
					return false, err
				}
				success, err := addEmbeddedVector(authToken, splitTxt, storeName, textObject.Key)
				if err != nil {
					return false, err
				}
				if !success {
					return false, nil
				}
			}
		}
	}

	return true, nil
}

func getRelatedVectors(owner string) ([]*Vector, error) {
	vectors, err := GetVectors(owner)
	if err != nil {
		return nil, err
	}
	if len(vectors) == 0 {
		return nil, fmt.Errorf("no knowledge vectors found")
	}

	return vectors, nil
}

func GetNearestVectorText(authToken string, owner string, question string) (string, error) {
	qVector, err := ai.GetEmbeddingSafe(authToken, question)
	if err != nil {
		return "", err
	}
	if qVector == nil {
		return "", fmt.Errorf("no qVector found")
	}

	vectors, err := getRelatedVectors(owner)
	if err != nil {
		return "", err
	}

	var nVectors [][]float32
	for _, candidate := range vectors {
		nVectors = append(nVectors, candidate.Data)
	}

	i := ai.GetNearestVectorIndex(qVector, nVectors)
	return vectors[i].Text, nil
}

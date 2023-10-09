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

package embedding

import (
	"context"

	"github.com/casibase/casibase/proxy"
	"github.com/henomis/lingoose/embedder/huggingface"
)

type HuggingFaceEmbeddingProvider struct {
	subType   string
	secretKey string
}

func NewHuggingFaceEmbeddingProvider(subType string, secretKey string) (*HuggingFaceEmbeddingProvider, error) {
	return &HuggingFaceEmbeddingProvider{subType: subType, secretKey: secretKey}, nil
}

func (h *HuggingFaceEmbeddingProvider) QueryVector(text string, ctx context.Context) ([]float32, error) {
	client := huggingfaceembedder.New().WithToken(h.secretKey).WithModel(h.subType).WithHTTPClient(proxy.ProxyHttpClient)
	embed, err := client.Embed(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	return float64ToFloat32(embed[0]), nil
}

func float64ToFloat32(slice []float64) []float32 {
	newSlice := make([]float32, len(slice))
	for i, v := range slice {
		newSlice[i] = float32(v)
	}
	return newSlice
}

// Copyright 2024 The casbin Authors. All Rights Reserved.
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

	"github.com/casibase/casibase/embedding"
	"github.com/casibase/casibase/model"
	"github.com/casibase/casibase/util"
)

func getModelProviderFromName(owner string, providerName string) (*Provider, model.ModelProvider, error) {
	var provider *Provider
	var err error
	if providerName != "" {
		providerId := util.GetIdFromOwnerAndName(owner, providerName)
		provider, err = GetProvider(providerId)
	} else {
		provider, err = GetDefaultModelProvider()
	}
	if err != nil {
		return nil, nil, err
	}
	if provider == nil {
		if providerName != "" {
			return nil, nil, fmt.Errorf("The model provider: %s is not found", providerName)
		} else {
			return nil, nil, fmt.Errorf("Please add a model provider first")
		}
	}

	if provider.Category != "Model" {
		return nil, nil, fmt.Errorf("The model provider: %s is expected to be \"Model\" category, got: \"%s\"", provider.GetId(), provider.Category)
	}
	if provider.ClientSecret == "" && provider.Type != "Dummy" {
		return nil, nil, fmt.Errorf("The model provider: %s's client secret should not be empty", provider.GetId())
	}

	providerObj, err := provider.GetModelProvider()
	if err != nil {
		return nil, nil, err
	}

	return provider, providerObj, err
}

func getEmbeddingProviderFromName(owner string, providerName string) (*Provider, embedding.EmbeddingProvider, error) {
	var provider *Provider
	var err error
	if providerName != "" {
		providerId := util.GetIdFromOwnerAndName(owner, providerName)
		provider, err = GetProvider(providerId)
	} else {
		provider, err = GetDefaultEmbeddingProvider()
	}
	if err != nil {
		return nil, nil, err
	}
	if provider == nil {
		if providerName != "" {
			return nil, nil, fmt.Errorf("The embedding provider: %s is not found", providerName)
		} else {
			return nil, nil, fmt.Errorf("Please add an embedding provider first")
		}
	}

	if provider.Category != "Embedding" {
		return nil, nil, fmt.Errorf("The embedding provider: %s is expected to be \"Embedding\" category, got: \"%s\"", provider.GetId(), provider.Category)
	}
	if provider.ClientSecret == "" && provider.Type != "Dummy" {
		return nil, nil, fmt.Errorf("The embedding provider: %s's client secret should not be empty", provider.GetId())
	}

	providerObj, err := provider.GetEmbeddingProvider()
	if err != nil {
		return nil, nil, err
	}

	return provider, providerObj, err
}

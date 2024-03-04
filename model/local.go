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

package model

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/sashabaranov/go-openai"
)

type LocalModelProvider struct {
	typ              string
	subType          string
	deploymentName   string
	secretKey        string
	temperature      float32
	topP             float32
	frequencyPenalty float32
	presencePenalty  float32
	providerUrl      string
	apiVersion       string
}

func NewLocalModelProvider(typ string, subType string, secretKey string, temperature float32, topP float32, frequencyPenalty float32, presencePenalty float32, providerUrl string) (*LocalModelProvider, error) {
	p := &LocalModelProvider{
		typ:              typ,
		subType:          subType,
		secretKey:        secretKey,
		temperature:      temperature,
		topP:             topP,
		frequencyPenalty: frequencyPenalty,
		presencePenalty:  presencePenalty,
		providerUrl:      providerUrl,
	}
	return p, nil
}

func getLocalClientFromUrl(authToken string, url string) *openai.Client {
	config := openai.DefaultConfig(authToken)
	config.BaseURL = url

	c := openai.NewClientWithConfig(config)
	return c
}

func (p *LocalModelProvider) GetPricing() string {
	return `URL:
https://azure.microsoft.com/en-us/pricing/details/cognitive-services/openai-service/

Language models:

| Models                | Context | Input (Per 1,000 tokens) | Output (Per 1,000 tokens) |
|-----------------------|---------|--------------------------|--------------------------|
| GPT-3.5-Turbo-0125    | 16K     | $0.0005                  | $0.0015                  |
| GPT-3.5-Turbo-Instruct| 4K      | $0.0015                  | $0.002                   |
| GPT-4-Turbo           | 128K    | $0.01                    | $0.03                    |
| GPT-4-Turbo-Vision    | 128K    | $0.01                    | $0.03                    |
| GPT-4                 | 8K      | $0.03                    | $0.06                    |
| GPT-4                 | 32K     | $0.06                    | $0.12                    |

Image models:

| Models   | Quality | Resolution               | Price (per 100 images) |
|----------|---------|--------------------------|------------------------|
| Dall-E-3 | Standard| 1024 * 1024              | N/A                    |
|          | Standard| 1024 * 1792, 1792 * 1024 | $8                     |
| Dall-E-3 | HD      | 1024 * 1024              | N/A                    |
|          | HD      | 1024 * 1792, 1792 * 1024 | N/A                    |
| Dall-E-2 | Standard| 1024 * 1024              | N/A                    |
`
}

func (p *LocalModelProvider) calculatePrice(modelResult *ModelResult) error {
	model := p.subType
	var inputPricePerThousandTokens, outputPricePerThousandTokens float64
	switch {
	case strings.Contains(model, "gpt-3.5-turbo-instruct"):
		inputPricePerThousandTokens = 0.0015
		outputPricePerThousandTokens = 0.002
	case strings.Contains(model, "gpt-3.5-turbo"):
		inputPricePerThousandTokens = 0.0005
		outputPricePerThousandTokens = 0.0015
	case strings.Contains(model, "gpt-4-turbo"), strings.Contains(model, "gpt-4-vision"):
		inputPricePerThousandTokens = 0.01
		outputPricePerThousandTokens = 0.03
	case strings.Contains(model, "gpt-4") && strings.Contains(model, "32k"):
		inputPricePerThousandTokens = 0.06
		outputPricePerThousandTokens = 0.12
	case strings.Contains(model, "gpt-4"):
		inputPricePerThousandTokens = 0.03
		outputPricePerThousandTokens = 0.06
	case strings.Contains(model, "dall-e-3"):
		modelResult.TotalPrice = float64(modelResult.ImageCount) * 0.08
		return nil
	default:
		return fmt.Errorf("calculatePrice() error: unknown model type: %s", model)
	}

	inputPrice := getPrice(modelResult.PromptTokenCount, inputPricePerThousandTokens)
	outputPrice := getPrice(modelResult.ResponseTokenCount, outputPricePerThousandTokens)
	modelResult.TotalPrice = AddPrices(inputPrice, outputPrice)
	modelResult.Currency = "USD"
	return nil
}

func (p *LocalModelProvider) CalculatePrice(modelResult *ModelResult) error {
	return p.calculatePrice(modelResult)
}

func flushDataAzure(data string, writer io.Writer) error {
	flusher, ok := writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("writer does not implement http.Flusher")
	}
	for _, runeValue := range data {
		char := string(runeValue)
		_, err := fmt.Fprintf(writer, "event: message\ndata: %s\n\n", char)
		if err != nil {
			return err
		}

		flusher.Flush()

		delta := 0
		if !unicode.In(runeValue, unicode.Latin) {
			delta = 50
		}

		var delay int
		if char == "," || char == "，" {
			delay = 20 + rand.Intn(50) + delta
		} else if char == "." || char == "。" || char == "!" || char == "！" || char == "?" || char == "？" {
			delay = 50 + rand.Intn(50) + delta
		} else if char == " " || char == "　" || char == "(" || char == "（" || char == ")" || char == "）" {
			delay = 10 + rand.Intn(50) + delta
		} else {
			delay = rand.Intn(1 + delta*2/5)
		}

		if unicode.In(runeValue, unicode.Latin) {
			delay -= 20
		}

		if delay > 0 {
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}
	return nil
}

func flushDataOpenai(data string, writer io.Writer) error {
	flusher, ok := writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("writer does not implement http.Flusher")
	}
	if _, err := fmt.Fprintf(writer, "event: message\ndata: %s\n\n", data); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}

func (p *LocalModelProvider) QueryText(question string, writer io.Writer, history []*RawMessage, prompt string, knowledgeMessages []*RawMessage) (*ModelResult, error) {
	var client *openai.Client
	var flushData func(string, io.Writer) error
	if p.typ == "Local" {
		client = getLocalClientFromUrl(p.secretKey, p.providerUrl)
		flushData = flushDataOpenai
	} else if p.typ == "Azure" {
		client = getAzureClientFromToken(p.deploymentName, p.secretKey, p.providerUrl, p.apiVersion)
		flushData = flushDataAzure
	} else if p.typ == "OpenAI" {
		client = getOpenAiClientFromToken(p.secretKey)
		flushData = flushDataOpenai
	}

	ctx := context.Background()
	flusher, ok := writer.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("writer does not implement http.Flusher")
	}

	model := p.subType
	temperature := p.temperature
	topP := p.topP
	frequencyPenalty := p.frequencyPenalty
	presencePenalty := p.presencePenalty

	maxTokens := GetOpenAiMaxTokens(model)

	modelResult := &ModelResult{}
	if getOpenAiModelType(p.subType) == "Chat" {
		if p.subType == "dall-e-3" {
			reqUrl := openai.ImageRequest{
				Prompt:         question,
				Model:          openai.CreateImageModelDallE3,
				Size:           openai.CreateImageSize1024x1024,
				ResponseFormat: openai.CreateImageResponseFormatURL,
				N:              1,
			}

			respUrl, err := client.CreateImage(ctx, reqUrl)
			if err != nil {
				return nil, err
			}

			url := fmt.Sprintf("<img src=\"%s\" width=\"100%%\" height=\"auto\">", respUrl.Data[0].URL)
			fmt.Fprint(writer, url)
			flusher.Flush()

			modelResult.ImageCount = 1
			err = p.calculatePrice(modelResult)
			if err != nil {
				return nil, err
			}

			return modelResult, nil
		}

		rawMessages, err := OpenaiGenerateMessages(prompt, question, history, knowledgeMessages, model, maxTokens)
		if err != nil {
			return nil, err
		}

		var messages []openai.ChatCompletionMessage
		if p.subType == "gpt-4-vision-preview" {
			messages, err = OpenaiRawMessagesToGpt4VisionMessages(rawMessages)
			if err != nil {
				return nil, err
			}
		} else {
			messages = OpenaiRawMessagesToMessages(rawMessages)
		}

		// https://github.com/sashabaranov/go-openai/pull/223#issuecomment-1494372875
		promptTokenCount, err := OpenaiNumTokensFromMessages(messages, p.subType)
		if err != nil {
			return nil, err
		}

		modelResult.PromptTokenCount = promptTokenCount
		modelResult.TotalTokenCount = modelResult.PromptTokenCount + modelResult.ResponseTokenCount
		err = p.calculatePrice(modelResult)
		if err != nil {
			return nil, err
		}

		respStream, err := client.CreateChatCompletionStream(
			ctx,
			ChatCompletionRequest(model, messages, temperature, topP, frequencyPenalty, presencePenalty),
		)
		if err != nil {
			return nil, err
		}
		defer respStream.Close()

		isLeadingReturn := true
		for {
			completion, streamErr := respStream.Recv()
			if streamErr != nil {
				if streamErr == io.EOF {
					break
				}
				return nil, streamErr
			}

			if len(completion.Choices) == 0 {
				continue
			}

			data := completion.Choices[0].Delta.Content
			if isLeadingReturn && len(data) != 0 {
				if strings.Count(data, "\n") == len(data) {
					continue
				} else {
					isLeadingReturn = false
				}
			}

			err = flushData(data, writer)
			if err != nil {
				return nil, err
			}

			// https://github.com/sashabaranov/go-openai/pull/223#issuecomment-1494372875
			var responseTokenCount int
			responseTokenCount, err = GetTokenSize(p.subType, data)
			if err != nil {
				return nil, err
			}

			modelResult.ResponseTokenCount += responseTokenCount
			modelResult.TotalTokenCount = modelResult.PromptTokenCount + modelResult.ResponseTokenCount
			err = p.calculatePrice(modelResult)
			if err != nil {
				return nil, err
			}
		}

		return modelResult, nil
	} else if getOpenAiModelType(p.subType) == "Completion" {
		respStream, err := client.CreateCompletionStream(
			ctx,
			openai.CompletionRequest{
				Model:            model,
				Prompt:           question,
				Stream:           true,
				Temperature:      temperature,
				TopP:             topP,
				FrequencyPenalty: frequencyPenalty,
				PresencePenalty:  presencePenalty,
			},
		)
		if err != nil {
			return nil, err
		}
		defer respStream.Close()

		isLeadingReturn := true
		var response strings.Builder
		for {
			completion, streamErr := respStream.Recv()
			if streamErr != nil {
				if streamErr == io.EOF {
					break
				}
				return nil, streamErr
			}

			data := completion.Choices[0].Text
			if isLeadingReturn && len(data) != 0 {
				if strings.Count(data, "\n") == len(data) {
					continue
				} else {
					isLeadingReturn = false
				}
			}

			err = flushData(data, writer)
			if err != nil {
				return nil, err
			}

			_, err = response.WriteString(data)
			if err != nil {
				return nil, err
			}
		}

		modelResult, err = getDefaultModelResult(model, question, response.String())
		return modelResult, nil
	} else {
		return nil, fmt.Errorf("QueryText() error: unknown model type: %s", p.subType)
	}
}

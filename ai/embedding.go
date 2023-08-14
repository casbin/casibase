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

package ai

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/sashabaranov/go-openai"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

func splitTxt(f io.ReadCloser) []string {
	const maxLength = 210 * 3
	scanner := bufio.NewScanner(f)
	var res []string
	var temp string

	for scanner.Scan() {
		line := scanner.Text()
		if len(temp)+len(line) <= maxLength {
			temp += line
		} else {
			res = append(res, temp)
			temp = line
		}
	}

	if len(temp) > 0 {
		res = append(res, temp)
	}

	return res
}

func GetSplitTxt(f io.ReadCloser) []string {
	return splitTxt(f)
}

func getEmbedding(authToken string, input []string, timeout int) ([]float32, error) {
	client := getProxyClientFromToken(authToken)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30+timeout*2)*time.Second)
	defer cancel()

	resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: input,
		Model: openai.AdaEmbeddingV2,
	})
	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}

func GetEmbeddingSafe(authToken string, input []string) ([]float32, error) {
	var embedding []float32
	var err error
	for i := 0; i < 10; i++ {
		embedding, err = getEmbedding(authToken, input, i)
		if err != nil {
			if i > 0 {
				fmt.Printf("\tFailed (%d): %s\n", i+1, err.Error())
			}
		} else {
			break
		}
	}

	if err != nil {
		return nil, err
	} else {
		return embedding, nil
	}
}

func getMostSimilarVector(target []float64, candidates [][]float64) (int, float64) {
	return getMostSimilarVectorInternal(target, candidates)
}

func GetMostSimilarVectorFromFloat32(target []float32, candidates [][]float64) (int, float64) {
	return getMostSimilarVectorInternal(Float32To64(target), candidates)
}

func getMostSimilarVectorInternal(target []float64, candidates [][]float64) (int, float64) {
	targetVec := mat.NewVecDense(len(target), target)
	var mostSimilarIndex int
	var maxSimilarity float64 = -1

	for i, candidate := range candidates {
		candidateVec := mat.NewVecDense(len(candidate), candidate)
		similarity := cosineSimilarity(targetVec, candidateVec)
		if similarity > maxSimilarity {
			maxSimilarity = similarity
			mostSimilarIndex = i
		}
	}

	return mostSimilarIndex, maxSimilarity
}

func cosineSimilarity(vec1, vec2 *mat.VecDense) float64 {
	dotProduct := floats.Dot(vec1.RawVector().Data, vec2.RawVector().Data)
	normVec1 := floats.Norm(vec1.RawVector().Data, 2)
	normVec2 := floats.Norm(vec2.RawVector().Data, 2)

	if normVec1 == 0 || normVec2 == 0 {
		return 0.0
	}

	return dotProduct / (normVec1 * normVec2)
}

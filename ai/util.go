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

import "github.com/pkoukk/tiktoken-go"

func GetTokenSize(model string, prompt string) (int, error) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return 0, err
	}

	token := tkm.Encode(prompt, nil, nil)
	res := len(token)
	return res, nil
}

func Float32To64(f []float32) []float64 {
	res := make([]float64, len(f))
	for i, v := range f {
		res[i] = float64(v)
	}
	return res
}

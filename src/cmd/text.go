// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"strings"

	"github.com/cyy0523xc/nlp-corpus-tools/common"
	"github.com/cyy0523xc/nlp-corpus-tools/text"
	"github.com/spf13/cobra"
)

type TextParams struct {
	Fieldnames []string

	// iconv
	From string
	To   string

	// replace
	Old, New string

	// trim
	Cutset string
}

var textParams TextParams
var textActions = []string{"replace", "filter-space", "filter-newline-char", "iconv", "dbc2sbc", "trim"}

// textCmd represents the text command
var textCmd = &cobra.Command{
	Use:   "text",
	Short: "对文本的操作",
	Long: `对需要进行nlp处理的文本进行操作，例如编码转化，全角半角转换，字符串替换等。
For example:

    nlp-corpus-tools text --help

支持以下功能:
- [x] replace: 字符串替换
- [x] filter-space: 过滤掉文本中的空格，制表符等
- [x] filter-newline-char: 过滤掉文本中的换行符
- [x] iconv: 编码转换
- [x] dbc2sbc: 全角转半角
- [x] trim: 替换字符串前后的空格字符
`,
	Example: `
1. 将文本的gbk编码转化为utf8编码:

    nlp-corpus-tools text -a iconv -i test.csv --from=gbk --to=utf8
	cat test.csv|nlp-corpus-tools text -a iconv --from=gbk --to=utf8

2. 字符串替换：

    nlp-corpus-tools text -a replace -i test.csv --old="友邻茉莉花开" --new="hello"
	cat test.csv|nlp-corpus-tools text -a replace --old="友邻茉莉花开" --new="hello"

`,

	Run: func(cmd *cobra.Command, args []string) {
		rootParams.checkAction(textActions)

		r, err := rootParams.getInput()
		if err != nil {
			panic(err)
		}
		defer r.Close()

		w, err := rootParams.getOutput()
		if err != nil {
			panic(err)
		}
		defer w.Close()

		if common.Debug {
			fmt.Printf("text params: %+v\n", textParams)
		}

		funcMap := map[string]func(string) (string, error){
			"replace": TextReplace,
			"dbc2sbc": TextDBC2SBC,
		}
		if err = text.Transform(r, w, textParams.Fieldnames, funcMap[rootParams.action]); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(textCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	textCmd.PersistentFlags().StringVar(&textParams.From, "from", "gbk", "转换前的编码")
	textCmd.PersistentFlags().StringVar(&textParams.To, "to", "utf8", "转换后的编码")

	textCmd.PersistentFlags().StringVar(&textParams.Old, "old", "", "被替换的字符串")
	textCmd.PersistentFlags().StringVar(&textParams.New, "new", "", "新字符串")

	textCmd.PersistentFlags().StringVar(&textParams.Cutset, "cutset", "", "trim需要被替换的字符串")

	textCmd.PersistentFlags().StringArrayVar(&textParams.Fieldnames, "fields", []string{}, "需要转换的字段名，默认为全部")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// textCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func TextReplace(str string) (string, error) {
	return strings.Replace(str, textParams.Old, textParams.New, -1), nil
}

func TextFilterSpace(str string) (string, error) {
	str = strings.Replace(str, " ", "", -1)
	return strings.Replace(str, "\t", "", -1), nil
}

func TextFilterNewlineChar(str string) (string, error) {
	str = strings.Replace(str, "\r\n", " ", -1)
	return strings.Replace(str, "\n", " ", -1), nil
}

func TextTrim(str string) (string, error) {
	if textParams.Cutset != "" {
		return strings.Trim(str, textParams.Cutset)
	}
	return strings.TrimSpace(str)
}

func TextIconv(str string) (string, error) {

	return str, nil
}

func TextDBC2SBC(str string) (string, error) {
	return common.DBC2SBC(str), nil
}

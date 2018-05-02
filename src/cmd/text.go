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

3. 去掉前后的空格字符：

    nlp-corpus-tools text -a trim -i test.csv 
	cat test.csv|nlp-corpus-tools text -a trim

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

		switch rootParams.action {
		case "iconv":
			textParams.From = strings.ToLower(textParams.From)
			textParams.To = strings.ToLower(textParams.To)
			if !common.StringIn(textParams.From, common.Charsets) {
				panic("from参数不支持")
			}
			if !common.StringIn(textParams.To, common.Charsets) {
				panic("to参数不支持")
			}
			if textParams.From == textParams.To {
				panic("from和to参数应该相同")
			}
			if textParams.From != common.UTF8 && textParams.To != common.UTF8 {
				panic("from或者to参数至少需要一个为：" + common.UTF8)
			}
		}

		funcMap := map[string]func(string) (string, error){
			"replace":             TextReplace,
			"filter-space":        TextFilterSpace,
			"filter-newline_char": TextFilterNewlineChar,
			"dbc2sbc":             TextDBC2SBC,
			"iconv":               TextIconv,
			"trim":                TextTrim,
		}
		if err = text.Transform(r, w, textParams.Fieldnames, funcMap[rootParams.action]); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(textCmd)

	textCmd.PersistentFlags().StringArrayVar(&textParams.Fieldnames, "fields", []string{}, "指定转换的字段名，默认对全部的字段进行处理")

	textCmd.PersistentFlags().StringVar(&textParams.From, "from", "gbk", "转换前的编码，可以取值："+strings.Join(common.Charsets, ", "))
	textCmd.PersistentFlags().StringVar(&textParams.To, "to", "utf8", "转换后的编码，可以取值："+strings.Join(common.Charsets, ", "))

	textCmd.PersistentFlags().StringVar(&textParams.Old, "old", "", "被替换的字符串")
	textCmd.PersistentFlags().StringVar(&textParams.New, "new", "", "新字符串")

	textCmd.PersistentFlags().StringVar(&textParams.Cutset, "cutset", "", "trim需要被替换的字符串")
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
		return strings.Trim(str, textParams.Cutset), nil
	}
	return strings.TrimSpace(str), nil
}

func TextIconv(str string) (string, error) {
	if textParams.From == common.UTF8 {
		return common.UTF8_GB(str)
	}
	return common.GB_UTF8(str)
}

func TextDBC2SBC(str string) (string, error) {
	return common.DBC2SBC(str), nil
}

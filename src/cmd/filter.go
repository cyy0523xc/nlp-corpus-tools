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
	"github.com/cyy0523xc/nlp-corpus-tools/filter"
	"github.com/spf13/cobra"
)

type FilterParams struct {
	Fieldnames []string

	// in
	Ins []string

	// length
	Op     string
	Length int
}

var filterParams FilterParams
var filterActions = []string{"length", "null", "in"}
var supportOps = []string{"<", "<=", ">", ">=", "==", "!="}

// filterCmd represents the text command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "对文本的操作",
	Long: `对需要进行nlp处理的文本进行操作，例如编码转化，全角半角转换，字符串替换等。
For example:

    nlp-corpus-tools filter --help

支持以下功能:
- [x] length: 按字符串的长度进行过滤
- [x] null: 过滤掉空字符串的记录
- [x] in: 如果字段值等于某些值，则过滤
`,
	Example: `
1. 过滤掉字段值为空的记录:

    nlp-corpus-tools filter -a null -i test.csv --fields=content
	cat test.csv|nlp-corpus-tools filter -a null --fields=content

2. 过滤某些值的记录：

    nlp-corpus-tools filter -a in -i test.csv --fields=content --in=无 --in=沥林
	cat test.csv|nlp-corpus-tools filter -a in --fields=content --in=无 --in=沥林

3. 过滤长度小于3的记录：

    nlp-corpus-tools filter -a length -i test.csv --fields=content --length=3 --op="<"
	cat test.csv|nlp-corpus-tools filter -a length --fields=content --length=3 --op="<"

`,

	Run: func(cmd *cobra.Command, args []string) {
		rootParams.checkAction(filterActions)

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
			fmt.Printf("filter params: %+v\n", filterParams)
		}

		switch rootParams.action {
		case "length":
			if !common.StringIn(filterParams.Op, supportOps) {
				panic("op参数错误")
			}
			if filterParams.Length < 0 {
				panic("length参数错误")
			}
		case "in":
			if len(filterParams.Ins) < 1 {
				panic("in参数错误")
			}
		}

		funcMap := map[string]func(string) bool{
			"length": FilterLength,
			"null":   FilterNull,
			"in":     FilterIn,
		}
		if err = filter.Filter(r, w, filterParams.Fieldnames, funcMap[rootParams.action]); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(filterCmd)

	filterCmd.PersistentFlags().StringArrayVar(&filterParams.Fieldnames, "fields", []string{}, "指定转换的字段名，默认对全部的字段进行处理")

	filterCmd.PersistentFlags().StringArrayVar(&filterParams.Ins, "in", []string{}, "指定需要进行过滤的值")

	filterCmd.PersistentFlags().StringVar(&filterParams.Op, "op", "", "比较操作符，例如长度小于某个值的时候过滤，则可以设置op=\"<\"。支持的操作符："+strings.Join(supportOps, ", "))
	filterCmd.PersistentFlags().IntVar(&filterParams.Length, "length", 0, "比较的文本长度")
}

func FilterIn(s string) bool {
	return common.StringIn(s, filterParams.Ins)
}

func FilterNull(s string) bool {
	return s == ""
}

func FilterLength(s string) bool {
	l := len([]rune(s))
	switch filterParams.Op {
	case ">":
		return l > filterParams.Length
	case ">=":
		return l >= filterParams.Length
	case "<":
		return l < filterParams.Length
	case "<=":
		return l <= filterParams.Length
	case "==":
		return l == filterParams.Length
	case "!=":
		return l != filterParams.Length
	}
	return false
}

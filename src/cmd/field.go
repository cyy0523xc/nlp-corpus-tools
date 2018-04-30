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
	"io"
	"os"

	"github.com/cyy0523xc/nlp-corpus-tools/common"
	"github.com/cyy0523xc/nlp-corpus-tools/field"
	"github.com/spf13/cobra"
)

type FieldParams struct {
	fieldnames   []string
	oldFieldname string
	newFieldname string
}

var fieldParams FieldParams
var fieldActions = []string{"rename", "delete", "keep"}

// fieldCmd represents the field command
var fieldCmd = &cobra.Command{
	Use:   "field",
	Short: "对字段的操作",
	Long: `对字段进行操作，例如删除字段，保留字段，字段重命名等。
For example:

    nlp-corpus-tools field --help

支持的功能:
- [ ] rename: 字段重命名
- [ ] delete: 删除字段
- [x] keep: 只保留若干字段
- [ ] 
`,
	Example: `
1. 只保留若干字段

    nlp-corpus-tools field -a keep -i test.csv --fieldnames=content --fieldnames=house
	cat test.csv|nlp-corpus-tools field -a keep --fieldnames=content --fieldnames=house

	`,
	Run: func(cmd *cobra.Command, args []string) {
		rootParams.checkAction(fieldActions)

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

		switch rootParams.action {
		case "keep":
			fieldParams.keep(r, w)
		}
	},
}

func init() {
	rootCmd.AddCommand(fieldCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fieldCmd.PersistentFlags().String("foo", "", "A help for foo")
	fieldCmd.PersistentFlags().StringVar(&fieldParams.oldFieldname, "old-fieldname", "", "修改前的字段名")
	fieldCmd.PersistentFlags().StringVar(&fieldParams.newFieldname, "new-fieldname", "", "修改后的字段名")
	fieldCmd.PersistentFlags().StringArrayVar(&fieldParams.fieldnames, "fieldnames", []string{}, "字段名列表，每个fieldnames参数能指定一个字段名")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fieldCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (p FieldParams) keep(r io.Reader, w io.Writer) {
	if common.Debug {
		os.Stderr.WriteString(fmt.Sprintf("fieldnames: %+v\n", fieldParams.fieldnames))
	}
	if err := field.Keep(r, w, fieldParams.fieldnames); err != nil {
		panic(err)
	}
}

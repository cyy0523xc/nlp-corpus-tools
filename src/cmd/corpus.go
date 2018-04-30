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

	"github.com/cyy0523xc/nlp-corpus-tools/corpus"
	"github.com/spf13/cobra"
)

type CorpusParams struct {
	rate                int
	trainFile, testFile string
}

var corpusParams CorpusParams
var corpusActions = []string{"batch", "merge", "split", "count-records"}

// corpusCmd represents the corpus command
var corpusCmd = &cobra.Command{
	Use:   "corpus",
	Short: "对语料库文件的操作",
	Long: `对语料库文件进行操作，例如语料库文件的拆分，合并等。
For example:

    nlp-corpus-tools corpus --help

支持的功能如下：
- [ ] batch: 把一个大文件分拆成小的批量文件，可以方便进行标注
- [ ] merge: 例如将一个若干小的批量文件合并成大文件，例如训练集
- [x] split: 将文件按比例拆分成训练集和测试集
- [x] count-records: 统计csv文件的记录数（注意不是文件的行数）
`,
	Example: `
1. 统计语料库文件的记录数

    nlp-corpus-tools corpus -a count-records -i test.csv
	cat test.csv|nlp-corpus-tools corpus -a count-records

2. 将文件按比例拆分成训练集和测试集:

    nlp-corpus-tools corpus -a split -i test.csv --rate=0.8 --train-file=train.out --test-file=test.out
    cat test.csv|nlp-corpus-tools corpus -a split --rate=0.8 --train-file=train.out --test-file=test.out

`,
	Run: func(cmd *cobra.Command, args []string) {
		rootParams.checkAction(corpusActions)

		r, err := rootParams.getInput()
		if err != nil {
			panic(err)
		}
		defer r.Close()

		switch rootParams.action {
		case "split":
			corpusParams.split(r)
		case "count-records":
			corpusParams.countRecords(r)
		}

	},
}

func init() {
	rootCmd.AddCommand(corpusCmd)
	addActionFlag(corpusActions)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// corpusCmd.PersistentFlags().String("foo", "", "A help for foo")
	corpusCmd.PersistentFlags().IntVar(&corpusParams.rate, "rate", 5, "测试集的比例，默认该值为5，即测试集的比例为1/5")
	corpusCmd.PersistentFlags().StringVar(&corpusParams.trainFile, "train-file", "", "训练集的保存文件名")
	corpusCmd.PersistentFlags().StringVar(&corpusParams.testFile, "test-file", "", "测试集的保存文件名")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// corpusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func (p CorpusParams) countRecords(r io.Reader) {
	if n, err := corpus.CountRecords(r); err != nil {
		panic(err)
	} else {
		fmt.Printf("总记录数：%d\n", n)
	}
}

func (p CorpusParams) split(r io.Reader) {
	trainF, err := os.Create(corpusParams.trainFile)
	if err != nil {
		panic(err)
	}
	defer trainF.Close()

	testF, err := os.Create(corpusParams.testFile)
	if err != nil {
		panic(err)
	}
	defer testF.Close()

	if err = corpus.Split(r, trainF, testF, corpusParams.rate); err != nil {
		panic(err)
	}
}

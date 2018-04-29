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
	"os"

	"github.com/cyy0523xc/nlp-corpus-tools/common"
	"github.com/cyy0523xc/nlp-corpus-tools/corpus"
	"github.com/spf13/cobra"
)

type CorpusParams struct {
	rate                int
	trainFile, testFile string
}

var corpusParams CorpusParams
var corpusActions = []string{"batch", "merge", "split"}

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
- [ ] split: 将文件按比例拆分成训练集和测试集
- [ ] stat: 统计csv文件的记录数（注意不是文件的行数）
`,
	Example: `
1. 将文件按比例拆分成训练集和测试集:

    nlp-corpus-tools corpus -a split -i source_corpus.csv --rate=0.8

`,
	Run: func(cmd *cobra.Command, args []string) {
		if !common.StringIn(rootParams.action, corpusActions) {
			panic("不支持该action参数值")
		}

		r, err := os.Open(rootParams.infile)
		if err != nil {
			panic(err)
		}
		defer r.Close()

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

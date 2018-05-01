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
	"github.com/cyy0523xc/nlp-corpus-tools/corpus"
	"github.com/spf13/cobra"
)

type CorpusParams struct {
	per     int
	rate    int
	total   int
	prefix  string
	infiles []string
}

var corpusParams CorpusParams
var corpusActions = []string{"batch", "merge", "split", "stat"}

// corpusCmd represents the corpus command
var corpusCmd = &cobra.Command{
	Use:   "corpus",
	Short: "对语料库文件的操作",
	Long: `对语料库文件进行操作，例如语料库文件的拆分，合并等。
For example:

    nlp-corpus-tools corpus --help

支持的功能如下：
- [x] batch: 把一个大文件分拆成小的批量文件，可以方便进行标注。如果文件很大，还支持按比例抽取
- [x] merge: 例如将一个若干小的批量文件合并成大文件，例如训练集
- [x] split: 将文件按比例拆分成训练集和测试集
- [x] stat: 统计csv文件的记录数, 字段数，计算文件编码等
`,
	Example: `
1. 统计语料库文件的记录数

    nlp-corpus-tools corpus -a stat -i test.csv
	cat test.csv|nlp-corpus-tools corpus -a stat

2. 将文件按比例拆分成训练集和测试集，抽取1/4的样本数量作为测试集:

    nlp-corpus-tools corpus -a split -i test.csv --rate=4 --prefix=output
    cat test.csv|nlp-corpus-tools corpus -a split --rate=4 --prefix=output

3. 将一个大文件拆分成若干的小文件，只抽取原记录数的1/10，每个小文件500条记录，生成10个文件：

    nlp-corpus-tools corpus -a batch -i test.csv --rate=3 --prefix=output --per=500 --tatal=10

4. 将多个语料库文件合并成一个语料库文件

    nlp-corpus-tools corpus -a merge --infiles=test1.csv --infiles=test2.csv
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootParams.checkAction(corpusActions)

		r, err := rootParams.getInput()
		if err != nil {
			panic(err)
		}
		defer r.Close()

		switch rootParams.action {
		case "stat":
			corpusParams.stat(r)
		case "split":
			if corpusParams.rate < 1 {
				panic("rate参数错误，该值不能小于1")
			}
			if corpusParams.prefix == "" {
				panic("prefix参数错误")
			}
			corpusParams.split(r)
		case "batch":
			if corpusParams.rate < 0 {
				panic("rate参数错误")
			}
			if corpusParams.prefix == "" {
				panic("prefix参数错误")
			}
			corpusParams.batch(r)
		case "merge":
			w, err := rootParams.getOutput()
			if err != nil {
				panic(err)
			}
			if len(corpusParams.infiles) < 1 {
				panic("infiles参数错误")
			}
			defer w.Close()
			corpusParams.merge(w)
		}

	},
}

func init() {
	rootCmd.AddCommand(corpusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	corpusCmd.PersistentFlags().IntVar(&corpusParams.rate, "rate", 0, "当action=split: 测试集的比例; 当action=batch: 抽取样本的占比。例如该值为5，即比例为1/5，如果该值为0，则全部")
	corpusCmd.PersistentFlags().StringVar(&corpusParams.prefix, "prefix", "", "输出的语料库文件的文件名前缀")

	corpusCmd.PersistentFlags().IntVar(&corpusParams.per, "per", 500, "每个文件的保存记录数")
	corpusCmd.PersistentFlags().IntVar(&corpusParams.total, "total", 10, "总共保存为total个文件")

	corpusCmd.PersistentFlags().StringArrayVar(&corpusParams.infiles, "infiles", []string{}, "输入的文件名列表，可以指定多个infiles参数")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// corpusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (p CorpusParams) stat(r io.Reader) {
	if n, err := corpus.Stat(r); err != nil {
		panic(err)
	} else {
		fmt.Printf("总记录数：%d\n", n)
	}
}

func (p CorpusParams) split(r io.Reader) {
	var format = p.prefix + ".%s.csv"
	var fname = fmt.Sprintf(format, "train")
	trainF, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer trainF.Close()

	fname = fmt.Sprintf(format, "test")
	testF, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer testF.Close()

	if err = corpus.Split(r, trainF, testF, p.rate); err != nil {
		panic(err)
	}
}

func (p CorpusParams) batch(r io.Reader) {
	var err error
	var format = p.prefix + ".%03d.csv"
	var ws = make([]io.Writer, p.total)
	for i, _ := range ws {
		if f, err := os.Create(fmt.Sprintf(format, i)); err != nil {
			fmt.Printf("create file error: ", fmt.Sprintf(format, i))
			panic(err)
		} else {
			defer f.Close()
			ws[i] = io.Writer(f)
		}
	}

	if err = corpus.Batch(r, ws, p.rate, p.per); err != nil {
		panic(err)
	}
}

func (p CorpusParams) merge(w io.Writer) {
	if common.Debug {
		fmt.Printf("infiles: %+v\n", p.infiles)
	}
	var err error
	var rs = make([]io.Reader, len(p.infiles))
	for i, fname := range p.infiles {
		if f, err := os.Open(fname); err != nil {
			fmt.Printf("open file error: ", fname)
			panic(err)
		} else {
			defer f.Close()
			rs[i] = io.Reader(f)
		}
	}

	if err = corpus.Merge(rs, w); err != nil {
		panic(err)
	}
}

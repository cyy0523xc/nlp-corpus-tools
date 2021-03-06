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
	"os"

	"github.com/cyy0523xc/nlp-corpus-tools/common"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RootParams struct {
	infile  string
	outfile string
	action  string
}

var (
	rootParams RootParams
)
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nlp-corpus-tools",
	Short: ToolName,
	Long: ToolMsg + `

支持csv文件格式的语料库的字段处理，值处理，拆分合并等常用操作。

For help:
    nlp-corpus-tools --help

说明：
1. 除非有特殊说明，否则都支持语料库文件内容的管道输入和输出
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVar(&common.Debug, "debug", false, "全局测试状态")

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.src.yaml)")
	rootCmd.PersistentFlags().StringVarP(&rootParams.infile, "infile", "i", "", "输入的语料库文件")
	rootCmd.PersistentFlags().StringVarP(&rootParams.outfile, "outfile", "o", "", "输出的保存文件")

	rootCmd.PersistentFlags().IntVar(&common.Skip, "skip", 0, "读取csv时跳过前面若干行记录")
	rootCmd.PersistentFlags().IntVar(&common.Limit, "limit", 0, "读取csv时限制只读取若干行记录")

	rootCmd.PersistentFlags().StringVarP(&rootParams.action, "action", "a", "", "支持的子命令，每个命令支持的子命令是不一样的")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".src" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".src")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// getInput 获取输入的io
func (p RootParams) getInput() (*os.File, error) {
	if p.infile != "" {
		return os.Open(p.infile)
	} else {
		return os.Stdin, nil
	}
}

// getOutput 获取输出的io
func (p RootParams) getOutput() (*os.File, error) {
	if p.outfile != "" {
		return os.Create(p.outfile)
	} else {
		return os.Stdout, nil
	}
}

func (p RootParams) checkAction(actions []string) {
	if !common.StringIn(p.action, actions) {
		panic("不支持该action参数值")
	}
}

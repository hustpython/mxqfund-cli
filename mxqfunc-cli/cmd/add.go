/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mxqfund-cli/mxqfunc-cli/util"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加自选基金，可以添加多个",
	Long:  `添加自选基金，可以添加多个`,
	Run: func(cmd *cobra.Command, args []string) {
		AddKeyToJsonFile(args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addCmd.Flags().StringArray("v", []string{}, "add code")
}

func AddKeyToJsonFile(args []string) {
	before := util.SelectedFundsViper.GetStringSlice(util.SelectedFundsKey)
	for _, v := range args {
		var err bool
		if _, ok := util.FundsMap[v]; ok {
			for _, vs := range before {
				if vs == v {
					fmt.Println(v + "已经存在")
					err = true
				}
			}
			if err {
				continue
			}
			before = append(before, v)
			fmt.Println("增加" + v + "成功")
		} else {
			fmt.Println("无效的基金ID")
		}
	}
	util.SelectedFundsViper.Set(util.SelectedFundsKey, before)
	util.SelectedFundsViper.WriteConfigAs(util.SelectedFundsFile)
}

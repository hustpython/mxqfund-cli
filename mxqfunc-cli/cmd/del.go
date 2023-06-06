/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"mxqfund-cli/mxqfunc-cli/util"

	"github.com/spf13/cobra"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "删除自选基金，可以删除多个",
	Long:  `删除自选基金，可以删除多个`,
	Run: func(cmd *cobra.Command, args []string) {
		before := util.SelectedFundsViper.GetStringSlice(util.SelectedFundsKey)
		for _, v := range args {
			if _, ok := util.FundsMap[v]; ok {
				for _, vv := range args {
					for i, kk := range before {
						if vv == kk {
							before = append(before[:i], before[i+1:]...)
							fmt.Println("删除 " + vv + "成功")
						}
					}
				}
			} else {
				fmt.Println("无效的基金ID")
			}
		}
		util.SelectedFundsViper.Set(util.SelectedFundsKey, before)
		util.SelectedFundsViper.WriteConfigAs(util.SelectedFundsFile)
	},
}

func init() {
	rootCmd.AddCommand(delCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// delCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// delCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

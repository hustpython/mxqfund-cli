/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"mxqfund-cli/mxqfunc-cli/util"
	"strconv"
)

// topCmd represents the top command
var topCmd = &cobra.Command{
	Use:   "top",
	Short: "显示当日涨幅最大的基金，默认显示前10。可以通过-n控制显示的条数",
	Long:  `显示当日涨幅最大的基金，默认显示前10。可以通过-n控制显示的条数`,
	Run: func(cmd *cobra.Command, args []string) {
		var topNum = 10
		if len(args) >= 1 {
			topNum, _ = strconv.Atoi(args[0])
		}
		util.PrintTop(topNum)
	},
}

func init() {
	rootCmd.AddCommand(topCmd)
	topCmd.Flags().Int("n", 10, "控制显示的top个数")
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/guptarohit/asciigraph"
	"mxqfund-cli/mxqfunc-cli/util"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mxqfund-cli",
	Short: "可以显示自选的基金的基本信息",
	Long:  `可以显示自选的基金的基本信息`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		selectedFunds := util.SelectedFundsViper.GetStringSlice(util.SelectedFundsKey)
		for _, v := range selectedFunds {
			grd, sr := util.GetHistoryByCode(v)
			graph := asciigraph.Plot(grd, asciigraph.Precision(2))
			jz, ti := util.GetValueByCode(v)
			fmt.Println(util.FundsMap[v][2] + " [" + v + "] " + " " + ti + " 净值(%): " + jz)
			fmt.Println(sr)
			fmt.Println(graph)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

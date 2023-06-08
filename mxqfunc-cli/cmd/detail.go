/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"
	"mxqfund-cli/mxqfunc-cli/util"
)

// detailCmd represents the detail command
var detailCmd = &cobra.Command{
	Use:   "detail",
	Short: "输入基金编号，查询基金详情",
	Long:  `输入基金编号，查询基金详情`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			grd, junZhi, sr := util.GetHistoryByCode(v, util.DefaultDays)
			var asciiColor = asciigraph.Red
			if junZhi[0] < 0 {
				asciiColor = asciigraph.Green
			}
			graph := asciigraph.PlotMany([][]float64{grd, junZhi}, asciigraph.Precision(2),
				asciigraph.SeriesColors(
					asciigraph.White,
					asciiColor,
				))
			jz, ti := util.GetValueByCode(v)
			fmt.Println(util.FundsMap[v][2] + " [" + v + "] " + " " + ti + " 净值(%): " + jz)
			fmt.Println(sr)
			fmt.Println(len(junZhi), "日均线(%):", fmt.Sprintf("%.2f", junZhi[0]))
			fmt.Println(graph)
		}
	},
}

func init() {
	rootCmd.AddCommand(detailCmd)
	detailCmd.Flags().StringArray("v", []string{}, "add code")
}

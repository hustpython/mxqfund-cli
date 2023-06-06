/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"mxqfund-cli/mxqfunc-cli/util"
)

// searchkeyCmd represents the searchkey command
var searchkeyCmd = &cobra.Command{
	Use:   "search",
	Short: "通过关键词搜索基金列表",
	Long:  `通过关键词搜索基金列表`,
	Run: func(cmd *cobra.Command, args []string) {
		showTable(args[0])
	},
}

func showTable(arg string) {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "基金编号"},
			{Align: simpletable.AlignCenter, Text: "名称"},
			{Align: simpletable.AlignCenter, Text: "最新净值"},
		},
	}

	subtotal := 0
	for _, row := range util.FundsArray {
		if !strings.Contains(row[2], arg) {
			continue
		}
		jz, _ := util.GetValueByCode(row[0])
		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: row[0]},
			{Align: simpletable.AlignCenter, Text: row[2]},
			{Align: simpletable.AlignCenter, Text: jz},
		}

		table.Body.Cells = append(table.Body.Cells, r)
		subtotal += 1
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{Align: simpletable.AlignRight, Text: "合集(条)"},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", subtotal)},
		},
	}

	table.SetStyle(simpletable.StyleRounded)
	fmt.Println(table.String())
	// time.Sleep(time.Second * 10)
	// fmt.Print("\033[2J") 清除屏幕
}

func init() {
	searchkeyCmd.Flags().String("k", "", "your search key")
	//searchkeyCmd.MarkFlagRequired("k")
	rootCmd.AddCommand(searchkeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchkeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchkeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

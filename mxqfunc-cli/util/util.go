package util

import (
	"encoding/json"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var FundsArray [][]string
var FundsMap map[string][]string

type FundValue struct {
	Gszzl  string `json:"gszzl"`
	Gztime string `json:"gztime"`
}

var (
	SelectedFundsViper = viper.New()
)

const (
	SelectedFundsFile = "selectedfund.json"
	SelectedFundsKey  = "selectedFunds"
)

func init() {
	SelectedFundsViper.SetConfigFile(SelectedFundsFile)
	if err := SelectedFundsViper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	jsonFile, err := os.Open("allfunds.json")
	if err != nil {
		fmt.Println("error opening json file")
		return
	}
	defer jsonFile.Close()

	jsonByte, _ := io.ReadAll(jsonFile)
	// json处理
	err = json.Unmarshal(jsonByte, &FundsArray)
	if err != nil {
		fmt.Println("error Unmarshal json data")
		return
	}
	FundsMap = make(map[string][]string, len(FundsArray))
	for _, v := range FundsArray {
		FundsMap[v[0]] = v
	}
}

const (
	fundValueApi    = "http://fundgz.1234567.com.cn/js/"
	historyValueApi = "http://fund.eastmoney.com/f10/F10DataApi.aspx?type=lsjz&code="
	// code=110022&sdate=2018-02-22&edate=2018-03-02&per=20
)

func GetValueByCode(code string) (string, string) {
	resp, err := http.Get(fundValueApi + code + ".js?")
	if err != nil {
		return "", ""
	}

	s, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(s), "dwjz") {
		return "", ""
	}
	var t FundValue
	err = json.Unmarshal(s[8:len(s)-2], &t)
	if err != nil {
		return "", ""
	}
	f, _ := strconv.ParseFloat(t.Gszzl, 32)

	return GetColorStr(f), t.Gztime
}

func GetColorStr(f float64) string {
	var tmpShowValue string
	if f > 0 {
		tmpShowValue = fmt.Sprintf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, 0, 0, 31, fmt.Sprintf("%.2f", f), 0x1B)
	} else {
		tmpShowValue = fmt.Sprintf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, 0, 0, 32, fmt.Sprintf("%.2f", f), 0x1B)
	}
	return tmpShowValue
}

func GetHistoryByCode(code string) ([]float64, string) {
	resp, err := http.Get(historyValueApi + code)
	if err != nil {
		return []float64{}, ""
	}
	s, _ := io.ReadAll(resp.Body)
	ds := strings.Split(string(s), "<td>")
	var vf []float64

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "上涨/回调日期"},
			{Align: simpletable.AlignCenter, Text: "上涨/回调天数"},
			{Align: simpletable.AlignCenter, Text: "上涨/回调累计(%)"},
		},
	}
	last := float64(0)
	jur := float64(0)
	subtotal := float64(0)
	d := 0
	date11 := ""
	date12 := ""
	var allday = 0
	for i, dd := range ds[1 : len(ds)-2] {
		if strings.Contains(dd, "tor bold") {
			fs := strings.Split(dd, "</td>")
			values := strings.Split(fs[len(fs)-2], ">")
			value := values[len(values)-1]
			floatFund, _ := strconv.ParseFloat(strings.Trim(value, "%"), 64)
			allday++
			if len(date11) == 0 {
				date11 = fs[0]
			}
			if len(date12) == 0 {
				date12 = fs[0]
			}
			if last == 0 {
				last = floatFund
			}
			subtotal += floatFund
			vf = append(vf, floatFund)
			if (floatFund * last) > float64(0) {
				jur += floatFund
				d++
				if i == len(ds[1:len(ds)-2])-1 {
					r := []*simpletable.Cell{
						{Align: simpletable.AlignCenter, Text: date11 + "~" + fs[0]},
						{Align: simpletable.AlignCenter, Text: strconv.Itoa(d)},
						{Align: simpletable.AlignCenter, Text: GetColorStr(jur)},
					}
					table.Body.Cells = append(table.Body.Cells, r)
				}
			} else {
				r := []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: date11 + "~" + date12},
					{Align: simpletable.AlignCenter, Text: strconv.Itoa(d)},
					{Align: simpletable.AlignCenter, Text: GetColorStr(jur)},
				}
				jur = floatFund
				date11 = fs[0]
				d = 1

				if i == len(ds[1:len(ds)-2])-1 {
					r = []*simpletable.Cell{
						{Align: simpletable.AlignCenter, Text: date11 + "~" + date11},
						{Align: simpletable.AlignCenter, Text: strconv.Itoa(d)},
						{Align: simpletable.AlignCenter, Text: GetColorStr(jur)},
					}
				}
				table.Body.Cells = append(table.Body.Cells, r)
			}
			last = floatFund
			date12 = fs[0]

		}
	}
	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{Align: simpletable.AlignCenter, Text: fmt.Sprintf("近%d日涨幅:", allday)},
			{Align: simpletable.AlignCenter, Text: GetColorStr(subtotal)},
		},
	}
	table.SetStyle(simpletable.StyleRounded)

	return vf, table.String()
}

//上涨/回调日期
//上涨/回调天数
//上涨/回调累计

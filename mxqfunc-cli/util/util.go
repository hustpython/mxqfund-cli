package util

import (
	"encoding/json"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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
	DefaultDays     = "30"
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

type TopStruct struct {
	Code  string
	Name  string
	Time  string
	Value float64
}

func PrintTop(num int) {
	var tmpRes []*TopStruct
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "开始计算，预计需要耗时", len(FundsArray)/10, "秒")
	var tB = time.Now()
	for _, code := range FundsArray {
		resp, err := http.Get(fundValueApi + code[0] + ".js?")
		if err != nil {
			continue
		}

		s, _ := io.ReadAll(resp.Body)
		if !strings.Contains(string(s), "dwjz") {
			continue
		}
		var t FundValue
		err = json.Unmarshal(s[8:len(s)-2], &t)
		if err != nil {
			continue
		}
		f, _ := strconv.ParseFloat(t.Gszzl, 32)

		tmpRes = append(tmpRes, &TopStruct{
			Code:  code[0],
			Name:  code[2],
			Time:  t.Gztime,
			Value: f,
		})
	}
	sort.Slice(tmpRes, func(i, j int) bool {
		if tmpRes[i].Value > tmpRes[j].Value {
			return true
		}
		return false
	})
	fmt.Println("计算已完成，共用时:", time.Since(tB))

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "基金编号"},
			{Align: simpletable.AlignCenter, Text: "基金名称"},
			{Align: simpletable.AlignCenter, Text: "时间"},
			{Align: simpletable.AlignCenter, Text: "净值(%)"},
		},
	}

	for _, v := range tmpRes[:num] {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: v.Code},
			{Align: simpletable.AlignCenter, Text: v.Name},
			{Align: simpletable.AlignCenter, Text: v.Time},
			{Align: simpletable.AlignCenter, Text: GetColorStr(v.Value)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleRounded)
	fmt.Println(table.String())

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

func GetHistoryByCode(code string, days string) ([]float64, []float64, string) {
	resp, err := http.Get(historyValueApi + code + "&per=" + days)
	if err != nil {
		return []float64{}, []float64{}, ""
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
	var shouldAddLast bool
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
					shouldAddLast = true
				}
				table.Body.Cells = append(table.Body.Cells, r)
			}
			last = floatFund
			date12 = fs[0]

		}
	}
	if shouldAddLast {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: date11 + "~" + date11},
			{Align: simpletable.AlignCenter, Text: strconv.Itoa(d)},
			{Align: simpletable.AlignCenter, Text: GetColorStr(jur)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{Align: simpletable.AlignCenter, Text: fmt.Sprintf("近%d日涨幅:", allday)},
			{Align: simpletable.AlignCenter, Text: GetColorStr(subtotal)},
		},
	}
	table.SetStyle(simpletable.StyleRounded)
	var junZhi []float64
	junT := subtotal / float64(len(vf))
	for i := 0; i < len(vf); i++ {
		junZhi = append(junZhi, junT)
	}
	return vf, junZhi, table.String()
}

//上涨/回调日期
//上涨/回调天数
//上涨/回调累计

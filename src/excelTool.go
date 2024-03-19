package main

import (
	"WeatherKeker/Provider"
	"fmt"
	elsxlib "github.com/xuri/excelize/v2"
	"os"
	"strconv"
)

func WriteExcel(fileName, sheetName, postFix string, dataList []*Provider.WeatherOfCity) {
	xlFile, err := elsxlib.OpenFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			xlFile = elsxlib.NewFile()
		} else {
			fmt.Println(err)
			return
		}
	}
	defer func() {
		if xlFile != nil {
			if err := xlFile.Close(); err != nil {
				fmt.Println(err)
			}
		}
	}()

	const firstRowIndexStr = "1"
	const secRowIndexStr = "2"
	const cityColumnStr = "A"
	hasCreateSheet := false

	for cityIdx, cityData := range dataList {
		//第一行日期, 第二行最低温/最高温
		cityRowIndexStr := strconv.Itoa(cityIdx + 3)
		for dayIdx, day := range cityData.Days {
			minColIdxStr := getColumnIndexStr((dayIdx * 2) + 1)
			maxColIdxStr := getColumnIndexStr((dayIdx * 2) + 2)

			if !hasCreateSheet {
				if len(sheetName) == 0 {
					sheetName = day.Datetime + "更新"
				}

				sheetName = sheetName + "(" + postFix + ")"
				_, _ = xlFile.NewSheet(sheetName)

				sheetIdx, _ := xlFile.GetSheetIndex(sheetName)
				xlFile.SetActiveSheet(sheetIdx)

				hasCreateSheet = true
			}

			//设置日期头
			if cityIdx == 0 {
				minColStr := minColIdxStr + firstRowIndexStr
				maxColStr := maxColIdxStr + firstRowIndexStr

				mergeErr := xlFile.MergeCell(sheetName, minColStr, maxColStr)
				if mergeErr != nil {
					fmt.Println(mergeErr)
				}

				mergeErr = xlFile.SetCellStr(sheetName, minColStr, day.Datetime)
				mergeErr = xlFile.SetCellStr(sheetName, minColIdxStr+secRowIndexStr, "最低温")
				mergeErr = xlFile.SetCellStr(sheetName, maxColIdxStr+secRowIndexStr, "最高温")
				if mergeErr != nil {
					fmt.Println(mergeErr)
				}
			}

			//设置天气数据
			minDataColStr := minColIdxStr + cityRowIndexStr
			maxDataColStr := maxColIdxStr + cityRowIndexStr
			_ = xlFile.SetCellFloat(sheetName, minDataColStr, float64(day.TempMin), 2, 32)
			_ = xlFile.SetCellFloat(sheetName, maxDataColStr, float64(day.TempMax), 2, 32)
		}

		_ = xlFile.SetCellStr(sheetName, cityColumnStr+cityRowIndexStr, cityData.CityName)

	}

	if err := xlFile.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
}

func getColumnIndexStr(index int) string {
	if index < 0 {
		return ""
	}
	if index < 26 {
		return string(rune('A' + index))
	}
	return getColumnIndexStr(index/26-1) + getColumnIndexStr(index%26)

}

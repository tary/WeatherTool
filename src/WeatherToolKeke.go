package main

import (
	"WeatherKeker/Provider"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/eiannone/keyboard"
	"os"
	"syscall"
	"unsafe"
)

type WeatherServer struct {
	ExcelPath string `toml:"excelPath"`
	SheetName string `toml:"sheetName"`

	Provider []Provider.WeatherProvider `toml:"provider"`
}

type LocationList struct {
	LocationList []Provider.Location
}

type WeatherPerDayData struct {
	Datetime string  `toml:"datatime"`
	TempMax  float32 `toml:"tempmax"`
	TempMin  float32 `toml:"tempmin"`
}

type WeatherData struct {
	ResolvedAddress string `toml:"resolvedAddress"`

	Days []WeatherPerDayData `toml:"days"`
}

func main() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleTitle := kernel32.NewProc("SetConsoleTitleW")
	title := "小壳天气助手"
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	_, _, _ = setConsoleTitle.Call(uintptr(unsafe.Pointer(titlePtr)))

	fmt.Println("...小壳天气助手...")
	fmt.Println("")
	fmt.Println("...注意:部分数据源可能有使用次数限制...")
	fmt.Println("可以在 visualcrossing.com 和高德天气 注册个新账号,替换Config.toml中的key,获取更多次数.")
	fmt.Println("")

	var locData LocationList

	// Open the TOML file
	file, err := os.Open("Location.toml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the TOML file into the struct
	if _, err := toml.NewDecoder(file).Decode(&locData); err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	var svrConfig WeatherServer
	file2, err2 := os.Open("Config.toml")
	if err2 != nil {
		fmt.Println("Error opening file:", err2)
		return
	}
	defer file2.Close()

	// Decode the TOML file into the struct
	if _, err := toml.NewDecoder(file2).Decode(&svrConfig); err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	for _, provider := range svrConfig.Provider {
		if !provider.Enable {
			continue
		}
		fmt.Printf("当前数据源:%s\n", provider.ID)

		dataList, fetchErr := Provider.GetWeatherData(locData.LocationList, provider)
		if fetchErr != nil {
			println("执行失败", provider.ID, fetchErr.Error())
			continue
		}

		fmt.Printf("开始写入:%s\n", svrConfig.ExcelPath)
		WriteExcel(svrConfig.ExcelPath, svrConfig.SheetName, provider.PostfixName, dataList)
		fmt.Println("")
		fmt.Println("写入完成")
		fmt.Println("")
	}

	fmt.Println("")
	fmt.Println("程序执行完成, 按任意键关闭, 请查看天气数据文件:", svrConfig.ExcelPath)
	fmt.Println("")

	fmt.Println("按任意键退出")
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	_, _, err = keyboard.GetKey()
	if err != nil {
		panic(err)
	}
}

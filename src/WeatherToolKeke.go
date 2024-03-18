package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/eiannone/keyboard"
	"io"
	"net/http"
	"os"
	"syscall"
	"unsafe"
)

type WeatherServer struct {
	Key       string `toml:"key"`
	APIUrl    string `toml:"apiUrl"`
	UnitGroup string `toml:"unitGroup"`
	Content   string `toml:"content"`
	ExcelPath string `toml:"excelPath"`
	SheetName string `toml:"sheetName"`
}

type Location struct {
	Location string
}

type LocationList struct {
	LocationList []Location
}

type WeatherPerDayData struct {
	Datetime string  `toml:"datatime"`
	TempMax  float32 `toml:"tempmax"`
	TempMin  float32 `toml:"tempmin"`
	Humidity float32 `toml:"humidity"`
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
	fmt.Println("...注意:使用的是免费版天气数据服务每天只能获取1000次数据...")
	fmt.Println("可以在 https://www.visualcrossing.com/ 注册个新账号,替换Config.toml中的key,获取更多次数.")
	fmt.Println("")

	var locationList LocationList

	// Open the TOML file
	file, err := os.Open("Location.toml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the TOML file into the struct
	if _, err := toml.NewDecoder(file).Decode(&locationList); err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	var serverConfig WeatherServer
	file2, err2 := os.Open("Config.toml")
	if err2 != nil {
		fmt.Println("Error opening file:", err2)
		return
	}
	defer file2.Close()

	// Decode the TOML file into the struct
	if _, err := toml.NewDecoder(file2).Decode(&serverConfig); err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	//https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/Harbin%2CCN?unitGroup=metric&key=EB98LKGQRA5HT457J44SC4QBF&contentType=json"
	urlTemplate := fmt.Sprintf("%sCN%%%%2C%%s?unitGroup=%s&key=%s&contentType=%s",
		serverConfig.APIUrl, serverConfig.UnitGroup, serverConfig.Key, serverConfig.Content)

	dataList := make([]*WeatherData, 0, len(locationList.LocationList))

	for _, location := range locationList.LocationList {
		fullUrl := fmt.Sprintf(urlTemplate, location.Location)
		data := fetchData(fullUrl)
		if data == nil {
			return
		}

		dataList = append(dataList, data)
		fmt.Printf("成功获取%s的天气数据\n", data.ResolvedAddress)
	}

	fmt.Printf("开始写入:%s\n", serverConfig.ExcelPath)
	WriteExcel(serverConfig.ExcelPath, serverConfig.SheetName, dataList)
	fmt.Println("")
	fmt.Println("程序执行完成, 按任意键关闭, 请查看天气数据文件:", serverConfig.ExcelPath)
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

func fetchData(url string) *WeatherData {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error status code:", resp.StatusCode)
		return nil
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return nil
	}

	// Parse the JSON data
	var data = &WeatherData{}
	if err := json.Unmarshal(body, data); err != nil {
		fmt.Println("Error parsing the JSON data:", err)
		return nil
	}

	return data
}

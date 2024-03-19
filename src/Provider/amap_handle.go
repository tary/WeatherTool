package Provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
)

type amapWeatherPerDayData struct {
	Date      string `toml:"data"`
	Daytemp   string `toml:"daytemp"`
	Nighttemp string `toml:"nighttemp"`
}

type amapForecastsData struct {
	Casts []amapWeatherPerDayData `toml:"casts"`
	City  string                  `toml:"city"`
}

type amapWeatherData struct {
	Forecasts []amapForecastsData `toml:"forecasts"`
}

type amapProvider struct {
}

func (vc amapProvider) GetWeatherData(locList []Location, config WeatherProvider) ([]*WeatherOfCity, error) {

	var result = make([]*WeatherOfCity, 0)

	urlTemplate := fmt.Sprintf(config.UrlTemplate, config.APIUrl, config.Key)
	for _, loc := range locList {
		fmt.Printf("开始下载:%s的天气数据(数据源高德天气)\n", loc.Name)
		fullUrl := fmt.Sprintf(urlTemplate, loc.Name)

		// Read the response body
		body, err := getHttpBody(fullUrl)
		if err != nil || body == nil {
			fmt.Printf("获取%s的数据时时出错:%s\n", loc, err)
			return nil, err
		}

		// Parse the JSON data
		var amapData = &amapWeatherData{}
		if err := json.Unmarshal(body, amapData); err != nil {
			fmt.Printf("解析%s的数据时时出错:%s\n", loc, err)
			return nil, err
		}

		if len(amapData.Forecasts) == 0 {
			fmt.Printf("解析%s的数据时时出错:%s\n", loc, err)
			return nil, errors.New("解析的数据时时出错")
		}

		forecast := amapData.Forecasts[0]

		weatherOfCity := &WeatherOfCity{loc.Name, make([]*WeatherPerDay, 0, len(forecast.Casts))}

		for _, amapDayData := range forecast.Casts {
			dayTemp, _ := strconv.ParseFloat(amapDayData.Daytemp, 64)
			nightTemp, _ := strconv.ParseFloat(amapDayData.Nighttemp, 64)
			var data = &WeatherPerDay{amapDayData.Date,
				float32(math.Max(dayTemp, nightTemp)),
				float32(math.Min(dayTemp, nightTemp))}

			weatherOfCity.Days = append(weatherOfCity.Days, data)
		}

		result = append(result, weatherOfCity)
		println("完成一个城市天气数据下载")
	}

	return result, nil
}

package Provider

import (
	"encoding/json"
	"fmt"
)

type vcWeatherPerDayData struct {
	Datetime string  `toml:"datatime"`
	TempMax  float32 `toml:"tempmax"`
	TempMin  float32 `toml:"tempmin"`
}

type vcWeatherData struct {
	Days []vcWeatherPerDayData `toml:"days"`
}

type visualCrossingProvider struct {
}

func (vc visualCrossingProvider) GetWeatherData(locList []Location, config WeatherProvider) ([]*WeatherOfCity, error) {

	var result = make([]*WeatherOfCity, 0)

	urlTemplate := fmt.Sprintf(config.UrlTemplate, config.APIUrl, config.Key)
	for _, loc := range locList {
		fmt.Printf("开始下载:%s的天气数据(数据源VisualCrossing)\n", loc.Name)
		fullUrl := fmt.Sprintf(urlTemplate, loc.Location)

		// Read the response body
		body, err := getHttpBody(fullUrl)
		if err != nil || body == nil {
			fmt.Printf("获取%s的数据时时出错:%s\n", loc, err)
			return nil, err
		}

		// Parse the JSON data
		var vcData = &vcWeatherData{}
		if err := json.Unmarshal(body, vcData); err != nil {
			fmt.Printf("解析%s的数据时时出错:%s\n", loc, err)
			return nil, err
		}

		weatherOfCity := &WeatherOfCity{loc.Name, make([]*WeatherPerDay, 0, len(vcData.Days))}

		for _, vcDayData := range vcData.Days {
			var data = &WeatherPerDay{vcDayData.Datetime, vcDayData.TempMax, vcDayData.TempMin}
			weatherOfCity.Days = append(weatherOfCity.Days, data)
		}

		result = append(result, weatherOfCity)
		println("完成一个城市天气数据下载")
	}

	return result, nil
}

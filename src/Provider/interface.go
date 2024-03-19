package Provider

type IWeatherProvider interface {
	GetWeatherData(locList []Location, config WeatherProvider) ([]*WeatherOfCity, error)
}

func GetWeatherData(locList []Location, config WeatherProvider) ([]*WeatherOfCity, error) {
	if config.ID == "visualcrossing" {
		var vc = visualCrossingProvider{}
		return vc.GetWeatherData(locList, config)
	} else if config.ID == "amap" {
		var vc = amapProvider{}
		return vc.GetWeatherData(locList, config)
	}

	return nil, nil
}

type WeatherProvider struct {
	ID          string `toml:"id"`
	Key         string `toml:"key"`
	APIUrl      string `toml:"apiUrl"`
	UrlTemplate string `toml:"urlTemplate"`
	PostfixName string `toml:"postfixName"`
	Enable      bool   `toml:"enable"`
}

type WeatherPerDay struct {
	Datetime string
	TempMax  float32
	TempMin  float32
}

type WeatherOfCity struct {
	CityName string
	Days     []*WeatherPerDay
}

type Location struct {
	Location string
	Name     string
}

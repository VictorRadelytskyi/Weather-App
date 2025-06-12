package models

var HotColdThreshold float64 = 290.15

type WeatherResponse struct {
	Main struct {
		Temperature float64 `json:"temp"`
		FeelsLike   float64 `json:"feels_like"`
		Humidity    float64 `json:"humidity"`
	} `json:"main"`

	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`

	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`

	HotColdThreshold float64 `json:"hot_cold_threshold"`
}

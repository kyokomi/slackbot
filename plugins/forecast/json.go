package forecast

type WeatherHacks struct {
	Copyright struct {
		Image struct {
			Height int    `json:"height"`
			Link   string `json:"link"`
			Title  string `json:"title"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"image"`
		Link     string `json:"link"`
		Provider []struct {
			Link string `json:"link"`
			Name string `json:"name"`
		} `json:"provider"`
		Title string `json:"title"`
	} `json:"copyright"`
	Description struct {
		PublicTime string `json:"publicTime"`
		Text       string `json:"text"`
	} `json:"description"`
	Forecasts []struct {
		Date      string `json:"date"`
		DateLabel string `json:"dateLabel"`
		Image     struct {
			Height int    `json:"height"`
			Title  string `json:"title"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"image"`
		Telop       string `json:"telop"`
		Temperature struct {
			Max struct {
				Celsius    string `json:"celsius"`
				Fahrenheit string `json:"fahrenheit"`
			} `json:"max"`
			Min interface{} `json:"min"`
		} `json:"temperature"`
	} `json:"forecasts"`
	Link     string `json:"link"`
	Location struct {
		Area       string `json:"area"`
		City       string `json:"city"`
		Prefecture string `json:"prefecture"`
	} `json:"location"`
	PinpointLocations []struct {
		Link string `json:"link"`
		Name string `json:"name"`
	} `json:"pinpointLocations"`
	PublicTime string `json:"publicTime"`
	Title      string `json:"title"`
}

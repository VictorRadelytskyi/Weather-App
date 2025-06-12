package utils

import (
	"Weather/internal/models"
	"fmt"
	"log"
	"net/url"
)

func ConstructURL(query models.Query) (fullURL string) {
	url_base := `https://api.openweathermap.org/data/2.5/weather`
	parsedURL, err := url.Parse(url_base)
	if err != nil {
		log.Fatalf("failed parsing url base string: %v: ", err)
	}
	queryParams := parsedURL.Query()
	queryParams.Set("lat", fmt.Sprintf("%f", query.Lat))
	queryParams.Set("lon", fmt.Sprintf("%f", query.Lon))
	queryParams.Set("appid", query.Appid)

	parsedURL.RawQuery = queryParams.Encode()

	return parsedURL.String()
}

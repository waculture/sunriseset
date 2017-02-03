package main

import (
	"net/http"

	"time"

	"strconv"

	"github.com/bwolf/suncal"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// SunRiseSet は、日の出時刻と日の入り時刻を保持する
type SunRiseSet struct {
	Latitude  string `json:"latitude" xml:"latitude" form:"latitude" query:"latitude"`
	Longitude string `json:"longitude" xml:"longitude" form:"longitude" query:"longitude"`
	Sunrise   string `json:"sunrise" xml:"sunrise" form:"sunrise" query:"sunrise"`
	Sunset    string `json:"sunset" xml:"sunset" form:"sunset" query:"sunset"`
	Date      string `json:"date" xml:"date" form:"date" query:"date"`
}

func main() {
	e := echo.New()

	// CORS restricted
	// Allows requests from any `https://labstack.com` or `https://labstack.net` origin
	// wth GET, PUT, POST or DELETE method.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/sunriseset", func(c echo.Context) error {
		u := new(SunRiseSet)
		if err := c.Bind(u); err != nil {
			return err
		}
		const timeLayout = "15:04:05"
		var latitude float64
		var longitude float64
		var date time.Time
		loc, _ := time.LoadLocation("Asia/Tokyo")

		if lat, err := strconv.ParseFloat(u.Latitude, 64); err == nil {
			latitude = lat
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "緯度不正")
		}

		if lng, err := strconv.ParseFloat(u.Longitude, 64); err == nil {
			longitude = lng
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "経度不正")
		}

		coords := suncal.GeoCoordinates{latitude, longitude}
		if dt, err := time.ParseInLocation("2006-01-02", u.Date, loc); err == nil {
			date = dt
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "日付不正")
		}
		sun := suncal.SunCal(coords, date)
		u.Sunrise = sun.Rise.Format(timeLayout)
		u.Sunset = sun.Set.Format(timeLayout)
		return c.JSON(http.StatusOK, u)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

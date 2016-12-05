package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
)

const (
	AppId = ""
)

const (
	ErrApiRet = "Failed to get result from API"
)

const (
	OwmApiBaseUrl = "http://api.openweathermap.org/data/"
	OwmVersion    = "2.5"
	OwmApiUrl     = OwmApiBaseUrl + OwmVersion + "/weather"
)

type TCoord struct {
	Lon, Lat float64
}

type TWeather struct {
	Id                      int
	Main, Description, Icon string
}

type TMain struct {
	Temp               float64
	Pressure, Humidity int
	TempMin, TempMax   float64
}

type TWind struct {
	Speed, Deg int
}

type TClouds struct {
	All int
}

type TSys struct {
	Type, Id        int
	Message         float64
	Country         string
	Sunrise, Sunset int
}

type WeatherInfo struct {
	Coord TCoord

	Weather []TWeather
	Base    string
	Main    TMain

	Visibility int
	// todo: issue with wind info
	// Wind TWind

	Clouds TClouds

	Dt  int
	Sys TSys

	Id   int
	Name string
	Cod  int
}

func main() {
	flag.Parse()
	defer glog.Flush()

	ws := new(restful.WebService)
	ws.Route(ws.POST("/weather/proxy/api/q").Consumes("application/x-www-form-urlencoded").To(owmNameQuery))
	ws.Route(ws.POST("/weather/proxy/api/id").Consumes("application/x-www-form-urlencoded").To(owmIdQuery))

	restful.Add(ws)
	http.ListenAndServe(":8080", nil)
}

func encodeUrl(url string) string {
	return fmt.Sprintf("%s&appid=%s", url, AppId)
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		glog.Errorf("ERROR: %s Failed to get \"%s\"", err, url)
		return err
	} else {
		// glog.Infof("api resp body: %s\n", r.Body)
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

// todo: query city name
// api.openweathermap.org/data/2.5/weather?q={city name}
// api.openweathermap.org/data/2.5/weather?q=London

// api.openweathermap.org/data/2.5/weather?q={city name},{country code}
// api.openweathermap.org/data/2.5/weather?q=London,uk

// http://api.openweathermap.org/data/2.5/weather?q=London,uk&appid=d2f9a0e5896d7bc5a117055de7fcc46a
func owmNameQuery(req *restful.Request, resp *restful.Response) {
	city := req.QueryParameter("cityName")
	url := fmt.Sprintf("%s?q=%s", OwmApiUrl, city)
	country := req.QueryParameter("countryCode")
	if country != "" {
		url = fmt.Sprintf("%s,%s", url, country)
	}
	glog.Infof("form: %s owmNameQuery url: %s", req.Request.Form, url)

	rawInfo := WeatherInfo{}
	err := getJson(encodeUrl(url), &rawInfo)
	if err != nil {
		ret := map[string]interface{}{
			"err":      ErrApiRet,
			"ret_code": 5000,
		}
		resp.WriteAsJson(ret)
		return
	} else {
		glog.Infof("raw weather info: %s", rawInfo)
		// fmt.Println(rawInfo)
	}

	ret := map[string]interface{}{
		"ret_code": 0,
		"result":   rawInfo,
	}
	resp.WriteAsJson(ret)
}

// todo: how to query city id, may only support city name
// api.openweathermap.org/data/2.5/weather?id=2172797
func owmIdQuery(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	url := fmt.Sprintf("%s?id=%s", OwmApiUrl, id)
	glog.Infof("form: %s owmIdQuery url: %s ", req.Request.Form, url)

	rawInfo := WeatherInfo{}
	err := getJson(encodeUrl(url), &rawInfo)
	if err != nil {
		ret := map[string]interface{}{
			"err":      ErrApiRet,
			"ret_code": 5000,
		}
		resp.WriteAsJson(ret)
		return
	}

	ret := map[string]interface{}{
		"ret_code": 0,
		"result":   rawInfo,
	}
	resp.WriteAsJson(ret)
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

const DADATA_API_KEY = "ced67ee66aaf9f6df09e8e17e7ce3ffb56a05f8c"
const DADATA_SECRET_KEY = "d2ecbadfc616acaa12cbd48270e5fe685b8eb7fc"

type Address struct {
	Source       string `json:"source"`
	Result       string `json:"result"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Region       string `json:"region"`
	CityArea     string `json:"city_area"`
	CityDistrict string `json:"city_district"`
	Street       string `json:"street"`
	House        string `json:"house"`
	GeoLat       string `json:"geo_lat"`
	GeoLon       string `json:"geo_lon"`
	QcGeo        int64  `json:"qc_geo"`
}

type ResponseAddress struct {
	Addresses []*Address `json:"addresses"`
}

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

func searchAddress(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	var reqBody RequestAddressSearch
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Не указан запрос", http.StatusBadRequest)
		return
	}

	body := strings.NewReader(`[ "` + reqBody.Query + `" ]`)
	req, err := http.NewRequest("POST", "https://cleaner.dadata.ru/api/v1/clean/address", body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Token "+DADATA_API_KEY)
	req.Header.Add("X-Secret", DADATA_SECRET_KEY)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Сервис не доступен", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	var addresses []Address
	err = json.NewDecoder(resp.Body).Decode(&addresses)
	if err != nil {
		http.Error(w, "Ошибка при обработке ответа", 400)
		return
	}

	fmt.Fprintf(w, "%v", addresses)
}

func geocodeAddress(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	var reqBody RequestAddressGeocode
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	fmt.Println(reqBody)
	if err != nil || reqBody.Lat == "" || reqBody.Lng == "" {
		http.Error(w, "Не указан запрос", http.StatusBadRequest)
		return
	}

	body := strings.NewReader(`{ "lat": "` + reqBody.Lat + `", "lon": "` + reqBody.Lng + `" }`)
	req, err := http.NewRequest("POST", "http://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Token "+DADATA_API_KEY)
	req.Header.Add("X-Secret", DADATA_SECRET_KEY)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Сервис не доступен", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Ошибка при обработке ответа", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%v", string(response))
}

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	r.Post("/api/address/search", searchAddress)
	r.Post("/api/address/geocode", geocodeAddress)

	http.ListenAndServe(":8080", r)
}

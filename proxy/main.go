package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "test/docs"

	chi "github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var mx sync.RWMutex
var tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
var users = make(map[string]User)

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

// @Summary Поиск адреса
// @Description Этот эндпоинт ищет адрес через API DaData.
// @Tags address
// @Accept json
// @Produce json
// @Param request body RequestAddressSearch true "Запрос поиска адреса"
// @Success 200 {object} ResponseAddress
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Security BearerAuth
// @Router /address/search [post]
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

// @Summary Геокодинг адреса
// @Description Получение координат адреса через API DaData.
// @Tags address
// @Accept json
// @Produce json
// @Param request body RequestAddressGeocode true "Запрос на геокодинг"
// @Success 200 {object} ResponseAddress
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Security BearerAuth
// @Router /address/geocode [post]
func geocodeAddress(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	var reqBody RequestAddressGeocode
	err := json.NewDecoder(r.Body).Decode(&reqBody)
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

// @Summary Авторизация
// @Description Авторизация пользователя.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body User true "Пользователь"
// @Success 200 {string} string "Токен"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 200 {string} string "Логин или пароль неверный"
// @Router /login [post]
func login(w http.ResponseWriter, r *http.Request) {
	loginUser := User{}
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		http.Error(w, "Не указан запрос", http.StatusBadRequest)
		return
	}

	mx.Lock()
	defer mx.Unlock()
	user, ok := users[loginUser.Username]
	if !ok {
		http.Error(w, "Логин или пароль неверный", http.StatusOK)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		http.Error(w, "Логин или пароль неверный", http.StatusOK)
		return
	}

	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
		"username": user.Username,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	})
	fmt.Fprintf(w, "token: %s", tokenString)
}

// @Summary Регистрация
// @Description Регистрация пользователя.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body User true "Пользователь"
// @Success 200 {string} string "Пользователь добавлен"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /register [post]
func register(w http.ResponseWriter, r *http.Request) {
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Не указан запрос", http.StatusBadRequest)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
		return
	}

	mx.Lock()
	defer mx.Unlock()
	users[user.Username] = User{Username: user.Username, Password: string(hashPassword)}
	fmt.Fprintf(w, "add new user, username: %s", user.Username)
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	r.Post("/api/login", login)
	r.Post("/api/register", register)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
		r.Post("/api/address/search", searchAddress)
		r.Post("/api/address/geocode", geocodeAddress)
	})

	http.ListenAndServe(":8080", r)
}

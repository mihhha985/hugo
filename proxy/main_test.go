package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchAddressSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(searchAddress))
	defer ts.Close()

	query, _ := json.Marshal(RequestAddressSearch{Query: "москва сухонская 11"})
	req, err := http.Post(ts.URL+"/api/address/search", "application/json", bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	if req.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d. Got %d", http.StatusOK, req.StatusCode)
	}
}

func TestSearchAddressBadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(searchAddress))
	defer ts.Close()

	query, _ := json.Marshal(RequestAddressSearch{})
	req, err := http.Post(ts.URL+"/api/address/search", "application/json", bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	if req.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d. Got %d", http.StatusBadRequest, req.StatusCode)
	}
}

func TestGeocodeAddressSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(geocodeAddress))
	defer ts.Close()

	query, _ := json.Marshal(RequestAddressGeocode{Lat: "55.657931", Lng: "37.784606"})
	req, err := http.Post(ts.URL+"/api/address/geocode", "application/json", bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	if req.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d. Got %d", http.StatusOK, req.StatusCode)
	}
}

func TestGeocodeAddressBadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(geocodeAddress))
	defer ts.Close()

	query, _ := json.Marshal(RequestAddressGeocode{})
	req, err := http.Post(ts.URL+"/api/address/geocode", "application/json", bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	if req.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d. Got %d", http.StatusBadRequest, req.StatusCode)
	}
}

func TestRegisterUserSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(register))
	defer ts.Close()

	query, _ := json.Marshal(User{Username: "test3", Password: "test3"})
	req, err := http.Post(ts.URL+"/api/register", "application/json", bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	if req.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d. Got %d", http.StatusOK, req.StatusCode)
	}
}

func TestRegisterUserBadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(register))
	defer ts.Close()

	query, _ := json.Marshal(User{})
	req, err := http.Post(ts.URL+"/api/register", "application/json", bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	if req.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d. Got %d", http.StatusBadRequest, req.StatusCode)
	}
}

func TestLoginSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(login))
	defer ts.Close()

	query, _ := json.Marshal(User{Username: "test", Password: "test"})
	req, err := http.Post(ts.URL+"/api/login", "application/json", bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	if req.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d. Got %d", http.StatusOK, req.StatusCode)
	}
}

func TestLoginBadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(login))
	defer ts.Close()

	query, _ := json.Marshal(User{})
	req, err := http.Post(ts.URL+"/api/login", "application/json", bytes.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

	if req.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d. Got %d", http.StatusBadRequest, req.StatusCode)
	}
}

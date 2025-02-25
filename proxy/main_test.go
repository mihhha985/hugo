package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	fmt.Println(req.StatusCode)
	if req.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d. Got %d", http.StatusBadRequest, req.StatusCode)
	}
}

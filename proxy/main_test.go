package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyTestMiddleware_API(t *testing.T) {
	// Создаем тестовый запрос с путем "/api/test"
	req, err := http.NewRequest("GET", "/api/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Создаем mock-обработчик, который не будет вызываться
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	// Применяем middleware
	middleware := ProxyTest()
	handler := middleware(nextHandler)

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Проверяем тело ответа
	expected := "Proxy: /api/test\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}

func TestProxyTestMiddleware_NonAPI(t *testing.T) {
	// Создаем тестовый запрос с путем "/test"
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Создаем mock-обработчик, который не будет вызываться
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	// Применяем middleware
	middleware := ProxyTest()
	handler := middleware(nextHandler)

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем статус код (редирект)
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusFound)
	}

	// Проверяем заголовок Location для редиректа
	location := rr.Header().Get("Location")
	expectedLocation := "https://example.com"
	if location != expectedLocation {
		t.Errorf("Handler returned wrong redirect location: got %v, want %v", location, expectedLocation)
	}
}

func TestRootHandler(t *testing.T) {
	// Создаем тестовый запрос с путем "/"
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Создаем обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Проверяем тело ответа
	expected := "Hello, world!"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}

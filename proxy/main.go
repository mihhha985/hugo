package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Use(ProxyTest())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	//go WorkerTest()

	http.ListenAndServe(":8080", r)
}

const content = `123`

func WorkerTest() {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	var b byte = 0
	for range t.C {
		err := os.WriteFile("/app/static/_index.md", []byte(fmt.Sprintf("%s %d", content, b)), 0644)
		if err != nil {
			log.Println(err)
		}
		b++
	}
}

func ProxyTest() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if len(path) >= 4 && path[1:4] == "api" {
				w.Write([]byte("Hello from API"))
			} else {
				http.Redirect(w, r, "http://localhost:1313", http.StatusFound)
			}
		})
	}
}

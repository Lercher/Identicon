package main

import (
	"log"
	"net/http"
	"os"

	"github.com/lercher/identicon"
)

func getEnvOr(key string, orValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return orValue
	}
	return val
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/identicon/generate", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "no name given", http.StatusPreconditionFailed)
			return
		}

		log.Printf("generating identicon for %s", name)

		w.Header().Set("Content-Type", "image/png")
		if err := identicon.Generate([]byte(name)).WritePNGImage(w, 50, identicon.LightBackground(false)); err != nil {
			http.Error(w, "failed generating identicon", http.StatusInternalServerError)
		}
	})
	log.Printf("http://localhost:%s/identicon/generate?name=SomeName", getEnvOr("PORT", "8080"))
	if err := http.ListenAndServe(":"+getEnvOr("PORT", "8080"), mux); err != nil {
		log.Fatalf("failed listening to web server because %v", err)
	}
}

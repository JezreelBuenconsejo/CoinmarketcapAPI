package main

import (
	"coinmarketcapapi/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const apiBaseURL = "https://pro-api.coinmarketcap.com/v1"

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path[len("/api/"):]
	if endpoint == "" {
		http.Error(w, "Missing API endpoint in path", http.StatusBadRequest)
		return
	}

	targetURL, err := buildTargetURL(endpoint, r.URL.Query())
	if err != nil {
		http.Error(w, "Failed to build target URL", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-CMC_PRO_API_KEY", os.Getenv("CMC_API_KEY"))
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	copyResponse(w, resp)
}

func buildTargetURL(endpoint string, queryParams url.Values) (string, error) {
	targetURL, err := url.Parse(fmt.Sprintf("%s/%s", apiBaseURL, endpoint))
	if err != nil {
		return "", err
	}

	query := targetURL.Query()
	for key, values := range queryParams {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	targetURL.RawQuery = query.Encode()
	return targetURL.String(), nil
}

func copyResponse(w http.ResponseWriter, resp *http.Response) {
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func init() {
	if utils.IsLocalEnvironment() {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		} else {
			log.Println(".env file loaded successfully (local development)")
		}
	} else {
		log.Println("Running in production environment (skipping .env)")
	}
}

func main() {
	http.HandleFunc("/api/", proxyRequest)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

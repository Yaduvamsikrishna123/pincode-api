// main.go
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type PostOffice struct {
	Name           string `json:"Name"`
	Description    string `json:"Description"`
	BranchType     string `json:"BranchType"`
	DeliveryStatus string `json:"DeliveryStatus"`
	Circle         string `json:"Circle"`
	District       string `json:"District"`
	Division       string `json:"Division"`
	Region         string `json:"Region"`
	Block          string `json:"Block"`
	State          string `json:"State"`
	Country        string `json:"Country"`
	Pincode        string `json:"Pincode"`
}

type PostalResponse struct {
	Message    string       `json:"Message"`
	Status     string       `json:"Status"`
	PostOffice []PostOffice `json:"PostOffice"`
}

type PageData struct {
	Pincode  string
	Response *PostalResponse
	Error    string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/search", searchHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	pincode := r.FormValue("pincode")
	if pincode == "" {
		renderTemplate(w, "templates/index.html", PageData{
			Pincode: pincode,
			Error:   "Please enter a pincode",
		})
		return
	}

	url := fmt.Sprintf("https://api.postalpincode.in/pincode/%s", pincode)
	log.Printf("Making request to: %s", url)

	// Create a custom HTTP client with better configuration
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errMsg := fmt.Sprintf("Error creating request: %v", err)
		log.Println(errMsg)
		renderTemplate(w, "templates/index.html", PageData{
			Pincode: pincode,
			Error:   errMsg,
		})
		return
	}

	// Add headers to mimic a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	// Try the request with retries
	var resp *http.Response
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		reqWithCtx := req.WithContext(ctx)

		resp, err = client.Do(reqWithCtx)
		if err == nil {
			cancel()
			break
		}

		cancel()

		if i == maxRetries-1 {
			errMsg := fmt.Sprintf("Error making request after %d attempts: %v", maxRetries, err)
			log.Println(errMsg)
			renderTemplate(w, "templates/index.html", PageData{
				Pincode: pincode,
				Error:   "Unable to connect to the postal service. Please try again later.",
			})
			return
		}

		time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
	}
	if err != nil {
		errMsg := fmt.Sprintf("Error making request: %v", err)
		log.Println(errMsg)
		renderTemplate(w, "templates/index.html", PageData{
			Pincode: pincode,
			Error:   errMsg,
		})
		return
	}
	defer resp.Body.Close()

	log.Printf("Response status: %s", resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Error reading response: %v", err)
		log.Println(errMsg)
		renderTemplate(w, "templates/index.html", PageData{
			Pincode: pincode,
			Error:   errMsg,
		})
		return
	}

	log.Printf("Response body: %s", string(body))

	var responses []PostalResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		errMsg := fmt.Sprintf("Error parsing JSON: %v", err)
		log.Println(errMsg)
		renderTemplate(w, "templates/index.html", PageData{
			Pincode: pincode,
			Error:   errMsg,
		})
		return
	}

	if len(responses) == 0 || responses[0].Status != "Success" {
		errMsg := "No data found for the provided pincode"
		if len(responses) > 0 && responses[0].Message != "" {
			errMsg = responses[0].Message
		}
		log.Println(errMsg)
		renderTemplate(w, "templates/index.html", PageData{
			Pincode: pincode,
			Error:   errMsg,
		})
		return
	}

	renderTemplate(w, "templates/index.html", PageData{
		Pincode:  pincode,
		Response: &responses[0],
	})
}

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"Id"`
	OrignalURL   string    `json:"orignal-url"`
	ShortURL     string    `string:"short_url"`
	CreatingDate time.Time `json:"creatin_time"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OrignalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OrignalURL)) // It converts the orignalURL string to a byte slice
	fmt.Println("hasher: ", hasher)
    

	// give slice of data
	data := hasher.Sum(nil)
	fmt.Println("hasher data: ", data)
    
	//encode into string
	hash := hex.EncodeToString(data)
	fmt.Println("EncodeToString: ", hash) //generate a string URL 

	//We want this URL string upto 8 character
	fmt.Println("final String: ", hash[:8])
	return hash[:8]

}

func createURL(orignalURL string) string{
	shortURL := generateShortURL(orignalURL)
	id := shortURL  // Use the Short URL as the ID for Simplicity
	urlDB[id] = URL{
		ID: id,
		OrignalURL: orignalURL,
		ShortURL: shortURL,
		CreatingDate: time.Now(),
	} 
	return shortURL	
}

func getURL(id string) (URL,error) {
	url, ok := urlDB[id]

	if !ok {
		return URL{} , errors.New("URL not found")
	} 
	return url,nil 
}

func RootPageURL(w http.ResponseWriter,r *http.Request){
	fmt.Fprintf(w,"Hello World I am Akriti")
}

func ShortURLHandler(w http.ResponseWriter,r *http.Request) {
	var data struct {
		URL string 	`json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w,"Invalid request body", http.StatusBadRequest)
		return 
	}
	
	shortURL_ := createURL((data.URL))
	// fmt.Fprintf(w,shortURL)

	response := struct {
		ShortURL string `json:"short_url`
	}{ShortURL: shortURL_}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(response)

}

func redirectURLHandler(w http.ResponseWriter,r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url,err := getURL(id)
	if err != nil{
		http.Error(w,"Invalid request", http.StatusNotFound)
	}

	http.Redirect(w,r,url.OrignalURL,http.StatusFound)
	
}

func main() {
	// fmt.Println("Starting URL Shortner......")
	// OrignalURL := "https://github.com/techaakritisha";
	// generateShortURL(OrignalURL)

	//Register the handler function to handle all requests to the root URL ("/")
	http.HandleFunc("/",RootPageURL)

	//START THE HTTP SERVER ON PORT 8080
	fmt.Println("Server Starting on port 3000....")
	err := http.ListenAndServe(":3000", nil)
    if err != nil {
		fmt.Println("Error on starting server:", err)
	}

	http.HandleFunc("Shorten ", ShortURLHandler)

}

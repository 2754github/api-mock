package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type response struct {
	Status    int                    `json:"status"`
	Header    map[string]string      `json:"header"`
	Body      map[string]interface{} `json:"body"`
	NoContent bool                   `json:"no_content"`
}

type serverError struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func main() {
	setTimeZone()

	port := ":" + os.Getenv("API_PORT")
	entryPoint := os.Getenv("API_ENTRY_POINT")
	log.Printf("Server started. => http://localhost%s%s", port, entryPoint)
	http.HandleFunc(entryPoint, requestHandler)
	http.ListenAndServe(port, nil)
}

func setTimeZone() {
	loc, err := time.LoadLocation(os.Getenv("TIME_ZONE"))
	if err != nil {
		log.Print("Load TIME_ZONE failed!")
		log.Print("Set time zone to UTC and continue...")
		return
	}
	time.Local = loc
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Request received.")
	defer log.Print("Waiting for next request...")
	defer log.Print("Response returned.")

	res, errJson := loadResponseJson()
	if errJson != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(errJson)
		return
	}

	fmt.Println("[Request Detail]")
	dump, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(dump))

	for k, v := range res.Header {
		w.Header().Set(k, v)
	}
	if res.NoContent {
		return
	}
	w.WriteHeader(res.Status)
	resBody, _ := json.Marshal(res.Body)
	w.Write(resBody)
}

func loadResponseJson() (res response, errJson []byte) {
	f, err := os.Open("response.json")
	if err != nil {
		errJson, _ = json.Marshal(&serverError{
			Title:   "Open `response.json` failed!",
			Message: "Make sure that `response.json` exists.",
			Detail:  err.Error(),
		})
		return
	}
	defer f.Close()

	b, _ := ioutil.ReadAll(f)
	if err := json.Unmarshal(b, &res); err != nil {
		errJson, _ = json.Marshal(&serverError{
			Title:   "Unmarshal `response.json` failed!",
			Message: "The `response.json` setting may be incorrect.",
			Detail:  err.Error(),
		})
		return
	}
	return
}

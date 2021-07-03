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
	"unicode"
)

type response struct {
	Status  int                    `json:"status"`
	Header  map[string]string      `json:"header"`
	Body    map[string]interface{} `json:"body"`
}

type httpErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func main() {
	setTimeZone()

	port := ":" + os.Getenv("API_PORT")
	entryPoint := os.Getenv("API_ENTRY_POINT")
	log.Print("Server started. => http://localhost" + port + entryPoint)
	http.HandleFunc(entryPoint, makeHandler(requestHandler))
	http.ListenAndServe(port, nil)
}

func setTimeZone() {
	loc, err := time.LoadLocation(os.Getenv("TIME_ZONE"))
	if err != nil {
		log.Fatal("Load TIME_ZONE failed!")
	}
	time.Local = loc
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, response)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, loadResponseJson())
	}
}

func loadResponseJson() (res response) {
	f, err := os.Open("response.json")
	if err != nil {
		log.Fatal("Open response.json failed!")
	}
	defer f.Close()

	b, _ := ioutil.ReadAll(f)
	if err := json.Unmarshal(b, &res); err != nil {
		log.Fatal("Unmarshal response.json failed!")
	}
	return
}

func requestHandler(w http.ResponseWriter, r *http.Request, res response) {
	defer log.Print("Waiting for next request...")
  log.Print("Request received.")

  fmt.Println("[Request Detail]")
	dump, _ := httputil.DumpRequest(r, true)
  fmt.Print(string(dump))

	q := r.URL.Query()
	if len(q) != 0 {
		fmt.Println("[Query-Parameters]")
		i := 0
		for k, v := range q {
			i++
			fmt.Printf("%2d) %s = %s\n", i, k, v[0])
		}
	}

	for k, v := range res.Header {
		w.Header().Set(k, v)
	}

	m := r.Method
	if m != http.MethodGet {
		detail := m + " method not supported."
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(marshalJson(&httpErrorResponse{Code: 405, Message: "Method Not Allowed", Detail: detail}))
		log.Print(detail)
		return
	}

	w.WriteHeader(res.Status)
	w.Write(marshalJson(res.Body))
	log.Print("Response returned.")
}

func marshalJson(x interface{}) []byte {
	json, err := json.Marshal(x)
	if err != nil {
		log.Fatal("Marshal json failed!")
	}
	return json
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Ping struct {
	Status int
	Result string
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	// 変数pingを定義
	ping := Ping{http.StatusNotFound, "OK"}
	// 変数resにjson 出力
	res, err := json.Marshal(ping)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// 書き込みjson
	w.Write(res)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/hello", pingHandler)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":3000", nil)
}

func main() {
	handleRequests()
}

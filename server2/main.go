package main

import (
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("server2: " + r.RemoteAddr))
}
func main() {
	http.HandleFunc("/", Index)
	err := http.ListenAndServe(":8877", nil)
	log.Fatal(err)
}

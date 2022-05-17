package main

import (
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("server1"))
}
func main() {
	http.HandleFunc("/", Index)
	err := http.ListenAndServe(":7788", nil)
	log.Fatal(err)
}

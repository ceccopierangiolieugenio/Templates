// +build ignore

/* 
    this source is extracted from:
    https://golang.org/doc/articles/wiki/#tmp_2
*/

package main

import (
    "fmt"
    "log"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling: %s", r.URL.Path[1:])
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Staring WebServer on 5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

func main() {

	var port int

	flag.IntVar(&port, "p", 8880, "Port to listen on. Default port: 8880")
	flag.Parse()

	http.HandleFunc("/", reflectHandler)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

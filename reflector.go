package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	var port string

	if os.Getenv("HTTP_PLATFORM_PORT") != "" {
		port = os.Getenv("HTTP_PLATFORM_PORT")
	}

	if port == "" {
		flag.StringVar(&port, "p", "8880", "Port to listen on. Default port: 8880")
		flag.Parse()
	}

	if port == "" {
		log.Fatal("error - can not get listening port details\n")
	}

	http.HandleFunc("/", reflectHandler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

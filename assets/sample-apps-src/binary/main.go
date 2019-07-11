package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port := 8080
	if strValue, ok := os.LookupEnv("PORT"); ok {
		if intValue, err := strconv.Atoi(strValue); err == nil {
			port = intValue
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Homeport!")
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

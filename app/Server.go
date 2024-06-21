package main

import (
	"encoding/json"
	"fmt"
	"io"
	"main/handler"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	// GET method
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "test GET")
	})

	// POST method
	mux.HandleFunc("/post", func(writer http.ResponseWriter, request *http.Request) {
		data, err := io.ReadAll(request.Body)
		handler.ErrorHandler(err)

		var requestData map[string]string

		err = json.Unmarshal(data, &requestData)
		handler.ErrorHandler(err)

		fmt.Println("Recived Request: ", requestData)

		fmt.Fprint(writer, "Json berhasil diterima")
	})

	server := http.Server{
		Addr:    ":9001",
		Handler: mux,
	}

	err := server.ListenAndServeTLS("../cert.pem", "../key.pem")
	handler.ErrorHandler(err)
}

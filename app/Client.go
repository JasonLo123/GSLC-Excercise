package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/handler"
	"net/http"
	"os"
	"time"
)

func main() {
	choice := 0

	for {
		fmt.Println("Welcome")
		fmt.Println("1. Get Method")
		fmt.Println("2. Post Method")
		fmt.Println("3. Exit")
		fmt.Print(">> ")
		fmt.Scanf("%d", &choice)

		stdin := bufio.NewReader(os.Stdin)
		stdin.ReadString('\n')

		switch choice {
		case 1:
			getMethod()
		case 2:
			postMethod()
		case 3:
			return
		}
	}
}

func getMethod() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
	handler.ErrorHandler(err)

	req.Close = true

	response, err := http.DefaultClient.Do(req)
	handler.ErrorHandler(err)

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	handler.ErrorHandler(err)

	fmt.Println("Server Said: ", string(data))
}

func postMethod() {
	data := map[string]string{
		"Name": "Daniel",
	}

	JsonData, err := json.Marshal(data)
	handler.ErrorHandler(err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8080", bytes.NewBuffer(JsonData))
	handler.ErrorHandler(err)

	req.Close = true

	response, err := http.DefaultClient.Do(req)
	handler.ErrorHandler(err)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	handler.ErrorHandler(err)

	fmt.Println("Response: ", string(body))
}

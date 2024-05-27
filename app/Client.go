package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"main/handler"
	"net"
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
		fmt.Scanf("%d\n", &choice)

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

	client := createTLSClient()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://localhost:8080", nil)
	handler.ErrorHandler(err)

	req.Close = true

	response, err := client.Do(req)
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

	client := createTLSClient()
	req, err := http.NewRequestWithContext(ctx, "POST", "https://localhost:8080/post", bytes.NewBuffer(JsonData))
	handler.ErrorHandler(err)

	req.Close = true

	response, err := client.Do(req)
	handler.ErrorHandler(err)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	handler.ErrorHandler(err)

	fmt.Println("Response: ", string(body))
}

func createTLSClient() *http.Client {
	certPool := x509.NewCertPool()
	cert, err := os.ReadFile("../cert.pem")
	if err != nil {
		fmt.Println("Error reading cert.pem:", err)
		os.Exit(1)
	}
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		fmt.Println("Failed to append cert to certPool")
		os.Exit(1)
	}

	tlsConfig := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := tls.Dial(network, addr, tlsConfig)
			if err != nil {
				return nil, err
			}
			state := conn.ConnectionState()
			printTLSInfo(&state)
			return conn, nil
		},
	}

	return &http.Client{
		Transport: transport,
	}
}

func printTLSInfo(state *tls.ConnectionState) {
	if state == nil {
		fmt.Println("No TLS connection state available")
		return
	}

	fmt.Printf("TLS Version: %s\n", tlsVersionToString(state.Version))
	fmt.Printf("Cipher Suite: %s\n", tls.CipherSuiteName(state.CipherSuite))
	if len(state.PeerCertificates) > 0 {
		issuer := state.PeerCertificates[0].Issuer
		fmt.Printf("Issuer Organization: %s\n", issuer.Organization)
	}
}

func tlsVersionToString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	}
	return "Unknown"
}

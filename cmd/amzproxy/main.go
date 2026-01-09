package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

var version = "dev"

var (
	service     string
	region      string
	host        string
	showVersion bool
)

func init() {
	flag.StringVar(&service, "service", "", "AWS service name")
	flag.StringVar(&region, "region", "", "AWS region")
	flag.StringVar(&host, "host", "", "Target AWS host")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("amzproxy %s\n", version)
		os.Exit(0)
	}

	if service == "" || region == "" || host == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	signer := v4.NewSigner()

	http.HandleFunc("/", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		proxyRequest(w, r, signer, cfg.Credentials)
	}))

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func proxyRequest(w http.ResponseWriter, r *http.Request, signer *v4.Signer, creds aws.CredentialsProvider) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	targetURL := &url.URL{
		Scheme:   "https",
		Host:     host,
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
	}

	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, "failed to create proxy request", http.StatusInternalServerError)
		return
	}

	copyHeaders(r.Header, proxyReq.Header)
	proxyReq.Header.Set("connection", "close")

	cred, err := creds.Retrieve(context.Background())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve AWS credentials: %v", err), http.StatusInternalServerError)
		return
	}

	hash := sha256.Sum256(bodyBytes)
	payloadHash := hex.EncodeToString(hash[:])

	err = signer.SignHTTP(context.Background(), cred, proxyReq, payloadHash, service, region, time.Now().UTC())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to sign request: %v", err), http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultTransport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("proxy error: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeaders(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeaders(src, dst http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

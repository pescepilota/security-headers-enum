package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gosuri/uitable"
)

var securityHeaders = map[string]string{
	"Strict-Transport-Security":         "max-age=6307200; includeSubdomains",
	"X-XSS-Protection":                  "1; mode=block",
	"X-Frame-Options":                   "DENY",
	"X-Content-Type-Options":            "nosniff",
	"Content-Security-Policy":           "script-src 'self'; object-src 'self'",
	"X-Permitted-Cross-Domain-Policies": "none",
	"Referrer-Policy":                   "update needed!",
	"Expect-CT":                         "update needed!",
	"Feature-Policy":                    "camera: 'none'; payment: 'none'; microphone: 'none'",
}

var (
	flagTarget = flag.String("t", "", `Target host. Provide a URL: "https://example.com"`)
	flagHelp   = flag.Bool("h", false, "Print use instructions and exit.")
)

func main() {
	var temp int = 0

	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("HEADER NAME", "PARAMETERS", "SUGGESTIONS (if possible)")

	flag.Parse()

	if *flagHelp {
		fmt.Println("A shitty script to enum headers.")
		fmt.Println("Usage: ")
		fmt.Println("seheeh -t target URL")
		fmt.Println("")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *flagTarget == "" {
		log.Fatal("Missing target (-t https://example.com).")
	}

	resp, err := http.Get(*flagTarget)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("####  RESPONSE HEADERS #### ")
	for k, v := range resp.Header {
		fmt.Print(k)
		fmt.Print(" : ")
		fmt.Println(v)
	}

	fmt.Println("#### SECURITY HEADERS ####")
	for k, v := range securityHeaders {
		for ke, va := range resp.Header {
			if strings.EqualFold(k, ke) {
				table.AddRow(k, va, v)
				break
			} else {
				temp++
			}
			if temp > len(resp.Header) {
				temp = 0
				table.AddRow(k, "NOT SET!", v)
			}
		}
	}
	fmt.Println(table)
}

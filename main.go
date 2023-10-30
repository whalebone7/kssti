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

	
		processURL(userAgent, parsedURL)
	}
}

func processURL(userAgent string, parsedURL *url.URL) {
	client := &http.Client{}
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

	injectedValues := []string{"{{7*7}}", "<%= 7*7 %>", "${{7*7}}", "#{7*7}", "*{7*7}", "${7*7}."}

	for _, value := range injectedValues {
		if strings.Contains(string(body), "49") {
			fmt.Printf("URL: %s, Value: %s found in the response!\n", parsedURL, value)
		}
	}
}

func generateRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/52.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Safari/602.4.8",
		// Add more user agents as needed
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(userAgents))

	return userAgents[randomIndex]

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/twmb/murmur3"
)

const fix = 76

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "i", "", "File with list of URLs to hash")
	flag.Parse()

	if inputFile == "" {
		flag.PrintDefaults()
		log.Fatal("Input file not set")
	}

	urls, err := readInputFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	for _, url := range urls {
		hash := getShodanHash(url)
		fmt.Printf("%s:%d\n", url, hash)
	}
}

func getShodanHash(url string) int32 {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var buf bytes.Buffer
	var s []string

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting favicon: %v", err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading favicon content: %v", err)
	}
	str := base64.RawStdEncoding.EncodeToString(content)

	for len(str) > fix {
		s = append(s, str[:fix])
		str = str[fix:]
	}
	s = append(s, str)

	for _, ss := range s {
		buf.WriteString(ss)
		buf.WriteString("\n")
	}
	str = buf.String()

	mm3 := murmur3.StringSum32(str)

	return int32(mm3)
}

func readInputFile(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var urls []string
	for _, line := range lines {
		if line == "" {
			continue
		}
		if !strings.HasSuffix(line, "/favicon.ico") {
			line += "/favicon.ico"
		}
		urls = append(urls, line)
	}

	return urls, nil
}

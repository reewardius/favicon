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
	var url string
	flag.StringVar(&url, "u", "", "http(s)://example.com/favicon.ico")
	flag.Parse()

	if url == "" || !strings.HasSuffix(url, "/favicon.ico") {
		flag.PrintDefaults()
		log.Fatal("URL not set or does not end with /favicon.ico")
	}

	fmt.Println(getShodanHash(url))
}

// Credits @sshell_
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

	// slice up string
	for len(str) > fix {
		s = append(s, str[:fix])
		str = str[fix:]
	}
	s = append(s, str)

	// put it all together
	for _, ss := range s {
		buf.WriteString(ss)
		buf.WriteString("\n")
	}
	str = buf.String()

	// do murmurhash3 stuff
	mm3 := murmur3.StringSum32(str)

	// convert uint32 to int32
	return int32(mm3)
}

package main

import (
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

func main() {
	var url string
	flag.StringVar(&url, "u", "", "http(s)://example.com/favicon.ico")
	flag.Parse()

	if url == "" {
		flag.PrintDefaults()
		log.Fatal("URL not set")
	}
	if !strings.HasSuffix(url, "/favicon.ico") {
		log.Fatal("URL not ending with /favicon.ico")
	}

	fmt.Println(getShodanHash(url))
}

func getShodanHash(url string) int32 {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	final := ""
	fix := 76
	s := make([]string, 0)

	f, _ := http.Get(url)
	content, _ := ioutil.ReadAll(f.Body)
	str := base64.StdEncoding.EncodeToString(content)

	for i := 0; i*fix+fix < len(str); i++ {
		it := str[i*fix : i*fix+fix]
		s = append(s, it)
	}

	findlen := len(s) * fix
	last := str[findlen:] + "\n"

	for _, s := range s {
		final = final + s + "\n"
	}
	str = final + last

	mm3 := murmur3.StringSum32(str)

	return int32(mm3)
}

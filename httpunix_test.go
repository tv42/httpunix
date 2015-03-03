package httpunix_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/tv42/httpunix"
)

func Example_standalone() {
	u := &httpunix.Transport{
		DialTimeout:           100 * time.Millisecond,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	u.RegisterLocation("foo", "sock")

	var client = http.Client{
		Transport: u,
	}

	resp, err := client.Get("http+unix://foo/bar")
	if err != nil {
		log.Fatal(err)
	}
	buf, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", buf)
	resp.Body.Close()
}

func Example_integrated() {
	u := &httpunix.Transport{
		DialTimeout:           100 * time.Millisecond,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	u.RegisterLocation("foo", "sock")

	// If you want to use http: with the same client:
	t := &http.Transport{}
	t.RegisterProtocol(httpunix.Scheme, u)
	var client = http.Client{
		Transport: t,
	}

	resp, err := client.Get("http+unix://foo/bar")
	if err != nil {
		log.Fatal(err)
	}
	buf, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", buf)
	resp.Body.Close()
}

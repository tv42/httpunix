package httpunix

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"
)

const Scheme = "http+unix"

type Transport struct {
	DialTimeout           time.Duration
	RequestTimeout        time.Duration
	ResponseHeaderTimeout time.Duration

	mu sync.Mutex
	// map a URL "hostname" to a UNIX domain socket path
	loc map[string]string
}

func (t *Transport) RegisterLocation(loc string, path string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.loc == nil {
		t.loc = make(map[string]string)
	}
	if _, exists := t.loc[loc]; exists {
		panic("location " + loc + " already registered")
	}
	t.loc[loc] = path
}

var _ http.RoundTripper = (*Transport)(nil)

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL == nil {
		return nil, errors.New("http+unix: nil Request.URL")
	}
	if req.URL.Scheme != Scheme {
		return nil, errors.New("unsupported protocol scheme: " + req.URL.Scheme)
	}
	if req.URL.Host == "" {
		return nil, errors.New("http+unix: no Host in request URL")
	}
	t.mu.Lock()
	path, ok := t.loc[req.URL.Host]
	t.mu.Unlock()
	if !ok {
		return nil, errors.New("unknown location: " + req.Host)
	}

	c, err := net.DialTimeout("unix", path, t.DialTimeout)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(c)
	if t.RequestTimeout > 0 {
		c.SetWriteDeadline(time.Now().Add(t.RequestTimeout))
	}
	if err := req.Write(c); err != nil {
		return nil, err
	}
	if t.ResponseHeaderTimeout > 0 {
		c.SetReadDeadline(time.Now().Add(t.ResponseHeaderTimeout))
	}
	resp, err := http.ReadResponse(r, req)
	return resp, err
}

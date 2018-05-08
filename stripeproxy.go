package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Proxy struct {
	client *http.Client
}

func NewProxy() *Proxy {
	return &Proxy{
		client: &http.Client{
			Timeout: time.Second * 60,
		},
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request

	req, err = http.NewRequest(r.Method, "https://api.stripe.com"+r.RequestURI, r.Body)
	for name, value := range r.Header {
		req.Header.Set(name, value[0])
	}
	resp, err = p.client.Do(req)
	r.Body.Close()

	if err != nil {
		log.Printf("%v %v ERROR: %s", r.Method, r.RequestURI, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("%v %v %d", r.Method, r.RequestURI, resp.StatusCode)

	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()

}

func main() {
	proxy := NewProxy()
	listenLinda := ":46969"
	fmt.Println("Stripe proxy listening on " + listenLinda)
	err := http.ListenAndServe(listenLinda, proxy)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

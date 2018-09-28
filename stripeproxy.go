package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Proxy struct {
	client     *http.Client
	endpoint   string
	listenPort string
}

func NewProxy(endpoint, listenPort string) *Proxy {
	return &Proxy{
		client: &http.Client{
			Timeout: time.Second * 60,
		},
		endpoint:   endpoint,
		listenPort: listenPort,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request

	req, err = http.NewRequest(r.Method, p.endpoint+r.RequestURI, r.Body)
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

// ListenAndServe on listenPort
func (p *Proxy) ListenAndServe() error {
	return http.ListenAndServe(p.listenPort, p)
}

func main() {
	go func() {
		connectProxy := NewProxy("https://connect.stripe.com", ":46970")
		fmt.Println(connectProxy.endpoint + " listening on " + connectProxy.listenPort)
		err := connectProxy.ListenAndServe()
		if err != nil {
			log.Fatal("connectProxy ERROR: ", err.Error())
		}
	}()
	apiProxy := NewProxy("https://api.stripe.com", ":46969")
	fmt.Println(apiProxy.endpoint + " listening on " + apiProxy.listenPort)
	err := apiProxy.ListenAndServe()
	if err != nil {
		log.Fatal("apiProxy ERROR: ", err.Error())
	}
}

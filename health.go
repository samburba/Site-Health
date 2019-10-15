package main

import(
	"net"
	"net/http"
	"os"
	"time"
	"log"
	"strconv"
)

type Transport struct {
	rtp			http.RoundTripper
	dialer 		*net.Dialer
	connStart	time.Time
	connEnd		time.Time
	reqStart	time.Time
	reqEnd		time.Time
}

func newTransport() *Transport {
	tr := &Transport {
		dialer: &net.Dialer {
			Timeout: 	30 * time.Second,
			KeepAlive: 	30 * time.Second,
		},
	}
	tr.rtp = &http.Transport {
		Proxy: 					http.ProxyFromEnvironment,
		Dial: 					tr.dial,
		TLSHandshakeTimeout: 	10 * time.Second,
	}
	return tr
}

func (tr *Transport) RoundTrip (r *http.Request) (*http.Response, error) {
	tr.reqStart = time.Now()
	resp, err := tr.rtp.RoundTrip(r)
	tr.reqEnd = time.Now()
	return resp, err
}

func (tr *Transport) dial (network, addr string) (net.Conn, error) {
	tr.connStart = time.Now()
	cn, err := tr.dialer.Dial(network, addr)
	tr.connEnd = time.Now()
	return cn, err
}

func (tr *Transport) ReqDuration() time.Duration {
	return tr.Duration() - tr.ConnDuration()
}

func (tr *Transport) ConnDuration() time.Duration {
    return tr.connEnd.Sub(tr.connStart)
}

func (tr *Transport) Duration() time.Duration {
    return tr.reqEnd.Sub(tr.reqStart)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage health <uri>")
		os.Exit(0)
	}
	count := 1
	if len(os.Args) > 2 && len(os.Args) < 5 && os.Args[2] == "-c" {
		input_count, err := strconv.Atoi(os.Args[3])
		count = input_count;
		if err != nil {
			log.Fatalln(err)
			os.Exit(0)
		}
	}
	url := os.Args[1]
	tp := newTransport();
	client := &http.Client{Transport: tp}
	for i := 0; i < count; i++ {
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalln(err)
		os.Exit(0)
	}
		log.Printf("Response from %s: %s", url, resp.Status)
		log.Println("Duration:", tp.Duration())
		log.Println("Request duration:", tp.ReqDuration())
		log.Println("Connection duration:", tp.ConnDuration())
	}
}
package main

import(
	"net"
	"net/http"
	"os"
	"time"
	"log"
	"strconv"
	"fmt"
	"strings"
)

const WaitTime = 3

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

func getURI(uri string, client *http.Client) *http.Response {
	resp, err := client.Get(uri)
	if err != nil {
		log.Fatalln(err)
		os.Exit(0)
	}	
	return resp
}

func printStandard(uri string, resp *http.Response, tp *Transport) {
	log.Printf("Response from %s: %s", uri, resp.Status)
	log.Println("Duration:", tp.Duration())
	log.Println("Request duration:", tp.ReqDuration())
	log.Println("Connection duration:", tp.ConnDuration())
}

func printCondensed(uri string, resp *http.Response, tp *Transport) {
	fmt.Printf("\rResponse from %s: %s - Duration: %s", uri, resp.Status, tp.Duration())
}
func printGraphic(uri string, resp *http.Response, tp *Transport) {
	percent := int(tp.Duration() / 10000000)
	total := 100 - percent
	if total > 0 {
		fmt.Printf("\r%s%s\t%s\n", strings.Repeat("*", percent), strings.Repeat("-", total), tp.Duration())
	} else {
		fmt.Printf("\r%s\t%s\n", strings.Repeat("*", percent),tp.Duration())
	}
	fmt.Printf("\rResponse from %s: %s - Duration: %s", uri, resp.Status, tp.Duration())
}

func main() {
	tp := newTransport();
	client := &http.Client{Transport: tp}
	if len(os.Args) == 2 {
		resp := getURI(os.Args[1], client)
		printStandard(os.Args[1], resp, tp)
	} else if len(os.Args) == 4 && os.Args[2] == "-c" {
		count, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalln(err)
			os.Exit(0)
		}	
		for i := 0; i < count; i++ {
			resp := getURI(os.Args[1], client)
			printStandard(os.Args[1], resp, tp)
		}
	} else if len(os.Args) == 3 && os.Args[2] == "-r" {
		for {
			resp := getURI(os.Args[1], client)
			printCondensed(os.Args[1], resp, tp)
			time.Sleep(WaitTime * time.Second)
		}		
	} else if len(os.Args) == 3 && os.Args[2] == "-g" {
		for {
			resp := getURI(os.Args[1], client)
			printGraphic(os.Args[1], resp, tp)
			time.Sleep(WaitTime * time.Second)
		}		
	} else {
		log.Fatalln("Usage: health <uri> [-c <count> | -r | -g]")
		os.Exit(0)
	}
}
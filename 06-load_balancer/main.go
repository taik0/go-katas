package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type Stock struct {
	close       chan interface{}
	ProductList []struct {
		Sku      string `xml:"sku" json:"sku"`
		Quantity int    `xml:"quantity" json:"quantity"`
	} `xml:"Product" json:"products"`
}

type Subscription struct {
	data    chan []byte
	closing chan bool
}

func (s *Subscription) Close() {
	select {
	case s.closing <- true:
		close(s.closing)
	default:
		return
	}
}

type LoadBalancer interface {
	Merge(urls ...string) chan Stock
}


type RandomLoadBalancer int

func (lb *RandomLoadBalancer) Merge(urls ...string) chan Stock {
	n := make([]string, *lb)
	for i, _ := range n {
		n[i] = urls[rand.Intn(len(urls))]
	}
	return MergeSubscription(n...)
}

type RRLoadBalancer struct {
	Reqs          int
	Weights       []int
	Last          int
	CurrentWeight int
}

func (lb *RRLoadBalancer) Merge(urls ...string) chan Stock {
	n := make([]string, lb.Reqs)
	for i, _ := range n {
		n[i] = urls[lb.next(len(urls))]
	}
	return MergeSubscription(n...)
}

func (lb *RRLoadBalancer) next(numBackends int) int {
	for {
		lb.Last = (lb.Last + 1) % numBackends
		if lb.Last == 0 {
			lb.CurrentWeight = lb.CurrentWeight - lb.gcd(numBackends)
			if lb.CurrentWeight <= 0 {
				lb.CurrentWeight = lb.max(numBackends)
				if lb.CurrentWeight == 0 {
					return 0
				}
			}
		}
		if lb.Weights[lb.Last] >= lb.CurrentWeight {
			return lb.Last
		}
	}
}

func (lb *RRLoadBalancer) max(numBackends int) int {
	max := 0
	weights := lb.Weights[:numBackends]
	for i := range weights {
		if weights[i] > max {
			max = weights[i]
		}
	}
	return max
}

func (lb *RRLoadBalancer) gcd(numBackends int) int {
	weights := lb.Weights[:numBackends]

	for len(weights) > 2 {
		weights = append(weights[2:], gcdHelper(weights[0], weights[1]))
	}
	if len(weights) == 1 {
		return weights[0]
	}
	return gcdHelper(weights[0], weights[1])
}

func NewWRRLoadBalancer(reqs int, weights []int) LoadBalancer {
	return &RRLoadBalancer{Reqs: reqs, Weights: weights, Last: -1, CurrentWeight: 0}
}

func NewRRLoadBalancer(reqs int) LoadBalancer {
	weights := make([]int, 100)
	for i := range weights {
		weights[i] = 1
	}
	return &RRLoadBalancer{Reqs: reqs, Weights: weights, Last: -1, CurrentWeight: 0}
}

func gcdHelper(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func GetBackendData(url string) Subscription {
	sub := Subscription{data: make(chan []byte), closing: make(chan bool)}
	go func() {
		defer close(sub.data)
		resp, err := http.Get(url)
		defer resp.Body.Close()

		if nil != err {
			return
		}

		if resp.StatusCode != 200 {
			return
		}

		products, err := ioutil.ReadAll(resp.Body)
		if nil != err {
			return
		}

		select {
		case <-sub.closing:
		case sub.data <- products:
		}
	}()

	return sub
}

func MergeSubscription(urls ...string) chan Stock {
	results := Subscription{data: make(chan []byte), closing: make(chan bool)}
	subs := make([]Subscription, len(urls))
	for i, url := range urls {
		subs[i] = GetBackendData(url)
	}
	for _, sub := range subs {
		go func(sub Subscription) {
			data, ok := <-sub.data
			if ok == true {
				results.data <- data
				results.Close()
			}
		}(sub)
	}

	go func(subs ...Subscription) {
		<-results.closing
		for _, sub := range subs {
			sub.Close()
		}
	}(subs...)

	done := make(chan Stock)
	go func() {
		data, ok := <-results.data
		if ok == true {
			stock, err := Parse(data)
			if nil == err {
				done <- stock
			}
		}
	}()
	return done
}

const globalTimeout = 250 * time.Millisecond

func main() {

	rand.Seed(time.Now().Unix())
	port := flag.Int("port", 8080, "Listen port for the flaky backend.")
	flag.Parse()

	fmt.Printf("Port flag value: %d\n", *port)

	//var lb RandomLoadBalancer = 2
	lb := NewRRLoadBalancer(5)
	//lb := NewWRRLoadBalancer(5, []int{3,1,1,1,1,1,1})

	r := gin.Default()
	r.GET("/product", func(c *gin.Context) {

		products := lb.Merge("http://localhost:8081/product", "http://localhost:8082/product", "http://localhost:8083/product")
		select {
		case <-time.After(globalTimeout):
			c.String(500, "Timeout ")
			return
		case data, ok := <-products:
			if ok == true {
				c.JSON(200, data)
			} else {
				c.String(500, "Error from backend ")
			}
		}

	})
	r.Run(fmt.Sprintf(":%d", *port))

}

func Parse(products []byte) (Stock, error) {
	stock := Stock{}
	err := xml.Unmarshal(products, &stock)
	return stock, err
}

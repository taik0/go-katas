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
	close(s.closing)
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

	r := gin.Default()
	r.GET("/product", func(c *gin.Context) {

		products := MergeSubscription("http://localhost:8081/product", "http://localhost:8081/product", "http://localhost:8081/product")
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

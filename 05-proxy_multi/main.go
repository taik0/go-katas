package main

import (
	"encoding/json"
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
	s.closing <- true
}

func GetBackendData(url string) Subscription {
	sub := Subscription{data: make(chan []byte), closing: make(chan bool)}
	go func() {
		defer close(sub.data)
		resp, err := http.Get(url)
		if nil != err {
			fmt.Println("Cannot Get data from backend")
			return
		}

		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Println("Invalid Status Code")
			return
		}

		products, err := ioutil.ReadAll(resp.Body)
		if nil != err {
			fmt.Println("Error reading the body")
			return
		}

		select {
		case <-sub.closing:
		case sub.data <- products:
		}
	}()
	return sub
}

func MergeSubscription(urls ...string) Subscription {
	results := Subscription{data: make(chan []byte, len(urls)), closing: make(chan bool)}
	subs := make([]Subscription, len(urls))
	for i, url := range urls {
		subs[i] = GetBackendData(url)
	}
	for _, sub := range subs {
		go func(sub Subscription) {
			select {
			case data, ok := <-sub.data:
				if ok == true {
					results.data <- data
				}
			}
		}(sub)

	}
	go func() {
	select {
	case <-results.closing:
		for _, sub := range subs {
				sub.Close()
		}
	}
}()
	return results
}

const globalTimeout = 250 * time.Millisecond

func main() {

	rand.Seed(time.Now().Unix())
	port := flag.Int("port", 8080, "Listen port for the flaky backend.")
	flag.Parse()

	fmt.Printf("Port flag value: %d\n", *port)

	done := make(chan []byte)

	r := gin.Default()
	r.GET("/product", func(c *gin.Context) {

		products := MergeSubscription("http://localhost:8081/product", "http://localhost:8081/product")
		select {
		case <-time.After(globalTimeout):
			c.String(500, "Timeout ")
			return
		case xdata, ok := <-products.data:
			if ok == true {
			products.Close()
			go Parse(xdata, done)
			data := <-done
			c.String(200, string(data))
			} else {
				c.String(500, "Error from backend ")
			}
		}

	})
	r.Run(fmt.Sprintf(":%d", *port))

}

func Parse(products []byte, done chan<- []byte) {

	stock := Stock{}

	err := xml.Unmarshal(products, &stock)
	if err != nil {
		fmt.Println("Error Unmarshaling XML:", err.Error())
		return
	}

	data, err := json.Marshal(stock)
	if nil != err {
		fmt.Println("Error Marshaling to JSON:", err.Error())
		return
	}

	done <- data
}

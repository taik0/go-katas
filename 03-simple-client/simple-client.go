package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/xml"
	"encoding/json"
	"math/rand"
	"time"
)

type Stock struct {
	close       chan interface{}
	ProductList []struct {
		Sku      string `xml:"sku" json:"sku"`
		Quantity int    `xml:"quantity" json:"quantity"`
	} `xml:"Product" json:"products"`
}

func main() {

	done := make(chan interface{})
	for i := 0; i < 10; i++ {
		go Parse("http://localhost:8081/product", done)
	}

	for i := 0; i < 10; i++ {
		data := <-done
		fmt.Printf("data: %s\n", data)
	}

}

func Parse(url string, done chan interface{}) {

	stock := Stock{}
	resp, err := http.Get(url)
	if nil != err {
		fmt.Println("Error trying to get URL: ", err.Error())
		return
	}

	if resp.Status == "500" {
		fmt.Println("Error 500")
		return
	}

	products, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if nil != err {
		fmt.Println("Error reading data: ", err.Error())
		return
	}

	if len(products) == 0 {
		fmt.Println("Data len 0")
		return
	}

	time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
	err = xml.Unmarshal(products, &stock)
	if err != nil {
		fmt.Println("Error Unmarshaling XML:", err.Error())
		return
	}

	time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
	data, err := json.Marshal(stock)
	if nil != err {
		fmt.Println("Error Marshaling to JSON:", err.Error())
		return
	}

	time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
	done <- data
}

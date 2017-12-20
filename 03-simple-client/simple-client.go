package main

import "net/http"
import "fmt"
import "io/ioutil"
import "encoding/xml"

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
		fmt.Printf("data: %v\n", data)
	}

}

func Parse(url string, done chan interface{}) {
	stock := Stock{}
	resp, err := http.Get(url)
	if nil != err {
		return
	}

	products, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if nil != err {
		return
	}

	err = xml.Unmarshal(products, &stock)
	if err != nil {
		return
	}

	done <- stock
}

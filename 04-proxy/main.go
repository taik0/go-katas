package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
  "github.com/gin-gonic/gin"
)

type Stock struct {
	close       chan interface{}
	ProductList []struct {
		Sku      string `xml:"sku" json:"sku"`
		Quantity int    `xml:"quantity" json:"quantity"`
	} `xml:"Product" json:"products"`
}

func main() {

	rand.Seed(time.Now().Unix())

	done := make(chan []byte)

	r := gin.Default()
	r.GET("/product", func(c *gin.Context) {
    resp, err := http.Get("http://localhost:8081/product")
  	if nil != err {
      c.AbortWithStatus(500)
  	}

  	defer resp.Body.Close()
  	if resp.StatusCode == 500 {
  		c.AbortWithStatus(500)
  	}

  	products, err := ioutil.ReadAll(resp.Body)

  	if nil != err {
  		fmt.Println("Error reading data: ", err.Error())
  		c.AbortWithStatus(500)
  	}

  	for i := 0; i < 10; i++ {
  		go Parse(products, done)
  	}

  	data := <-done
		c.String(200, string(data))
	})
	r.Run()

}

func Parse(products []byte, done chan<- []byte) {

	stock := Stock{}

	time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
	err := xml.Unmarshal(products, &stock)
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

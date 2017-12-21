package main

import (
	"encoding/json"
	"encoding/xml"
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

func main() {

	rand.Seed(time.Now().Unix())
	port := flag.Int("port", 8081, "Listen port for the flaky backend.")
	flag.Parse()

	fmt.Printf("Port flag value: %d\n", *port)

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

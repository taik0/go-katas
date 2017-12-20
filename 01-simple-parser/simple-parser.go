package main

import (
  "encoding/json"
  "encoding/xml"
  "fmt"
  "time"
  "math/rand"
)

type Stock struct {
    close chan interface{}
    Name string
    ProductList []struct {
        Sku      string `xml:"sku" json:"sku"`
        Quantity int    `xml:"quantity" json:"quantity"`
    } `xml:"Product" json:"products"`
}

func main() {
  xmlData := []byte(`<?xml version="1.0" encoding="UTF-8" ?>
  <ProductList>
  <Product>
      <sku>ABC123</sku>
      <quantity>2</quantity>
  </Product>
  <Product>
      <sku>ABC124</sku>
      <quantity>20</quantity>
  </Product>
  </ProductList>`)

  done := make(chan interface{})
  for i := 0; i < 10; i++ {
    stock := Stock{close: done, Name: fmt.Sprintf("#%d", i)}
    go stock.Parse(xmlData)

  }

   for i := 0; i < 10; i++ {
     data := <-done
     fmt.Printf("Json %d data: %s\n", i, data)
   }
}

func (s Stock) Parse(xmlData []byte) {

  time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)

  err := xml.Unmarshal(xmlData, &s)
  if nil != err {
    fmt.Println("Error Unmarshaling XML:", err.Error())
    return
  }

  time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)

  data, err := json.Marshal(s)
  if nil != err {
    fmt.Println("Error Marshaling to JSON:", err.Error())
    return
  }

  time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
  fmt.Println(s.Name)
  s.close <- data
}

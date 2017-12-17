package main

import (
  "encoding/json"
  "encoding/xml"
  "fmt"
)

type Stock struct {
    ProductList []struct {
        Sku      string `xml:"sku" json:"sku"`
        Quantity int    `xml:"quantity" json:"quantity"`
    } `xml:"Product" json:"products"`
}

 // func (stock Stock) UnmarshalXML(b []byte) error {
 //   return "hola"
 // }

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

  data, err := parse(xmlData)
  if err != nil {
    panic(err)
  }
  //fmt.Printf("Struct data: %s\n", stock.ProductList[0].Sku)
  fmt.Printf("Json data: %s\n", data)
}

func parse(xmlData []byte) ([]byte, error) {
  stock := &Stock{}
  err := xml.Unmarshal(xmlData, stock)
  if err != nil {
    fmt.Println("Error Unmarshaling XML:", err.Error())
    return nil, err
  }
  data, err := json.Marshal(stock)
  if err != nil {
    fmt.Println("Error Marshaling to JSON:", err.Error())
    return nil, err
  }
  return data, nil
}

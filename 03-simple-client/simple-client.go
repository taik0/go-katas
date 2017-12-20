package main

import "net/http"
import "fmt"
import "io/ioutil"
import "encoding/xml"

type Stock struct {
    ProductList []struct {
        Sku      string `xml:"sku" json:"sku"`
        Quantity int    `xml:"quantity" json:"quantity"`
    } `xml:"Product" json:"products"`
}

func main() {

  stock := Stock{}
  err := Parse("http://localhost:8081/product", &stock)
  if err != nil {
    panic(err)
  }
  fmt.Println(stock.ProductList[9].Sku)

}

func Parse(url string, stock *Stock) error {
  resp, err := http.Get(url)
  if nil != err {
    return err
  }

  products, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if nil != err {
    return err
  }

  err = xml.Unmarshal(products, stock)
  if err != nil {
    return err
  }

  return nil
}

package main

import "flag"
import "fmt"
import "math/rand"
import "time"
import "encoding/xml"

type Product struct {
	Sku      string `xml:"sku"`
	Quantity int    `xml:"quantity"`
}

func main() {

	rand.Seed(time.Now().Unix())
	var sku_len int = 40
	port := flag.Int("port", 8081, "Listen port for the flaky backend.")
	flag.Parse()

	fmt.Printf("Port flag value: %d\n", *port)
	x := randStrGen(sku_len)
	fmt.Println(x)

	prod, err := randProductGen()
	if nil != err {
		panic(err)
	}
	fmt.Println(string(prod))
}

func randStrGen(lenght int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-.")
	random_array := make([]rune, lenght)
	for i, v := range rand.Perm(len(letters))[:lenght] {
		random_array[i] = rune(letters[v])
	}
	return string(random_array)
}

func randProductGen() ([]byte, error) {
	prod := Product{Sku: randStrGen(40), Quantity: rand.Intn(100)}
	return xml.MarshalIndent(prod, "", "\t")
}

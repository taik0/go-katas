package main

import "flag"
import "fmt"
import "math/rand"
import "time"
import "github.com/gin-gonic/gin"

type Product struct {
	Sku      string `xml:"sku"`
	Quantity int    `xml:"quantity"`
}

type Stock struct {
	ProductList []Product `xml:"Product"`
}

func main() {

	rand.Seed(time.Now().Unix())
	port := flag.Int("port", 8081, "Listen port for the flaky backend.")
	flag.Parse()

	fmt.Printf("Port flag value: %d\n", *port)

	r := gin.Default()
	r.GET("/product", func(c *gin.Context) {
		d := rand.Intn(100)
		if d < 10 {
			c.AbortWithStatus(500)
		} else {
			randDelay()
			c.XML(200, randProductListGen(10))
		}
	})
	r.Run(fmt.Sprintf(":%d", *port))
}

func randStrGen(lenght int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-.")
	random_array := make([]rune, lenght)
	for i, v := range rand.Perm(len(letters))[:lenght] {
		random_array[i] = rune(letters[v])
	}
	return string(random_array)
}

func randProductGen() Product {
	return Product{Sku: randStrGen(40), Quantity: rand.Intn(100)}
}

func randProductListGen(l int) Stock {
	my_stock := make([]Product, l)
	for i := 0; i < l; i++ {
		my_stock[i] = randProductGen()
	}
	return Stock{my_stock}
}

func randDelay() {
	n := rand.Intn(100)
	if n < 20 {
		time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
	} else if n < 70 {
		time.Sleep(time.Duration(rand.Int31n(50)+50) * time.Millisecond)
	} else if n < 95 {
		time.Sleep(time.Duration(rand.Int31n(500)+200) * time.Millisecond)
	}
	return
}

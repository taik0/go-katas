package main

import "flag"
import "fmt"
import "math/rand"
import "time"

func main() {

	rand.Seed(time.Now().Unix())
	var sku_len int = 40
	port := flag.Int("port", 8081, "Listen port for the flaky backend.")
	flag.Parse()

	fmt.Printf("Port flag value: %d\n", *port)
	x := randStrGen(sku_len)
	fmt.Println(x)
}

func randStrGen(lenght int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-.")
	random_array := make([]rune, lenght)
	for i, v := range rand.Perm(len(letters))[:lenght] {
		random_array[i] = rune(letters[v])
	}
	return string(random_array)
}

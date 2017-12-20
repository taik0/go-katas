package main

import "flag"
import "fmt"
import "math/rand"
import "time"


func main() {

  rand.Seed(time.Now().Unix())
  var sku_len int = 6
	port := flag.Int("port", 8081, "Listen port for the flaky backend.")
	flag.Parse()

	fmt.Printf("Port flag value: %d\n", *port)
	x := randStrGen(sku_len)
	fmt.Println(x)
}

func randStrGen(len int) string {
	random_array := make([]rune, len)
	for i, v := range rand.Perm(26)[:len] {
		random_array[i] = rune(v) + 65
	}
	return string(random_array)
}

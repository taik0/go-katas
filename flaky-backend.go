package main

import "flag"
import "fmt"

func main() {

  var port int
  flag.IntVar(&port, "port", 8081, "Listen port for the flaky backend.")
  flag.Parse()

  fmt.Printf("Port flag value: %d\n", port)
}

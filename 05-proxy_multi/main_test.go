package main

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
)

func Test(t *testing.T) {
    testServer1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, client")
    }))
    testServer2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "", 500)
    }))

    defer testServer1.Close()
    defer testServer2.Close()

    result := MergeSubscription(testServer1.URL, testServer2.URL)

    fmt.Println(result)
    fmt.Println(<-result.data)
    fmt.Println("done")
}

package main

import (
     "fmt"
     "log"
     "net/http"
)

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
        fmt.Fprintln(w, "bank-transfer-api is running")
    })

    log.Println("server starting on :8080")

    if err := http.ListenAndServe(":8080",nil); err!=nil {
        log.Fatal(err)
    }




}

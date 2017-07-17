package main

import (
    "fmt"
    "log"

    "github.com/sajari/docconv"
)

func main() {
    res, err := docconv.ConvertPath("1-s2.0-S1877050910005004-main.pdf")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(res)
}

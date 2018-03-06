package tools

import (
    "fmt"
    "os"
)

// check consolidates a common error checking mechanism into one call
func Check(err error) {
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
}


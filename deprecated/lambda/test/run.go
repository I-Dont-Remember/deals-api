package main

import (
    "fmt"
    "bufio"
    "os/exec"
)

// https://stackoverflow.com/questions/30725751/streaming-commands-output-progress
func main() {
    cmd := exec.Command("./build-update-test.sh")
    stdout, _ := cmd.StdoutPipe()
    cmd.Start()

    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        m := scanner.Text()
        fmt.Println(m)
    }

    cmd.Wait()
}

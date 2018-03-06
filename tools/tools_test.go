package tools

import (
    "testing"
    "os"
    "os/exec"
    "errors"
)

// https://stackoverflow.com/questions/26225513/how-to-test-os-exit-scenarios-in-go#33404435
func TestCheck(t *testing.T) {
    // test not failing on nil
    var nilError error
    if os.Getenv("BE_CHECK") == "1" {
        Check(nilError)
        return
    }
    cmd := exec.Command(os.Args[0], "-test.run=TestCheck")
    cmd.Env = append(os.Environ(), "BE_CHECK=1")
    err := cmd.Run()
    if e, _ := err.(*exec.ExitError); e == nil {
        return
    } else {
        t.Fatalf("want process to return nil, got %v", err)
    }

    // test failing on error
    goodError := errors.New("Test error for check")

    if os.Getenv("BE_CHECK") == "1" {
        Check(goodError)
        return
    }
    cmd = exec.Command(os.Args[0], "-test.run=TestCheck")
    cmd.Env = append(os.Environ(), "BE_CHECK=1")
    err = cmd.Run()
    if e, ok := err.(*exec.ExitError); ok && !e.Success() {
        return
    } else {
        t.Fatalf("process ran with err %v, want exit status 1", err)
    }
}

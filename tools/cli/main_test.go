package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
)


func Test_main(t *testing.T) {
    assert.NotPanics(t, main, "Main panics but shouldn't")
}

func Test_getCommands(t *testing.T) {
    retval := getCommands()
    assert.NotNil(t, retval, "Should be []cli.Command but is nil")
}

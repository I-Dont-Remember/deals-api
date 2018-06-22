package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createLocationID(t *testing.T) {
	var cases = []struct {
		name     string
		address  string
		expected string
	}{
		// To check these, use echo -n "string" | md5sum (-n necessary because
		// echo automatically adds a newline)
		{"A name", "815 Adress Ln", "0130a417ba0117e2649e3e577afb5b56"},
		{"", "", "d41d8cd98f00b204e9800998ecf8427e"},
	}

	for _, test := range cases {
		id := createLocationID(test.name, test.address)
		assert.Exactly(t, id, test.expected, "Computed incorrect checksum")
	}
}

func Test_createDealID(t *testing.T) {
	var cases = []struct {
		lid      string
		deal     string
		expected string
	}{
		// To check these, use echo -n "string" | md5sum (-n necessary because
		// echo automatically adds a newline)
		{"0130a417ba0117e2649e3e577afb5b56", "super good deal", "7c2d047d24a4dbe46d59dbd2ac359843"},
		{"d41d8cd98f00b204e9800998ecf8427e", "another great deal", "84daee89560ae8ff3430b89addb804e1"},
	}

	for _, test := range cases {
		id := createLocationID(test.lid, test.deal)
		assert.Exactly(t, id, test.expected, "Computed incorrect checksum")
	}
}

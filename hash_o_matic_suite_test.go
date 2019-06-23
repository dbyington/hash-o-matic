package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHashOMatic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HashOMatic Suite")
}

package test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSlurpCwsImages(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SlurpCwsImages Suite")
}

package e2e_tests_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2eTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2eTests Suite")
}

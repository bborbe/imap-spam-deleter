package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Imap Spam Deleter", func() {
	It("Compiles", func() {
		var err error
		_, err = gexec.Build("github.com/bborbe/imap-spam-deleter", "-mod=vendor")
		Expect(err).NotTo(HaveOccurred())
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Imap Spam Deleter Suite")
}

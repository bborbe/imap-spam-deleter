package pkg_test

import (
	"context"

	"github.com/emersion/go-imap"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bborbe/imap-spam-deleter/pkg"
)

var _ = Describe("Subject", func() {

	var ctx context.Context
	var result string
	var input string
	var err error
	BeforeEach(func() {
		ctx = context.Background()

	})
	JustBeforeEach(func() {
		msg := &imap.Message{
			Envelope: &imap.Envelope{
				Subject: input,
			},
		}
		result, err = pkg.Subject(ctx, msg)
	})
	Context("simple subject", func() {
		BeforeEach(func() {
			input = "hello world"
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns correct subject", func() {
			Expect(result).To(Equal("hello world"))
		})
	})
	Context("simple subject", func() {
		BeforeEach(func() {
			input = "=?windows-1251?B?QWxmYXpvbmUgLSBiZWZyZWllIGRlaW4gd2lsZGVzIHRpZXIh?="
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns correct subject", func() {
			Expect(result).To(Equal("Alfazone - befreie dein wildes tier!"))
		})
	})
})

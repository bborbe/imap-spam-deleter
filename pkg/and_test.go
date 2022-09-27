package pkg_test

import (
	"context"

	"github.com/emersion/go-imap"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bborbe/imap-spam-deleter/pkg"
)

var _ = Describe("And", func() {
	var rules []pkg.Rule
	var err error
	var delete bool
	var ctx context.Context
	var msg *imap.Message
	BeforeEach(func() {
		msg = &imap.Message{}
	})
	JustBeforeEach(func() {
		delete, err = pkg.And(rules...).Delete(ctx, msg)
	})
	Context("empty", func() {
		BeforeEach(func() {
			rules = []pkg.Rule{}
		})
		It("return delete", func() {
			Expect(delete).To(BeTrue())
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
	})
	Context("true", func() {
		BeforeEach(func() {
			rules = []pkg.Rule{pkg.True()}
		})
		It("return delete", func() {
			Expect(delete).To(BeTrue())
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
	})
	Context("false", func() {
		BeforeEach(func() {
			rules = []pkg.Rule{pkg.False()}
		})
		It("return delete", func() {
			Expect(delete).To(BeFalse())
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
	})
	Context("true and false", func() {
		BeforeEach(func() {
			rules = []pkg.Rule{pkg.True(), pkg.False()}
		})
		It("return delete", func() {
			Expect(delete).To(BeFalse())
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
	})
})

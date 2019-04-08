package cf_test

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/homeport/gonut/internal/gonut/cf"
)

var _ = Describe("Cloud Foundry JSON structs and contracts", func() {
	Context("Cloud Foundry API result JSON", func() {
		It("should parse Cloud Foundry API app details", func() {
			data, err := ioutil.ReadFile("../../../assets/test/cf-curl/v2/apps/nodejs-app.json")
			Expect(err).ToNot(HaveOccurred())

			var app AppDetails
			Expect(json.Unmarshal(data, &app)).ToNot(HaveOccurred())
			Expect(app.Entity.DetectedBuildpackGUID).To(BeEquivalentTo("6b70e2d7-1c63-4af9-b06d-37ae841ca8ae"))
		})

		It("should parse Cloud Foundry API buildpacks details", func() {
			data, err := ioutil.ReadFile("../../../assets/test/cf-curl/v2/buildpacks/nodejs-buildpack.json")
			Expect(err).ToNot(HaveOccurred())

			var buildpack BuildpackDetails
			Expect(json.Unmarshal(data, &buildpack)).ToNot(HaveOccurred())
			Expect(buildpack.Entity.Name).To(BeEquivalentTo("nodejs_buildpack"))
		})

		It("should parse Cloud Foundry API stacks details", func() {
			data, err := ioutil.ReadFile("../../../assets/test/cf-curl/v2/stacks/cflinuxfs3.json")
			Expect(err).ToNot(HaveOccurred())

			var stack StackDetails
			Expect(json.Unmarshal(data, &stack)).ToNot(HaveOccurred())
			Expect(stack.Entity.Name).To(BeEquivalentTo("cflinuxfs3"))
			Expect(stack.Entity.Description).To(BeEquivalentTo("Cloud Foundry Linux-based filesystem (Ubuntu 18.04)"))
		})
	})
})

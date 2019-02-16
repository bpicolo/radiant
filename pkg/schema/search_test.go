package schema_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	. "github.com/bpicolo/radiant/pkg/schema"
)

var queryWithSchema = `---
backend: main
index: shakespeare
name: shakespeare/LinesBySpeaker
schema:
  type: object
  properties:
    speaker:
      type: string
      minLength: 2
      maxLength: 5
    query:
      type: string
  required:
    - speaker
source: |
  bool:
    filter:
      term:
        speaker: "{{ .speaker }}"
  {{ if .query}}
    must:
      terms:
        text_entry: {{ .query | toJson}}
  {{ end }}
`

var queryWithoutSchema = `---
backend: main
index: shakespeare
name: shakespeare/LinesBySpeaker
source: |
  bool:
    filter:
      term:
        speaker: "{{ .speaker }}"
  {{ if .query}}
    must:
      terms:
        text_entry: {{ .query | toJson}}
  {{ end }}
`

var _ = Describe("Search", func() {
	var search Search
	var queryDefinition QueryDefinition

	Context("with an empty query schema", func() {
		BeforeEach(func() {
			err := yaml.Unmarshal([]byte(queryWithoutSchema), &queryDefinition)
			if err != nil {
				panic(err)
			}
			search = Search{
				QueryDefinition: &queryDefinition,
				Context: map[string]interface{}{
					"name": "dave",
				},
				From: 0,
				Size: 20,
			}
		})

		It("should return no validation errors", func() {
			Expect(search.Validate()).To(BeNil())
		})
	})

	Context("when validated with a query schema", func() {
		BeforeEach(func() {
			err := yaml.Unmarshal([]byte(queryWithSchema), &queryDefinition)
			if err != nil {
				panic(err)
			}
			search = Search{
				QueryDefinition: &queryDefinition,
				Context: map[string]interface{}{
					"speaker": "dave",
				},
				From: 0,
				Size: 20,
			}
		})

		It("should return nil when the search is valid", func() {
			Expect(search.Validate()).To(BeNil())
		})

		It("should return an error when missing a required field", func() {
			search.Context = map[string]interface{}{}
			Expect(search.Validate()).To(HaveOccurred())
		})

		It("should return an error when a field is invalid", func() {
			search.Context = map[string]interface{}{
				"speaker": "too long",
			}
			Expect(search.Validate()).To(HaveOccurred())
		})
	})

})

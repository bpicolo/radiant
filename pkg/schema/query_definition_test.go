package schema_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	. "github.com/bpicolo/radiant/pkg/schema"
)

var queryWithAGarbageSchema = `---
backend: main
index: shakespeare
name: shakespeare/LinesBySpeaker
schema: garbage
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

var _ = Describe("QueryDefinition", func() {
	Context("when unmarshalled with invalid query schema", func() {

		It("returns an error", func() {
			var qd QueryDefinition
			Expect(yaml.Unmarshal([]byte(queryWithAGarbageSchema), &qd)).To(HaveOccurred())
		})
	})
})

package schema

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// Search represents a search using a particular query definition
type Search struct {
	QueryDefinition *QueryDefinition
	Context         interface{}
	From            int
	Size            int
}

func (s *Search) Validate() error {
	schema := s.QueryDefinition.Schema()
	if schema == nil {
		return nil
	}
	dataLoader := gojsonschema.NewGoLoader(s.Context)
	result, err := schema.Validate(dataLoader)

	if err != nil {
		return fmt.Errorf("Failed to validate search context: %s", err)
	}
	if !result.Valid() {
		errs := make([]string, len(result.Errors()))
		for i, jerr := range result.Errors() {
			errs[i] = jerr.String()
		}
		return fmt.Errorf("Invalid query context: %s", strings.Join(errs, ", "))
	}
	return nil
}

package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bpicolo/radiant/pkg/query"
	"github.com/bpicolo/radiant/pkg/schema"
	"github.com/spf13/cobra"
)

// generateSearchCmd represents the generateSearch command
var queryCheck = &cobra.Command{
	Use:   "query-check [/path/to/query] [context]",
	Short: "Check the template output of a given query`",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		engine := query.NewEngine()
		qry, err := readSearchDefinition(args[0])
		if err != nil {
			log.Fatalf("Error reading query definition: %s", err)
		}
		ctx := make(map[string]interface{})
		err = json.Unmarshal([]byte(args[1]), &ctx)
		if err != nil {
			log.Fatalf("Error parsing context: %s", err)
		}
		search, err := engine.Interpret(&schema.Search{QueryDefinition: qry, Context: ctx})
		if err != nil {
			log.Fatalf("Error templating query: %s", err)
		}
		fmt.Println(search.ESQuery)
	},
}

func init() {
	rootCmd.AddCommand(queryCheck)
}

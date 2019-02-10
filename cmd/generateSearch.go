package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateSearchCmd represents the generateSearch command
var generateSearchCmd = &cobra.Command{
	Use:   "generateSearch [name]",
	Short: "Generate a search plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Println("generateSearch called for name %s", name)
	},
}

func init() {
	rootCmd.AddCommand(generateSearchCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateSearchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	"github.com/bpicolo/radiant/pkg/config"
	"github.com/bpicolo/radiant/pkg/server"
)

// generateSearchCmd represents the generateSearch command
var serveCmd = &cobra.Command{
	Use:  "serve [configPath]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath := args[0]
		if cfgPath == "" {
			log.Fatalf("No config file was found")
			os.Exit(1)
		}
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			log.Fatalf("The specified radiant config file was not found")
			os.Exit(1)
		}
		cfg, err := readConfig(cfgPath)
		if err != nil {
			log.Fatalf("Failed to parse radiant configuration: %s", err)
			os.Exit(1)
		}

		s, err := server.NewServer(cfg)
		if err != nil {
			log.Fatalf("Failed to start radiant server: %s", err)
		}
		bind, _ := cmd.Flags().GetString("bind")
		s.Serve(bind)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("bind", "b", "127.0.0.1:5000", "Host to bind radiant server to")
}

func readConfig(cfgPath string) (*config.RadiantConfig, error) {
	cfg := config.RadiantConfig{}
	dat, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dat, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

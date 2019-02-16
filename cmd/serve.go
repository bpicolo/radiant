package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	"github.com/bpicolo/radiant/pkg/config"
	"github.com/bpicolo/radiant/pkg/schema"
	"github.com/bpicolo/radiant/pkg/server"
)

// generateSearchCmd represents the generateSearch command
var serveCmd = &cobra.Command{
	Use:  "serve [configPath]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath := args[0]

		cfg, err := readConfig(cfgPath)
		if err != nil {
			log.Fatalf("Failed to parse radiant configuration: %s", err)
			os.Exit(1)
		}

		searchDir, _ := cmd.Flags().GetString("search_dir")
		configureSearches(cfg, searchDir)

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
	serveCmd.Flags().StringP("search_dir", "d", "./searches", "Directory containing radiant search definitions")
}

func configureSearches(cfg *config.RadiantConfig, searchDir string) {
	err := filepath.Walk(
		searchDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if err != nil {
				log.Println("Error walking search directory", err)
				return nil
			}
			qry, err := readSearchDefinition(path)
			if err != nil {
				log.Printf("Error binding search at path: %s\n%s", path, err)
				return nil
			}
			cfg.Queries = append(cfg.Queries, qry)
			return nil
		},
	)
	if err != nil {
		log.Println("Error walking search directory", err)
	}
}

func readSearchDefinition(path string) (*schema.QueryDefinition, error) {
	search := schema.QueryDefinition{}
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dat, &search)
	if err != nil {
		return nil, err
	}

	return &search, nil
}

func readConfig(cfgPath string) (*config.RadiantConfig, error) {
	if cfgPath == "" {
		log.Fatalf("No config file was specified")
	}
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("The specified radiant config file was not found")
		os.Exit(1)
	}
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

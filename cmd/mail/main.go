package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/kataras/mail"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Version is the current semantic version of the `mail` package.
const Version = mail.Version

var rootCmd = &cobra.Command{
	Use:                        "mail [command] [flags]",
	Example:                    "",
	Short:                      "mail is the command line tool for the mail package, which provides a simple way to send mail through terminal.",
	Version:                    Version,
	SilenceErrors:              true,
	TraverseChildren:           true,
	SuggestionsMinimumDistance: 1,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// unmarshal will try to check if a flag value is a file
// if so, then it will parse its contents, decode them and set to the `outPtr`,
// otherwise it will decode the flagvalue using the json unmarshaler and send the result to the `outPtr`.
func unmarshal(flagValue string, outPtr interface{}) error {
	// remove if @ (unix users usually do that).
	if flagValue[0] == '@' {
		flagValue = flagValue[1:]
	}

	ext := filepath.Ext(flagValue)
	// no file.
	if ext == "" {
		return json.Unmarshal([]byte(flagValue), outPtr)
	}

	// read file contents.
	contents, err := ioutil.ReadFile(flagValue)
	if err != nil {
		return err
	}

	// unmarshal to the outPtr based on the file extension.
	switch ext {
	case ".yml", ".yaml":
		return yaml.Unmarshal(contents, outPtr)
	case ".tml", ".toml":
		return toml.Unmarshal(contents, outPtr)
	case ".xml":
		return xml.Unmarshal(contents, outPtr)
	case ".json":
		return json.Unmarshal(contents, outPtr)
	}

	return fmt.Errorf("unsuported file extension: available formats are json, xml, yaml and toml")
}

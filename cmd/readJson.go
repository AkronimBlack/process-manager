package cmd

import (
	"fmt"
	"github.com/AkronimBlack/process-manager/pkg/parser"
	"github.com/AkronimBlack/process-manager/shared"
	"github.com/spf13/cobra"
	"os"
)

var (
	fileLocation string
)

// readJsonCmd represents the readJson command
var readJsonCmd = &cobra.Command{
	Use:   "read:parser",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		parser := parser.Parser{}
		err := parser.LoadFile(fileLocation)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		validationErrors := parser.Validate()
		if len(validationErrors) != 0 {
			fmt.Printf("invalid actions \n%s\n", shared.ToJsonPrettyString(validationErrors))
			os.Exit(1)
		}
		fmt.Println(shared.ToJsonPrettyString(parser.Actions()))
	},
}

func init() {
	rootCmd.AddCommand(readJsonCmd)
	readJsonCmd.Flags().StringVarP(&fileLocation, "file-location", "f", "", "location of parser file to parse")
}

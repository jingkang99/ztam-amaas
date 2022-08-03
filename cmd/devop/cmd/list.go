package cmd

import (
	"os"
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List project files and modules",
	Long:  `Search files or modules with a pattern`,
	Args:  cobra.OnlyValidArgs,
	ValidArgs: []string{"file", "module"},
	PreRunE: func(cmd *cobra.Command, args []string) error {
        if len(args) == 0 {
            cmd.Help()
            os.Exit(0)
        }
        return nil
    },
	Run: func(cmd *cobra.Command, args []string) {
		if(Debug) {
			fmt.Println("list called")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

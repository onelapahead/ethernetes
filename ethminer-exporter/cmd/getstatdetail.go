package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// getstatdetailCmd represents the getstatdetail command
var getstatdetailCmd = &cobra.Command{
	Use:   "getstatdetail",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.TODO()
		initErr := initClient(ctx, cmd.Parent())
		if initErr != nil {
			return initErr
		}

		stats, statErr := ethapi.GetStatDetail(ctx)
		if statErr != nil {
			return statErr
		}

		statsBytes, marshalErr := json.MarshalIndent(stats, "", "  ")
		if marshalErr != nil {
			return marshalErr
		}
		fmt.Println(string(statsBytes))

		return nil
	},
}

func init() {
	clientCmd.AddCommand(getstatdetailCmd)
}

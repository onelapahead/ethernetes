package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
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

		pong, pingErr := ethapi.Ping(ctx)
		if pingErr != nil {
			return pingErr
		}
		fmt.Printf("%s!\n", pong)

		return nil
	},
}

func init() {
	clientCmd.AddCommand(pingCmd)
}

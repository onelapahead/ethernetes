package cmd

import (
	"context"

	"github.com/hfuss/ethernetes/ethminer-exporter/pkg/ethminer/client"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var ethapi *client.ApiClient

func init() {
	rootCmd.AddCommand(clientCmd)
	networkedFlags(clientCmd)
}

func networkedFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("hostname", "127.0.0.1", "ethminer API hostname")
	cmd.PersistentFlags().Int("port", 3333, "ethminer API port")
}

func getNetworkedFlags(cmd *cobra.Command) (string, int, error) {
	hostname, hostnameErr := cmd.PersistentFlags().GetString("hostname")
	if hostnameErr != nil {
		return "", 0, hostnameErr
	}

	port, portErr := cmd.PersistentFlags().GetInt("port")
	if portErr != nil {
		return "", 0, portErr
	}

	return hostname, port, nil
}

func initClient(ctx context.Context, cmd *cobra.Command) error {
	hostname, port, flagsErr := getNetworkedFlags(cmd)
	if flagsErr != nil {
		return flagsErr
	}

	ethapi = &client.ApiClient{
		Host:         hostname,
		Port:         port,
		ConnPoolSize: 1,
	}
	ethapi.Init(ctx)
	return nil
}

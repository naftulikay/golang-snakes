package cmd

import (
	"github.com/naftulikay/golang-snakes/cmd/users"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCommand = &cobra.Command{
		Use: "golang-snakes",
		Long: "A viper/cobra demo.",
	}
)

func Execute() error {
	return rootCommand.Execute()
}

func init() {
	rootCommand.AddCommand(users.Commands()...)

	// setup env variables
	viper.AutomaticEnv()
}
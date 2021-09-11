package users

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	createUserCommand = &cobra.Command{
		Use: "create",
		Short: "Create a user.",
		PreRun: func(cmd *cobra.Command, args []string) {
			// NOTE the binding of viper to flags needs to happen here otherwise the bindings are global and overwrite each other
			const EnvKey = "api_key"
			const CliKey = "api-key"

			// bind flags to env variables
			if err := viper.BindPFlag(EnvKey, cmd.Flags().Lookup(CliKey)); err != nil {
				log.Fatalf("Unable to bind flag: %s", err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			type Config struct {
				APIKey string `mapstructure:"api_key"`
			}

			var config Config

			if err := viper.Unmarshal(&config); err != nil {
				log.Fatalf("Unable to unmarshal viper config: %s", err)
			}

			log.Printf("Create User: %+v", config)
		},
	}
)

func init() {
	const EnvKey = "api_key"
	const CliKey = "api-key"

	// bind environment variables
	if err := viper.BindEnv(EnvKey); err != nil {
		log.Fatalf("Unable to bind environment: %s", err)
	}

	// setup flags
	flags := createUserCommand.Flags()
	flags.StringP(CliKey, "k", "", "The API key.")
}
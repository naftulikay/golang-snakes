package users

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	listUsersCommand = &cobra.Command{
		Use: "list",
		Short: "List users.",
		PreRun: func(cmd *cobra.Command, args []string) {
			const EnvKey = "api_key"
			const CliKey = "api-key"

			// bind flags to env variables
			// NOTE the binding of viper to flags MUST happen here otherwise the bindings are global and overwrite each other
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
				log.Fatalf("Unable to unmarshal config: %s", err)
			}

			log.Printf("Listing users: %+v", config)
		},
	}
)

func init() {
	const EnvKey = "api_key"
	const CliKey = "api-key"

	// bind environment variables
	// NOTE this can occur here, because this is essentially just telling viper "hey you should know about env variable X"
	if err := viper.BindEnv(EnvKey); err != nil {
		log.Fatalf("Unable to bind environment: %s", err)
	}

	// setup flags
	// NOTE this can occur here, because flags are specific to each command
	flags := listUsersCommand.Flags()
	flags.StringP(CliKey, "k", "", "The API key.")
}

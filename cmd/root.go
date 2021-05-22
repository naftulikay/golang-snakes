package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// this struct will hold all of our CLI/environment/config file variables
type config struct {
	Env           string `mapstructure:"env"`
	MySQLHost     string `mapstructure:"mysql_host"`
	MySQLPort     uint16 `mapstructure:"mysql_port"`
	MySQLUser     string `mapstructure:"mysql_user"`
	MySQLPassword string `mapstructure:"mysql_password"`
	MySQLDatabase string `mapstructure:"mysql_database"`
	JWTKey        string `mapstructure:"jwt_key"`
}

// Config a global variable which will store the application config once parsed
var Config config

const DefaultEnv = "dev"

var (
	rootCmd = &cobra.Command{
		Use:   "golang-snakes", // name of the binary for usage docs
		Short: "This is the short description.",
		Long:  "This is the long description.",
		Run: func(cmd *cobra.Command, args []string) {
			// NOTE important to load/unmarshal viper here, if you do it in init(), it won't work.

			// attempt to load from disk
			if err := viper.ReadInConfig(); err != nil {
				fmt.Printf("INFO: Unable to load config file: %s\n", err)
			}

			// attempt to unmarshal config into the Config global variable
			if err := viper.Unmarshal(&Config); err != nil {
				// crash here if we failed
				log.Fatalf("ERROR: Unable to unmarshal config: %s\n", err)
			}

			// simply print out the parsed configuration for verification
			fmt.Printf("Parsed: %+v\n", &Config)
		},
	}
)

// Execute call this from your main() function to parse the flags and kick off the application, which will execute the
// Run function defined in our root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// just as a golang reminder, the init function in a module is called when the module is imported. consider it
	// similar to the main global scope in Python outside of functions and classes. as a result, this function runs
	// really early before anything else can be called in the module. therefore, here, we setup viper and cobra.

	// config file
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".") // load config.json, if present, in the current working directory

	// environment
	viper.AutomaticEnv() // not so automatic, you have to BindEnv() to make this actually happen, see spf13/viper#188

	viper.BindEnv("env") // presumably this has to happen because converting camelCase into uppercase snake case
	// isn't trivial for things like MySQLHost
	viper.BindEnv("mysql_host")
	viper.BindEnv("mysql_port")
	viper.BindEnv("mysql_user")
	viper.BindEnv("mysql_password")
	viper.BindEnv("mysql_database")
	viper.BindEnv("jwt_key")

	// define flags in cobra
	// FIXME we need to define and test making arguments *required* and set to nil if not defined
	rootCmd.Flags().StringP("environment", "e", DefaultEnv,
		`The execution environment for the application. Expected values are dev for local development and prod for production deployments.`)
	rootCmd.Flags().StringP("mysql-host", "", "", "MySQL Database Host")
	rootCmd.Flags().Uint16P("mysql-port", "", 3306, "MySQL Database Port")
	rootCmd.Flags().StringP("mysql-user", "", "", "MySQL Username")
	rootCmd.Flags().StringP("mysql-password", "", "", "MySQL Password")
	rootCmd.Flags().StringP("mysql-database", "", "", "MySQL Database")
	// wtf is StringP and why StringP vs String or StringVar or StringVarP?
	//
	//  1. The P is for POSIX, namely supporting short arguments like -P *and* long options like --port. If you don't
	//     need short options, use String()
	//  2. String and StringP return a *string to where the variable value will be stored, as opposed to StringVar and
	//     StringVarP which force you to pass a reference to an address in which to store the value.
	//
	// I lost a lot of time to trying to figure this out, namely because the documentation is sparse and is many levels
	// deep. viper integrates with cobra, which internally uses pflags, which is a drop-in replacement to augment the
	// standard library's flags, which don't do the POSIX thing.
	rootCmd.Flags().StringP("jwt-key", "", "", "JWT Key in Base-64 Format")

	// bind cobra flags to viper
	viper.BindPFlag("env", rootCmd.Flags().Lookup("environment"))
	viper.BindPFlag("mysql_host", rootCmd.Flags().Lookup("mysql-host"))
	viper.BindPFlag("mysql_port", rootCmd.Flags().Lookup("mysql-port"))
	viper.BindPFlag("mysql_user", rootCmd.Flags().Lookup("mysql-user"))
	viper.BindPFlag("mysql_password", rootCmd.Flags().Lookup("mysql-password"))
	viper.BindPFlag("mysql_database", rootCmd.Flags().Lookup("mysql-database"))
	viper.BindPFlag("jwt_key", rootCmd.Flags().Lookup("jwt-key"))
}

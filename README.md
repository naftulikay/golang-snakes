# golang-snakes ![Build Status][build.svg]

Structured and actually working example of a Golang CLI application using [spf13/cobra][cobra] and [spf13/viper][viper].

Why? Because it's needlessly hard to figure out on your own, the quality of the documentation leaves a lot to be
desired, and the endless examples out there on forums and StackOverflow and blog posts simply don't cover building
something reasonable like this. I lost an absurd amount of time tinkering with what amounts to not more than a few
hundred lines of code just to get CLI flag parsing and environment variable binding working.

**It should not be this hard to do something so basic.**

## Goals

The goals of this demo project are very simple:

 - Provide structured, nested commands for taking actions on certain data-types from the command-line.
 - Try to detect configuration with the following priority order:
   1. CLI flags
   2. Environment Variables
   3. An optional config file.
 - Serve as an example of how to structure a CLI utility's code.

## Key Takeaways

 1. Every Go file can have an `init` function which runs exactly once for each file for the lifetime of the process.
 2. You'll want to run **Cobra** flag configuration on a per-command basis in an `init` function.
 3. You'll want to run _some_ **Viper** configuration in an `init` function, perhaps in `cmd/root.go`, as it'll only
    need to run once.
 4. You absolutely **must** run **Viper** configuration on a per-command basis in the command's `PreRun` function.

In your `cmd/root.go`, in your `init` function, this should be run:

```golang
func init() {
	// make viper aware of environment variables
	viper.AutomaticEnv()
}
```

In, for example, `cmd/users/createUser.go`, which is the `./golang-snakes users create` command, you'll want to run the
following to setup the CLI flags:

```golang
var createUserCommand = &cobra.Command{ /* ... */ }

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
```

Take note of `viper.BindEnv`: these calls can be hoisted as high up in your command call stack as you'd like, and can be
collapsed to find a union of all environment variables that your application knows about. It seems that this method is
also idempotent, so if you call it with the same environment variable names in multiple `init` methods, there doesn't
seem to be a problem with this.

Finally, and this is critical: your binding _between_ Cobra and Viper **must** occur in your command's `PreRun` handler:

```golang
var createUserCommand = &cobra.Command{
	Use: "create",
	PreRun: func(cmd *cobra.Command, args []string) {
            // NOTE the binding of viper to flags needs to happen here otherwise the bindings are global and overwrite each other
            const EnvKey = "api_key"
            const CliKey = "api-key"
            
            // bind flags to env variables
            if err := viper.BindPFlag(EnvKey, cmd.Flags().Lookup(CliKey)); err != nil {
                log.Fatalf("Unable to bind flag: %s", err)
            }
    },
}
```

To finish everything up, use a struct and `viper.Unmarshal` to load from environment variables and CLI flags:

```golang
var createUserCommand = &cobra.Command{
	Use: "create",
	PreRun: /*...*/,
	Run: func(cmd *cobra.Command, args []string) {
            type Config struct {
                APIKey string `mapstructure:"api_key"`
            }
            
            var config Config
            
            if err := viper.Unmarshal(&config); err != nil {
                log.Fatalf("Unable to unmarshal viper config: %s", err)
            }
            
            log.Printf("Create User: %+v", config)
	}
}
```

The fact that this is so specific and yet is entirely not covered by the documentation and available examples is why
this repository exists.

## License

Licensed at your discretion under either:

 - [Apache Software License, Version 2.0](./LICENSE-APACHE)
 - [MIT License](./LICENSE-MIT)

 [build.svg]: https://github.com/naftulikay/golang-snakes/actions/workflows/ci.yml/badge.svg
 [cobra]: https://github.com/spf13/cobra
 [viper]: https://github.com/spf13/viper
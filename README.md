# golang-snakes

![Fucking Snakes](./images/fucking-snakes.gif)

A demo project in Golang to demonstrate _how to actually setup [`cobra`][cobra] and [`viper`][viper]_ because **it took
me #$@7ing ten hours to arrive at something that works.** In spite of these libraries being presumably the most used
libraries in Golang for parsing CLI arguments, environment variables, and config files, it was extremely difficult to
get something basic like this started.

I have many words to describe how frustrating this was, but I will leave these words to the reader's imagination. The
very first hurdle to clear in a language should not need to be a blind polevault. Note that I have written large
applications in Golang before, but did not need to parse and merge config together as I did here. "Just read the docs"
wasn't helpful either, I had to actually read source code and traverse GitHub issue trackers to get things working, and
I was unable to find an example demonstrating what I actually wanted to do.

## Goal

The goal for this project was fairly simple, or so I thought. I needed to throw together a demonstration of being able
to accept application configuration in a variety of different ways, namely and in priority order:

 1. CLI Flags (e.g. `-p`/`--mysql-port`)
 2. Environment Variables (e.g. `MYSQL_PORT`)
 3. A Configuration File (e.g. `{ "mysql_port": ... }`)

The priority here is important. CLI flags should be the highest priority, falling back to environment variables, and
finally falling back further to an on-disk configuration file. It was important to support all three avenues, as this
application should be runnable as a local binary as well as in a Docker container, which accepts environment variables
as the _de facto_ method of passing configuration to the container.

The application should merge all of these sources into one `struct` in a global variable (which, didn't we all agree
that global variables were a bad idea like >30 years ago?) to be accessed by the application.

The application does nothing but print out the configuration it has loaded and parsed.

## Building and Running

`go build` should be all that is required to build the application.

A `config.default.json` is present in the repository, copy this to `config.json` to have the application load this file
for configuration values.

> **NOTE:** The string values in `config.default.json` are prefixed with `from_cfg:`. This prefix is meaningless, it
> only exists to make it easier for you to see that it is successfully falling back to the config file values.

Run `./golang-snakes` to execute the application. As an example, we will focus on setting the MySQL port for the
application. Set `mysql_port` in `config.json` to `1234`, and execute the application to see that this is what it
prints out:

```shell
$ go build
$ ./golang-snakes
&{Env:from_cfg:dev MySQLHost:from_cfg:localhost MySQLPort:1234 MySQLUser:from_cfg:root MySQLPassword:from_cfg:password MySQLDatabase:from_cfg:default JWTKey:}
```

Next, let's override that with an environment variable:

```shell
$ MYSQL_PORT=5678 ./golang-snakes
Parsed: &{Env:from_cfg:dev MySQLHost:from_cfg:localhost MySQLPort:5678 MySQLUser:from_cfg:root MySQLPassword:from_cfg:password MySQLDatabase:from_cfg:default JWTKey:}
```

Finally, let's demonstrate that CLI flags do override everything:

```shell
$ MYSQL_PORT=5678 ./golang-snakes --mysql-port 9999
Parsed: &{Env:from_cfg:dev MySQLHost:from_cfg:localhost MySQLPort:9999 MySQLUser:from_cfg:root MySQLPassword:from_cfg:password MySQLDatabase:from_cfg:default JWTKey:}
```

Yup, it works!

## Why Did this Take Ten Hours?

It was many, many things. My goal was to bind to a struct, but other than a few snippets in the viper README and a few
really contrived examples online, it felt like this was a feature that was only _kinda_ supported. Also, it didn't help
that as per [spf13/viper#188](https://github.com/spf13/viper/issues/188), in spite of viper saying that it supports
environment variables, it kind of doesn't, unless you force it to.

It tells you to call:

```go
viper.AutomaticEnv()
```

And then everything should be magical! Except, it doesn't seem to do anything. I lost a lot of time to this until I
tried actually binding the environment variables myself:

```go
viper.AutomaticEnv() // not so automatic
viper.BindEnv("mysql_port")
```

With this, it finally started actually using environment variables.

By the way, my struct that I'm binding to looks like this, with fields omitted for brevity:

```go
type config struct {
	MySQLPort uint16 `mapstructure:"mysql_port"`
}
```

So environment variables were just not supported by default. It can't do the magical thing where it reads the struct
and infers environment variable names from it. Presumably this is related to Golang using `PascalCase`, which makes
conversion to `SCREAMING_SNAKE_CASE` for environment variable names impossible for certain cases. For instance, how do
you reach `MYSQL_PORT` from `MySQLPort`? It _might_ be possible to infer this, but the logic would be complex and
there might be edge cases. I'm too tired to give this more thought.

Part of the time suck was probably my fault. Global variables are a disease to me and in every other programming
language, they are heavily discouraged. If I would have just given up and did things like `StringVarP` and bound to
global variables, it might have saved me some time.

In retrospect, it _is_ possible to avoid doing global variables altogether, but at some point I did just give up and
was willing to compromise to have a global struct variable and have the cobra command be a global variable too. At this
point I just don't care anymore.

Next were hurdles with cobra, trying to figure out when/how/why to use the different flag binding methods and what they
did. Viper integrates with cobra, which uses pflags, so I went very far down a rabbit hole of source code and
documentation. 

I started a separate project just to test cobra, and it has very opinionated ways about how you are to run things using
a `Run` closure on the `*cobra.Command` struct. Once that was working, I began trying to integrate cobra and viper, and
this was also extremely frustrating.

I tried to do the config reading and the unmarshalling within the `init` function in my `cmd` module, and I kept getting
`nil` pointer errors and crashes that I couldn't figure out. On a whim, I decided to instead do the viper config loading
and unmarshalling within the `Run` closure attached to the root command, and this seemed to work.

I'm _pretty sure_ that everything is working, but I'd really like to build a bunch of shell test cases to see that
variable precedence is working in all cases.

## License

This unfortunate adventure is licensed permissively under your choice of:

 - [Apache Software License, Version 2.0](./LICENSE-APACHE)
 - [MIT License](./LICENSE-MIT)

 [cobra]: https://github.com/spf13/cobra
 [viper]: https://github.com/spf13/viper
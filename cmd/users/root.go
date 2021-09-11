package users

import "github.com/spf13/cobra"

var (
	usersCommand = &cobra.Command{
		Use: "users",
		Short: "Interact with user objects.",
	}
)

func Commands() []*cobra.Command {
	return []*cobra.Command{usersCommand}
}

func ChildCommands() []*cobra.Command {
	return []*cobra.Command{createUserCommand,listUsersCommand}
}

func init() {
	usersCommand.AddCommand(ChildCommands()...)
}
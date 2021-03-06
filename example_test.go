package cmdmux

import (
	"fmt"
	"strings"
)

// This running example is a todo list manager.
func Example() {
	CreateTodo := func(name string, args []string) {
		fmt.Println("created:", args[0])
	}
	ListTodos := func(name string, args []string) {
		fmt.Println("milk")
		fmt.Println("eggs")
		fmt.Println("motor boat")
	}

	RegisterFunc("new", CreateTodo)
	RegisterFunc("list", ListTodos)

	Exec("todo", []string{"new", "quinoa"}) // usually os.Args[0], os.Args[1:]
	// Output: created: quinoa
}

// Using the running example...
func ExampleCommands() {
	commands := Commands()
	commandstr := strings.Join(commands, ", ")
	fmt.Println("available commands:", commandstr)
	// Output: available commands: new, list
}

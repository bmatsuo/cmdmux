package cmdmux

import (
	"fmt"
)

func Example() {
	CreateTodo := func(name string, args []string) {
		fmt.Println("created:", args)
	}

	ListTodos := func(name string, args []string) {
		fmt.Println("milk", "eggs", "motor boat")
	}

	RegisterFunc("new", CreateTodo)
	RegisterFunc("list", ListTodos)

	Exec("todo", []string{"new", "quinoa"}) // usually os.Args[0], os.Args[1:]
	// Output:
	// created: [quinoa]
}

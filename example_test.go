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
	fmt.Println(Commands())
	Exec("todo", []string{"new", "quinoa"})
	// Output:
	// [new list]
	// created: [quinoa]
}

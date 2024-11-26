package main

import (
	"flag"
	"github.com/sajib/mowtodo/pkg/todo"
)

func main() {
	t := todo.Init()

	addPtr := flag.String("a", "", "add a task")
	priorityPtr := flag.String("p", "low", "set task priority (low, medium, high)")
	dueDatePtr := flag.String("d", "", "set task due date (YYYY-MM-DD)")
	remPtr := flag.Int("r", 0, "remove a task")
	togglePtr := flag.Int("t", 0, "toggle done for a task")
	editPtr := flag.Bool("e", false, "edit todo file")
	listDonePtr := flag.Bool("ld", false, "list done tasks")
	listUndonePtr := flag.Bool("lu", false, "list undone tasks")
	hideProgressPtr := flag.Bool("hp", false, "hide progress bar")

	flag.Parse()

	newTask := *addPtr
	priority := *priorityPtr
	dueDate := *dueDatePtr
	remTaskNum := *remPtr
	toggleTaskNum := *togglePtr
	editFlag := *editPtr

	t.ListDone = !*listUndonePtr
	t.ListUndone = !*listDonePtr
	t.ShowProgress = !*hideProgressPtr

	switch {
	case remTaskNum != 0:
		t.RemTask(remTaskNum)
	case newTask != "":
		t.AddTask(newTask, priority, dueDate)
	case toggleTaskNum != 0:
		t.ToggleTask(toggleTaskNum)
	case editFlag:
		t.OpenEditor()
	case t.ListUndone:
		t.PrintList()
	case t.ListDone:
		t.PrintList()
	default:
		t.PrintList()
	}
}

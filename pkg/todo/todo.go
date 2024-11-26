package todo

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/sajib/mowtodo/pkg/file"
	"github.com/sajib/mowtodo/pkg/pprint"
)

type Task struct {
	Description string
	IsDone      bool
	Priority    string
	DueDate     string
}

type Todo struct {
	filePath     string
	doneCount    float64
	undoneCount  float64
	Tasks        []Task
	ListDone     bool
	ListUndone   bool
	ShowProgress bool
}

func Init() *Todo {
	todo := new(Todo)

	configDir := filepath.Join(os.Getenv("HOME"), ".config", "todo")
	filePath := filepath.Join(configDir, "todo.txt")

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			pprint.Error(fmt.Sprintf("Can't create config directory\n%s\n", err))
			return nil
		}

		file, err := os.Create(filePath)
		defer file.Close()
		if err != nil {
			pprint.Error(fmt.Sprintf("Can't create todo.txt file\n%s\n", err))
			return nil
		}
	} else if err != nil {
		pprint.Error(fmt.Sprintf("Can't check file\n%s\n", err))
		return nil
	}

	todo.filePath = filePath
	todo.ListDone = true
	todo.ListUndone = true
	todo.ShowProgress = true
	todo.loadTasks()

	return todo
}

func (t *Todo) OpenEditor() {
	// Retrieve the editor from environment variables
	editor := os.Getenv("EDITOR")

	// Check if the editor environment variable is set
	if editor == "" {
		pprint.Error("Can't open editor [no $EDITOR env]")
		return
	}

	// Run the editor with the todo file as an argument
	cmd := exec.Command(editor, t.filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		pprint.Error(fmt.Sprintf("Failed to open editor\n%s\n", err))
		return
	}
}

func (t *Todo) loadTasks() {
	file := file.Open(t.filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		task := parseTask(line)
		t.Tasks = append(t.Tasks, task)
	}

	if err := scanner.Err(); err != nil {
		pprint.Error("Can't read file")
	}
}

func parseTask(line string) Task {
	// Split the line by "|"
	parts := strings.Split(line, "|")

	// If the line doesn't have at least 3 parts, it's invalid
	if len(parts) < 3 {
		// You can handle this more gracefully, e.g., by skipping invalid lines or setting default values
		fmt.Printf("Invalid task format: %s\n", line)
		return Task{}
	}

	description := strings.TrimSpace(strings.TrimPrefix(parts[0], "[]"))
	priority := strings.TrimSpace(strings.TrimPrefix(parts[1], "Priority:"))
	dueDate := strings.TrimSpace(strings.TrimPrefix(parts[2], "Due:"))

	// Ensure missing values are handled properly
	if priority == "" {
		priority = "None"
	}
	if dueDate == "" {
		dueDate = "None"
	}

	return Task{Description: description, Priority: priority, DueDate: dueDate}
}

func (t *Todo) AddTask(task string, priority string, dueDate string) {
	// Validate priority before adding the task
	priority = normalizePriority(priority)
	if priority == "" {
		// If the priority is invalid, return early and show an error
		pprint.Error("Invalid priority. Allowed values are 'low', 'medium', 'high'")
		return
	}

	// Ensure due date is not empty
	if dueDate == "" {
		dueDate = "None"
	}

	// Format the task with priority and due date
	var ftask string
	if file.Size(t.filePath) == 0 {
		ftask = fmt.Sprintf("[] %s |  %s |  %s", task, priority, dueDate)
	} else {
		ftask = fmt.Sprintf("\n[] %s |  %s |  %s", task, priority, dueDate)
	}

	// Write the formatted task to the file
	file.Write(t.filePath, ftask, os.O_APPEND)

	// Print the updated list
	t.PrintList()
}

func (t Todo) ToggleTask(taskId int) {
	newLines := make([]string, 0)
	todoFile := file.Open(t.filePath)
	defer todoFile.Close()

	i := 1
	scanner := bufio.NewScanner(todoFile)
	for scanner.Scan() {
		line := scanner.Text()
		if i == taskId {
			if strings.Contains(line, "[X]") {
				line = strings.Replace(line, "[X", "[", 1)
			} else if strings.Contains(line, "[]") {
				line = strings.Replace(line, "[", "[X", 1)
			}
		}
		newLines = append(newLines, line)
		i++
	}

	if err := scanner.Err(); err != nil {
		pprint.Error(fmt.Sprintf("Can't read %s", t.filePath))
	}

	file.Write(t.filePath, strings.Join(newLines, "\n"), os.O_TRUNC)

	t.PrintList()
}

func (t *Todo) RemTask(taskId int) {
	if taskId > 0 && taskId <= len(t.Tasks) {
		// Remove the task from the list
		t.Tasks = append(t.Tasks[:taskId-1], t.Tasks[taskId:]...)
		t.saveToFile()
		t.PrintList()
	} else {
		pprint.Error("Invalid Task ID")
	}
}

// Modify PrintList to colorize priority

// Modify PrintList to colorize priority and handle errors
func (t *Todo) PrintList() {
	// Open the task file
	file, err := os.Open(t.filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Print the formatted header for the task list
	pprint.Print("      ToDo List      \n\n", color.FgCyan, color.Bold)
	fmt.Printf("%-3s %-50s %-10s %-15s\n", "ID", "Description", "Priority", "Due Date")
	fmt.Println(strings.Repeat("─", 80)) // Print a separator line

	// Variables to keep track of task ID and count of done/undone tasks
	i := 1
	for scanner.Scan() {
		line := scanner.Text()

		// Extract task description, priority, and due date
		taskDescription, priority, dueDate, err := parseTaskDetails(line)
		if err != nil {
			// If an error occurs during parsing, print the error and skip this line
			fmt.Println("Error parsing task:", err)
			continue
		}

		// Colorize priority
		var priorityColor color.Attribute
		switch strings.ToLower(priority) {
		case "high":
			priorityColor = color.FgRed
		case "medium":
			priorityColor = color.FgYellow
		case "low":
			priorityColor = color.FgGreen
		default:
			priorityColor = color.FgWhite
		}

		// Print each task with formatted output
		if strings.Contains(line, "[X]") && t.ListDone {
			// Task is done, show green checkmark icon only
			pprint.Print(fmt.Sprintf("%2d  ", i), color.FgWhite, color.Faint)
			pprint.Print("  ", color.FgGreen)
			fmt.Printf("%-50s ", taskDescription)
			pprint.Print(fmt.Sprintf("%-10s", priority), priorityColor)
			fmt.Printf("%-15s\n", dueDate)
			t.doneCount++
		} else if strings.Contains(line, "[]") && t.ListUndone {
			// Task is undone, show red checkbox icon
			pprint.Print(fmt.Sprintf("%2d  ", i), color.FgWhite, color.Faint)
			pprint.Print("⨀  ")
			fmt.Printf("%-50s ", taskDescription)
			pprint.Print(fmt.Sprintf("%-10s", priority), priorityColor)
			fmt.Printf("%-15s\n", dueDate)
			t.undoneCount++
		}
		i++
	}

	// Print progress bar (if applicable)
	t.printProgress()

	// Handle any errors that occurred while scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

// Helper function to parse task details (description, priority, and due date)
func parseTaskDetails(line string) (description, priority, dueDate string, err error) {
	// Set default values
	description = ""
	priority = "N/A"
	dueDate = "N/A"

	// Split the line to extract the task description, priority, and due date
	parts := strings.Split(line, "|")

	// Handle cases where the line does not have the expected parts
	if len(parts) > 0 {
		// Extract description (assuming description is everything before the first '|')
		description = strings.TrimSpace(parts[0][3:]) // Skip the checkbox "[ ]" or "[X]"

		// Check if there are priority and due date fields
		if len(parts) > 1 {
			priority = strings.TrimSpace(parts[1])
			// Normalize and validate the priority
			priority = normalizePriority(priority)
			if priority == "" {
				err = fmt.Errorf("invalid priority. Allowed values are 'low', 'medium', 'high'")
				return description, "", "", err
			}
		}
		if len(parts) > 2 {
			dueDate = strings.TrimSpace(parts[2])
		}
	}

	// Ensure priority and due date have default values if they are missing
	if priority == "" {
		priority = "none"
	}
	if dueDate == "" {
		dueDate = "none"
	}

	return description, priority, dueDate, nil
}

// Helper function to normalize the priority to 'Low', 'Medium', or 'High'
func normalizePriority(priority string) string {
	switch strings.ToLower(priority) {
	case "low", "l":
		return "Low"
	case "medium", "m":
		return "Medium"
	case "high", "h":
		return "High"
	default:
		return ""
	}
}

func (t Todo) printProgress() {
	if !t.ShowProgress || !t.ListDone || !t.ListUndone || (t.doneCount == 0 && t.undoneCount == 0) {
		return
	}

	lineWidth := 50.0
	total := t.doneCount + t.undoneCount
	doneWidth := lineWidth * (t.doneCount / total)
	undoneWidth := (lineWidth - doneWidth) - 1
	percentage := (t.doneCount / total) * 100

	fmt.Print("\n")
	for i := 0; i < int(doneWidth); i++ {
		pprint.Print("━", color.FgGreen)
	}

	if t.undoneCount != 0 {
		pprint.Print("╺", color.FgWhite, color.Faint)
	}

	for i := 0; i < int(undoneWidth); i++ {
		pprint.Print("━", color.FgWhite, color.Faint)
	}

	fmt.Printf(" %d%% Done\n", int(percentage))
}

func (t *Todo) saveToFile() {
	lines := []string{}
	for _, task := range t.Tasks {
		status := "[]"
		if task.IsDone {
			status = "[X]"
		}
		lines = append(lines, fmt.Sprintf("%s %s |  %s |  %s", status, task.Description, task.Priority, task.DueDate))
	}
	file.Write(t.filePath, strings.Join(lines, "\n"), os.O_TRUNC)
}

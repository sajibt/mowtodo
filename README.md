```markdown
# mowtodo Command Usage

**Table of Contents**

<!-- toc -->

- [About](#about)
  * [Installing](#installing)
  * [Uninstalling](#uninstalling)
  * [Build From Source](#build-from-source)
  * [Usage](#usage)

<!-- tocstop -->

## About

mowtodo is a simple command-line to-do list manager to manage your tasks. Written in GO.

### Installing

```console
GOBIN=<your install path> go install ./main/todo
```

### Uninstalling

```console
rm -rf <your install path>/mowtodo
```

### Build From Source

Install Go and build with this command:

```console
go build ./main/todo
```

### Usage

To add a task to the list:

```console
mowtodo -a "<Task Description>"
```

Toggle a task as done or undone:

```console
mowtodo -t <Task ID>
```

Remove a task from the list:

```console
mowtodo -r <Task ID>
```

Opens editor to edit the raw file of the list (it uses the `$EDITOR` environment variable):

```console
mowtodo -e
```

List done tasks:

```console
mowtodo -ld
```

List undone tasks:

```console
mowtodo -lu
```

Hide progress bar (can be used with other options):

```console
mowtodo -hp
```

## Command Reference

### Add a Task

```console
mowtodo -a "Task Description"
```
- **Description**: Adds a new task with the specified description.

### Edit Todo File

```console
mowtodo -e
```
- **Description**: Opens the todo file for manual editing.

### Hide Progress Bar

```console
mowtodo -hp
```
- **Description**: Hides the progress bar when displaying tasks.

### List Done Tasks

```console
mowtodo -ld
```
- **Description**: Lists all completed (done) tasks.

### List Undone Tasks

```console
mowtodo -lu
```
- **Description**: Lists all pending (undone) tasks.

### Remove a Task by ID

```console
mowtodo -r <task_id>
```
- **Description**: Removes a task by its ID. The task ID is obtained from the task listing.

### Toggle Task Completion by ID

```console
mowtodo -t <task_id>
```
- **Description**: Toggles the completion status of a task. If the task is undone, it will be marked as done, and vice versa.

## Example Usage

- Add a task: 

```console
mowtodo -a "Buy groceries"
```

- Remove a task by ID: 

```console
mowtodo -r 2
```

- List all undone tasks:

```console
mowtodo -lu
```

- Toggle task completion by ID: 

```console
mowtodo -t 3
```
```


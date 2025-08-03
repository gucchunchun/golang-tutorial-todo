# golang-tutorial-todo
## 📌 Project: ToDo CLI App
Goal: Build a command-line ToDo application in Go to practice CLI programming, input validation, time handling, and structured code organization.
### ✅ Core Requirements
Add Task

Command: todo add "Buy groceries" --due 2025-07-28

Store:

ID (auto-increment)

Task name

Status (pending or done)

Created time

Due date

Validate task name using regexp (Task 27)

List Tasks

Command: todo list

Show all tasks with:

ID | Name | Status | CreatedAt | DueDate | TimeLeft

Use time package (Task 28) to calculate time left.

Mark Task as Done

Command: todo done 3

Change status from pending → done.

Delete Task

Command: todo delete 3

Data Persistence

Store tasks in JSON file (tasks.json).

Use encoding/json for save/load.

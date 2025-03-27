package pages

import "baby-blog/types"

func Todos(pageData map[string]interface{}, db *types.Models) map[string]interface{} {
	// Get all todos from the database
	todos, err := db.Todo.GetAll()

	if err != nil {
		pageData["Failure"] = "Failed to load todos"
		return pageData
	}
	// Convert database journals to the format expected by the template
	todosData := make([]map[string]string, 0, len(todos))
	for _, todo := range todos {
		completedValue := "false"
		if todo.Completed {
			completedValue = "true"
		}
		todosData = append(todosData, map[string]string{
			"task":      todo.Task,
			"completed": completedValue,
			"date":      todo.CreatedAt.Format("2006-01-02"),
		})
	}
	// Return the data from the database
	pageData["todos"] = todosData
	pageData["todoCount"] = len(todosData)
	return pageData

}

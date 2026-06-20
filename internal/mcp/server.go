package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nouchka/mcp-tududi/internal/client"
)

// Server represents the MCP HTTP server
type Server struct {
	port   string
	client *client.TududuClient
	mux    *http.ServeMux
}

// NewServer creates a new MCP server
func NewServer(port string, tududuClient *client.TududuClient) *Server {
	s := &Server{
		port:   port,
		client: tududuClient,
		mux:    http.NewServeMux(),
	}
	
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/tools", s.handleListTools)
	s.mux.HandleFunc("/call_tool", s.handleCallTool)
}

// Start starts the MCP HTTP server
func (s *Server) Start() error {
	log.Printf("Starting MCP server on port %s", s.port)
	return http.ListenAndServe(":"+s.port, s.mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleListTools(w http.ResponseWriter, r *http.Request) {
	tools := s.getAvailableTools()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tools": tools,
	})
}

func (s *Server) handleCallTool(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name   string                 `json:"name"`
		Input  map[string]interface{} `json:"input"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read request body"})
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	result, err := s.callTool(req.Name, req.Input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": result,
	})
}

func (s *Server) getAvailableTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "tududi_list_tasks",
			"description": "List all tasks from Tududi",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "tududi_create_task",
			"description": "Create a new task in Tududi",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Task title",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Task description",
					},
					"projectId": map[string]interface{}{
						"type":        "string",
						"description": "Project ID to assign task to",
					},
					"areaId": map[string]interface{}{
						"type":        "string",
						"description": "Area ID to assign task to",
					},
					"dueDate": map[string]interface{}{
						"type":        "string",
						"description": "Due date in ISO format",
					},
					"priority": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"low", "medium", "high"},
						"description": "Task priority",
					},
				},
				"required": []string{"title"},
			},
		},
		{
			"name":        "tududi_update_task",
			"description": "Update an existing task",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Task ID",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "New task title",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "New task description",
					},
					"completed": map[string]interface{}{
						"type":        "boolean",
						"description": "Task completion status",
					},
					"priority": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"low", "medium", "high"},
						"description": "Task priority",
					},
				},
				"required": []string{"id"},
			},
		},
		{
			"name":        "tududi_delete_task",
			"description": "Delete a task",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Task ID to delete",
					},
				},
				"required": []string{"id"},
			},
		},
		{
			"name":        "tududi_complete_task",
			"description": "Mark a task as complete",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Task ID to complete",
					},
				},
				"required": []string{"id"},
			},
		},
		{
			"name":        "tududi_list_subtasks",
			"description": "List all subtasks for a parent task",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"parentId": map[string]interface{}{
						"type":        "string",
						"description": "Parent task ID",
					},
				},
				"required": []string{"parentId"},
			},
		},
		{
			"name":        "tududi_create_subtask",
			"description": "Create a new subtask under a parent task",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"parentId": map[string]interface{}{
						"type":        "string",
						"description": "Parent task ID",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Subtask title",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Subtask description",
					},
					"priority": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"low", "medium", "high"},
						"description": "Subtask priority",
					},
				},
				"required": []string{"parentId", "title"},
			},
		},
		{
			"name":        "tududi_update_subtask",
			"description": "Update a subtask",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"parentId": map[string]interface{}{
						"type":        "string",
						"description": "Parent task ID",
					},
					"subtaskId": map[string]interface{}{
						"type":        "string",
						"description": "Subtask ID",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "New subtask title",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "New subtask description",
					},
					"completed": map[string]interface{}{
						"type":        "boolean",
						"description": "Subtask completion status",
					},
				},
				"required": []string{"parentId", "subtaskId"},
			},
		},
		{
			"name":        "tududi_delete_subtask",
			"description": "Delete a subtask",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"parentId": map[string]interface{}{
						"type":        "string",
						"description": "Parent task ID",
					},
					"subtaskId": map[string]interface{}{
						"type":        "string",
						"description": "Subtask ID to delete",
					},
				},
				"required": []string{"parentId", "subtaskId"},
			},
		},
		{
			"name":        "tududi_list_projects",
			"description": "List all projects from Tududi",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "tududi_create_project",
			"description": "Create a new project in Tududi",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Project name",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Project description",
					},
				},
				"required": []string{"name"},
			},
		},
		{
			"name":        "tududi_update_project",
			"description": "Update a project",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Project ID",
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "New project name",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "New project description",
					},
				},
				"required": []string{"id"},
			},
		},
		{
			"name":        "tududi_delete_project",
			"description": "Delete a project",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Project ID to delete",
					},
				},
				"required": []string{"id"},
			},
		},
		{
			"name":        "tududi_list_areas",
			"description": "List all areas from Tududi",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "tududi_create_area",
			"description": "Create a new area in Tududi",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Area name",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Area description",
					},
				},
				"required": []string{"name"},
			},
		},
		{
			"name":        "tududi_update_area",
			"description": "Update an area",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Area ID",
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "New area name",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "New area description",
					},
				},
				"required": []string{"id"},
			},
		},
		{
			"name":        "tududi_delete_area",
			"description": "Delete an area",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Area ID to delete",
					},
				},
				"required": []string{"id"},
			},
		},
	}
}

func (s *Server) callTool(name string, input map[string]interface{}) (interface{}, error) {
	switch name {
	case "tududi_list_tasks":
		return s.listTasks()
	case "tududi_create_task":
		return s.createTask(input)
	case "tududi_update_task":
		return s.updateTask(input)
	case "tududi_delete_task":
		return s.deleteTask(input)
	case "tududi_complete_task":
		return s.completeTask(input)
	case "tududi_list_subtasks":
		return s.listSubtasks(input)
	case "tududi_create_subtask":
		return s.createSubtask(input)
	case "tududi_update_subtask":
		return s.updateSubtask(input)
	case "tududi_delete_subtask":
		return s.deleteSubtask(input)
	case "tududi_list_projects":
		return s.listProjects()
	case "tududi_create_project":
		return s.createProject(input)
	case "tududi_update_project":
		return s.updateProject(input)
	case "tududi_delete_project":
		return s.deleteProject(input)
	case "tududi_list_areas":
		return s.listAreas()
	case "tududi_create_area":
		return s.createArea(input)
	case "tududi_update_area":
		return s.updateArea(input)
	case "tududi_delete_area":
		return s.deleteArea(input)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func (s *Server) listTasks() (interface{}, error) {
	tasks, err := s.client.ListTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Server) createTask(input map[string]interface{}) (interface{}, error) {
	title, ok := input["title"].(string)
	if !ok {
		return nil, fmt.Errorf("title is required and must be a string")
	}

	task := client.Task{
		Title: title,
	}

	if desc, ok := input["description"].(string); ok {
		task.Description = desc
	}
	if projectID, ok := input["projectId"].(string); ok {
		task.ProjectID = projectID
	}
	if areaID, ok := input["areaId"].(string); ok {
		task.AreaID = areaID
	}
	if priority, ok := input["priority"].(string); ok {
		task.Priority = priority
	}

	created, err := s.client.CreateTask(task)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Server) updateTask(input map[string]interface{}) (interface{}, error) {
	id, ok := input["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	task := client.Task{}

	if title, ok := input["title"].(string); ok {
		task.Title = title
	}
	if desc, ok := input["description"].(string); ok {
		task.Description = desc
	}
	if completed, ok := input["completed"].(bool); ok {
		task.Completed = completed
	}
	if priority, ok := input["priority"].(string); ok {
		task.Priority = priority
	}

	updated, err := s.client.UpdateTask(id, task)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Server) deleteTask(input map[string]interface{}) (interface{}, error) {
	id, ok := input["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	err := s.client.DeleteTask(id)
	if err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted"}, nil
}

func (s *Server) completeTask(input map[string]interface{}) (interface{}, error) {
	id, ok := input["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	task, err := s.client.CompleteTask(id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Server) listSubtasks(input map[string]interface{}) (interface{}, error) {
	parentID, ok := input["parentId"].(string)
	if !ok {
		return nil, fmt.Errorf("parentId is required and must be a string")
	}

	subtasks, err := s.client.ListSubtasks(parentID)
	if err != nil {
		return nil, err
	}
	return subtasks, nil
}

func (s *Server) createSubtask(input map[string]interface{}) (interface{}, error) {
	parentID, ok := input["parentId"].(string)
	if !ok {
		return nil, fmt.Errorf("parentId is required and must be a string")
	}

	title, ok := input["title"].(string)
	if !ok {
		return nil, fmt.Errorf("title is required and must be a string")
	}

	subtask := client.Task{
		Title: title,
	}

	if desc, ok := input["description"].(string); ok {
		subtask.Description = desc
	}
	if priority, ok := input["priority"].(string); ok {
		subtask.Priority = priority
	}

	created, err := s.client.CreateSubtask(parentID, subtask)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Server) updateSubtask(input map[string]interface{}) (interface{}, error) {
	parentID, ok := input["parentId"].(string)
	if !ok {
		return nil, fmt.Errorf("parentId is required and must be a string")
	}

	subtaskID, ok := input["subtaskId"].(string)
	if !ok {
		return nil, fmt.Errorf("subtaskId is required and must be a string")
	}

	subtask := client.Task{}

	if title, ok := input["title"].(string); ok {
		subtask.Title = title
	}
	if desc, ok := input["description"].(string); ok {
		subtask.Description = desc
	}
	if completed, ok := input["completed"].(bool); ok {
		subtask.Completed = completed
	}

	updated, err := s.client.UpdateSubtask(parentID, subtaskID, subtask)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Server) deleteSubtask(input map[string]interface{}) (interface{}, error) {
	parentID, ok := input["parentId"].(string)
	if !ok {
		return nil, fmt.Errorf("parentId is required and must be a string")
	}

	subtaskID, ok := input["subtaskId"].(string)
	if !ok {
		return nil, fmt.Errorf("subtaskId is required and must be a string")
	}

	err := s.client.DeleteSubtask(parentID, subtaskID)
	if err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted"}, nil
}

func (s *Server) listProjects() (interface{}, error) {
	projects, err := s.client.ListProjects()
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *Server) createProject(input map[string]interface{}) (interface{}, error) {
	name, ok := input["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name is required and must be a string")
	}

	project := client.Project{
		Name: name,
	}

	if desc, ok := input["description"].(string); ok {
		project.Description = desc
	}

	created, err := s.client.CreateProject(project)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Server) updateProject(input map[string]interface{}) (interface{}, error) {
	id, ok := input["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	project := client.Project{}

	if name, ok := input["name"].(string); ok {
		project.Name = name
	}
	if desc, ok := input["description"].(string); ok {
		project.Description = desc
	}

	updated, err := s.client.UpdateProject(id, project)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Server) deleteProject(input map[string]interface{}) (interface{}, error) {
	id, ok := input["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	err := s.client.DeleteProject(id)
	if err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted"}, nil
}

func (s *Server) listAreas() (interface{}, error) {
	areas, err := s.client.ListAreas()
	if err != nil {
		return nil, err
	}
	return areas, nil
}

func (s *Server) createArea(input map[string]interface{}) (interface{}, error) {
	name, ok := input["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name is required and must be a string")
	}

	area := client.Area{
		Name: name,
	}

	if desc, ok := input["description"].(string); ok {
		area.Description = desc
	}

	created, err := s.client.CreateArea(area)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Server) updateArea(input map[string]interface{}) (interface{}, error) {
	id, ok := input["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	area := client.Area{}

	if name, ok := input["name"].(string); ok {
		area.Name = name
	}
	if desc, ok := input["description"].(string); ok {
		area.Description = desc
	}

	updated, err := s.client.UpdateArea(id, area)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Server) deleteArea(input map[string]interface{}) (interface{}, error) {
	id, ok := input["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	err := s.client.DeleteArea(id)
	if err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted"}, nil
}

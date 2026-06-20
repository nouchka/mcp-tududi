package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TududuClient struct {
	baseURL    string
	apiKey     string
	email      string
	password   string
	httpClient *http.Client
	token      string
}

// Task represents a Tududi task
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	Priority    string    `json:"priority,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	ProjectID   string    `json:"projectId,omitempty"`
	AreaID      string    `json:"areaId,omitempty"`
	ParentID    string    `json:"parentId,omitempty"`
	Subtasks    []Task    `json:"subtasks,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Subtask represents a subtask (task with parentId)
type Subtask struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	Priority    string    `json:"priority,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	ParentID    string    `json:"parentId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Project represents a Tududi project
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Area represents a Tududi area
type Area struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NewTududuClient creates a new Tududi client
func NewTududuClient(baseURL, apiKey, email, password string) *TududuClient {
	return &TududuClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		email:   email,
		password: password,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Authenticate authenticates the client (for email/password auth)
func (c *TududuClient) Authenticate() error {
	if c.apiKey != "" {
		// API key auth, no need to authenticate
		return nil
	}

	if c.email == "" || c.password == "" {
		return fmt.Errorf("either API key or email/password must be provided")
	}

	// For now, we'll assume token-based auth is not required
	// If needed, implement login endpoint
	return nil
}

// do makes an HTTP request with proper authentication
func (c *TududuClient) do(method, path string, body interface{}) (*http.Response, error) {
	url := c.baseURL + path
	
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	if c.email != "" {
		req.Header.Set("X-Email", c.email)
	}
	if c.password != "" {
		req.Header.Set("X-Password", c.password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// ListTasks returns all tasks
func (c *TududuClient) ListTasks() ([]Task, error) {
	resp, err := c.do("GET", "/api/tasks", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTask returns a specific task
func (c *TududuClient) GetTask(id string) (*Task, error) {
	resp, err := c.do("GET", "/api/tasks/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

// CreateTask creates a new task
func (c *TududuClient) CreateTask(task Task) (*Task, error) {
	resp, err := c.do("POST", "/api/tasks", task)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var created Task
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}
	return &created, nil
}

// UpdateTask updates a task
func (c *TududuClient) UpdateTask(id string, task Task) (*Task, error) {
	task.ID = id
	resp, err := c.do("PUT", "/api/tasks/"+id, task)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updated Task
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// DeleteTask deletes a task
func (c *TududuClient) DeleteTask(id string) error {
	resp, err := c.do("DELETE", "/api/tasks/"+id, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// CompleteTask marks a task as complete
func (c *TududuClient) CompleteTask(id string) (*Task, error) {
	task := Task{Completed: true}
	return c.UpdateTask(id, task)
}

// ListSubtasks returns all subtasks for a parent task
func (c *TududuClient) ListSubtasks(parentID string) ([]Subtask, error) {
	resp, err := c.do("GET", "/api/tasks/"+parentID+"/subtasks", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var subtasks []Subtask
	if err := json.NewDecoder(resp.Body).Decode(&subtasks); err != nil {
		return nil, err
	}
	return subtasks, nil
}

// CreateSubtask creates a new subtask under a parent task
func (c *TududuClient) CreateSubtask(parentID string, subtask Task) (*Task, error) {
	subtask.ParentID = parentID
	resp, err := c.do("POST", "/api/tasks/"+parentID+"/subtasks", subtask)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var created Task
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}
	return &created, nil
}

// UpdateSubtask updates a subtask
func (c *TududuClient) UpdateSubtask(parentID, subtaskID string, subtask Task) (*Task, error) {
	subtask.ID = subtaskID
	subtask.ParentID = parentID
	resp, err := c.do("PUT", "/api/tasks/"+parentID+"/subtasks/"+subtaskID, subtask)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updated Task
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// DeleteSubtask deletes a subtask
func (c *TududuClient) DeleteSubtask(parentID, subtaskID string) error {
	resp, err := c.do("DELETE", "/api/tasks/"+parentID+"/subtasks/"+subtaskID, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListProjects returns all projects
func (c *TududuClient) ListProjects() ([]Project, error) {
	resp, err := c.do("GET", "/api/projects", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}
	return projects, nil
}

// CreateProject creates a new project
func (c *TududuClient) CreateProject(project Project) (*Project, error) {
	resp, err := c.do("POST", "/api/projects", project)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var created Project
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}
	return &created, nil
}

// UpdateProject updates a project
func (c *TududuClient) UpdateProject(id string, project Project) (*Project, error) {
	project.ID = id
	resp, err := c.do("PUT", "/api/projects/"+id, project)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updated Project
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// DeleteProject deletes a project
func (c *TududuClient) DeleteProject(id string) error {
	resp, err := c.do("DELETE", "/api/projects/"+id, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListAreas returns all areas
func (c *TududuClient) ListAreas() ([]Area, error) {
	resp, err := c.do("GET", "/api/areas", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var areas []Area
	if err := json.NewDecoder(resp.Body).Decode(&areas); err != nil {
		return nil, err
	}
	return areas, nil
}

// CreateArea creates a new area
func (c *TududuClient) CreateArea(area Area) (*Area, error) {
	resp, err := c.do("POST", "/api/areas", area)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var created Area
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}
	return &created, nil
}

// UpdateArea updates an area
func (c *TududuClient) UpdateArea(id string, area Area) (*Area, error) {
	area.ID = id
	resp, err := c.do("PUT", "/api/areas/"+id, area)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updated Area
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// DeleteArea deletes an area
func (c *TududuClient) DeleteArea(id string) error {
	resp, err := c.do("DELETE", "/api/areas/"+id, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// Package backlog reads the workshop's Beans backlog through the pinned Beans CLI.
//
// Everything in this package is read-only except AddNote, which the presenter-only
// write demonstration uses and which the server refuses to expose unless it is
// explicitly enabled against a backlog marked disposable.
package backlog

// Task is the full task view returned by get_task.
type Task struct {
	ID        string   `json:"id"`
	Slug      string   `json:"slug,omitempty"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Status    string   `json:"status"`
	Type      string   `json:"type,omitempty"`
	Priority  string   `json:"priority,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	BlockedBy []string `json:"blocked_by,omitempty"`
	Blocking  []string `json:"blocking,omitempty"`
	Parent    string   `json:"parent,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
	// Version is the Beans etag for this task's current content. Callers record it
	// so a later read can detect that the accepted contract changed.
	Version string `json:"version"`
	// AcceptanceCriteria is the "## Acceptance criteria" section of the body, split
	// into lines, when the task has one. The body remains authoritative.
	AcceptanceCriteria []string `json:"acceptance_criteria,omitempty"`
	// Backlog identifies which backlog answered, without leaking a host path.
	Backlog string `json:"backlog"`
}

// TaskSummary is the bounded list view returned by list_tasks.
type TaskSummary struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Status    string   `json:"status"`
	Type      string   `json:"type,omitempty"`
	Priority  string   `json:"priority,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	BlockedBy []string `json:"blocked_by,omitempty"`
	Version   string   `json:"version"`
}

// ListFilter is the complete filter vocabulary of list_tasks. Anything not
// expressible here is not queryable through this adapter by design.
type ListFilter struct {
	Status string
	Tag    string
	Limit  int
}

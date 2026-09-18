package backlog

import "fmt"

// Code is a stable, machine-readable error code returned to MCP clients.
type Code string

const (
	CodeInvalidID       Code = "invalid_id"
	CodeInvalidFilter   Code = "invalid_filter"
	CodeTaskNotFound    Code = "task_not_found"
	CodeBacklogUnavail  Code = "backlog_unavailable"
	CodeBacklogTimeout  Code = "backlog_timeout"
	CodeBacklogTooLarge Code = "backlog_output_too_large"
	CodeBacklogBadData  Code = "backlog_unreadable"
	CodeWriteDisabled   Code = "write_disabled"
)

// Error carries a typed code plus a message that is safe to show an agent: it
// never includes the host's absolute paths or the raw CLI invocation.
type Error struct {
	Code    Code
	Message string
	// Detail is logged to stderr but not returned to the client.
	Detail string
}

func (e *Error) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func newError(code Code, message string, detail string) *Error {
	return &Error{Code: code, Message: message, Detail: detail}
}

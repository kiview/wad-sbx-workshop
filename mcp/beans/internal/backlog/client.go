package backlog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// Config pins every path the adapter uses. Nothing is discovered: the adapter never
// searches upward for a .beans.yml, so it cannot accidentally answer from the host
// user's unrelated backlog.
type Config struct {
	// BeansBin is an absolute path to the pinned Beans executable.
	BeansBin string
	// ConfigFile is an absolute path to the workshop's .beans.yml.
	ConfigFile string
	// DataDir is an absolute path to the workshop's .beans data directory.
	DataDir string
	// Name labels the backlog in responses instead of exposing DataDir.
	Name string
	// Timeout bounds every CLI invocation.
	Timeout time.Duration
	// MaxOutputBytes bounds how much CLI output the adapter will read.
	MaxOutputBytes int64
	// AllowNoteWrites enables the presenter-only add_task_note tool. It is only
	// honoured when DataDir also contains the disposable marker file.
	AllowNoteWrites bool

	Logger *slog.Logger
}

// DisposableMarker must exist in DataDir before any write is permitted. This is a
// real check on every write call, not a tool annotation.
const DisposableMarker = ".workshop-disposable"

const (
	defaultTimeout        = 10 * time.Second
	defaultMaxOutputBytes = 1 << 20 // 1 MiB
)

// Client runs the pinned Beans CLI against the configured workshop backlog.
type Client struct {
	noteMu      sync.Mutex
	cfg         Config
	beansVer    string
	writeArmed  bool
	writeReason string
}

// Open validates the configuration and records the Beans version. It fails fast so a
// misconfigured registration is reported at startup rather than on the first tool call.
func Open(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.Logger == nil {
		cfg.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultTimeout
	}
	if cfg.MaxOutputBytes <= 0 {
		cfg.MaxOutputBytes = defaultMaxOutputBytes
	}
	if cfg.Name == "" {
		cfg.Name = "workshop"
	}

	for label, p := range map[string]string{"beans executable": cfg.BeansBin, "beans config": cfg.ConfigFile, "beans data directory": cfg.DataDir} {
		if p == "" {
			return nil, newError(CodeBacklogUnavail, fmt.Sprintf("%s path is not configured", label), "")
		}
		if !filepath.IsAbs(p) {
			return nil, newError(CodeBacklogUnavail,
				fmt.Sprintf("%s path must be absolute", label), p)
		}
	}
	if info, err := os.Stat(cfg.BeansBin); err != nil || info.IsDir() {
		return nil, newError(CodeBacklogUnavail, "beans executable was not found", cfg.BeansBin)
	} else if runtime.GOOS != "windows" && info.Mode().Perm()&0o111 == 0 {
		return nil, newError(CodeBacklogUnavail, "beans executable is not executable", cfg.BeansBin)
	}
	if _, err := os.Stat(cfg.ConfigFile); err != nil {
		return nil, newError(CodeBacklogUnavail, "beans config file was not found", cfg.ConfigFile)
	}
	if info, err := os.Stat(cfg.DataDir); err != nil || !info.IsDir() {
		return nil, newError(CodeBacklogUnavail, "beans data directory was not found", cfg.DataDir)
	}

	c := &Client{cfg: cfg}

	out, err := c.run(ctx, []string{"version"})
	if err != nil {
		return nil, err
	}
	c.beansVer = strings.TrimSpace(string(out))

	if cfg.AllowNoteWrites {
		if _, statErr := os.Stat(filepath.Join(cfg.DataDir, DisposableMarker)); statErr == nil {
			c.writeArmed = true
		} else {
			c.writeReason = "the backlog is not marked disposable: " + DisposableMarker + " is missing"
			cfg.Logger.Warn("note writes requested but refused", "reason", c.writeReason)
		}
	} else {
		c.writeReason = "note writes were not enabled for this server"
	}
	return c, nil
}

// BeansVersion is the version string reported by the pinned CLI.
func (c *Client) BeansVersion() string { return c.beansVer }

// Name is the backlog label surfaced to clients.
func (c *Client) Name() string { return c.cfg.Name }

// WriteArmed reports whether the presenter write demonstration is actually available.
func (c *Client) WriteArmed() bool { return c.writeArmed }

// WriteRefusalReason explains why writes are unavailable.
func (c *Client) WriteRefusalReason() string { return c.writeReason }

// globalArgs are prepended to every invocation so the CLI can never fall back to
// upward config discovery or the BEANS_PATH environment variable.
func (c *Client) globalArgs() []string {
	return []string{"--config", c.cfg.ConfigFile, "--beans-path", c.cfg.DataDir}
}

// run executes the Beans CLI with an argument array, a timeout, and a bounded reader.
func (c *Client) run(ctx context.Context, args []string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.cfg.BeansBin, args...)
	cmd.Dir = c.cfg.DataDir
	// A deliberately minimal environment: no inherited BEANS_PATH, no credentials.
	cmd.Env = []string{"HOME=" + os.TempDir(), "NO_COLOR=1"}
	if runtime.GOOS == "windows" {
		// Windows has no Unix execute bits and uses different home/temp variables.
		// os/exec supplies SYSTEMROOT. Do not inherit credentials or BEANS_PATH.
		for _, key := range []string{"USERPROFILE", "APPDATA", "LOCALAPPDATA", "TEMP", "TMP"} {
			cmd.Env = append(cmd.Env, key+"="+os.TempDir())
		}
	} else {
		cmd.Env = append(cmd.Env, "PATH=/usr/bin:/bin")
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	stdout := &limitedWriter{w: &stdoutBuf, limit: c.cfg.MaxOutputBytes}
	stderr := &limitedWriter{w: &stderrBuf, limit: 64 << 10}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, newError(CodeBacklogTimeout,
			fmt.Sprintf("the backlog did not answer within %s", c.cfg.Timeout), strings.Join(args, " "))
	}
	if stdout.overflow {
		return nil, newError(CodeBacklogTooLarge, "the backlog returned more data than this tool will read", "")
	}
	if err != nil {
		// Beans writes a JSON error document to stdout for --json calls and exits 1.
		if code, msg, ok := decodeCLIError(stdoutBuf.Bytes()); ok {
			if code == "NOT_FOUND" {
				return nil, newError(CodeTaskNotFound, msg, "")
			}
			return nil, newError(CodeBacklogUnavail, msg, code)
		}
		return nil, newError(CodeBacklogUnavail, "the backlog command failed",
			fmt.Sprintf("%v: %s", err, strings.TrimSpace(stderrBuf.String())))
	}
	return stdoutBuf.Bytes(), nil
}

type cliError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code"`
}

func decodeCLIError(b []byte) (code string, message string, ok bool) {
	var doc cliError
	if err := json.Unmarshal(bytes.TrimSpace(b), &doc); err != nil {
		return "", "", false
	}
	if doc.Success || doc.Error == "" {
		return "", "", false
	}
	return doc.Code, doc.Error, true
}

// Get returns one task. The id must already have passed ValidateTaskID.
func (c *Client) Get(ctx context.Context, id string) (*Task, error) {
	if err := ValidateTaskID(id); err != nil {
		return nil, err
	}
	// Flags first, then "--", then the positional id: with the terminator in place a
	// value can never be reinterpreted as a flag.
	args := append(c.globalArgs(), "show", "--json", "--", id)
	out, err := c.run(ctx, args)
	if err != nil {
		return nil, err
	}

	var doc struct {
		ID        string   `json:"id"`
		Slug      string   `json:"slug"`
		Title     string   `json:"title"`
		Body      string   `json:"body"`
		Status    string   `json:"status"`
		Type      string   `json:"type"`
		Priority  string   `json:"priority"`
		Tags      []string `json:"tags"`
		BlockedBy []string `json:"blocked_by"`
		Blocking  []string `json:"blocking"`
		Parent    string   `json:"parent"`
		CreatedAt string   `json:"created_at"`
		UpdatedAt string   `json:"updated_at"`
		Etag      string   `json:"etag"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(out), &doc); err != nil {
		return nil, newError(CodeBacklogBadData, "the backlog returned output this tool could not parse", err.Error())
	}
	if doc.ID == "" {
		return nil, newError(CodeTaskNotFound, "task "+id+" was not found in the "+c.cfg.Name+" backlog", "")
	}
	body := strings.TrimSpace(doc.Body)
	return &Task{
		ID:                 doc.ID,
		Slug:               doc.Slug,
		Title:              doc.Title,
		Body:               body,
		Status:             doc.Status,
		Type:               doc.Type,
		Priority:           doc.Priority,
		Tags:               doc.Tags,
		BlockedBy:          doc.BlockedBy,
		Blocking:           doc.Blocking,
		Parent:             doc.Parent,
		CreatedAt:          doc.CreatedAt,
		UpdatedAt:          doc.UpdatedAt,
		Version:            doc.Etag,
		AcceptanceCriteria: ExtractAcceptanceCriteria(body),
		Backlog:            c.cfg.Name,
	}, nil
}

// ListResult is the bounded answer to list_tasks.
type ListResult struct {
	Tasks     []TaskSummary `json:"tasks"`
	Count     int           `json:"count"`
	Limit     int           `json:"limit"`
	Truncated bool          `json:"truncated"`
	Backlog   string        `json:"backlog"`
	Filter    struct {
		Status string `json:"status,omitempty"`
		Tag    string `json:"tag,omitempty"`
	} `json:"filter"`
}

// List returns a bounded list of tasks.
func (c *Client) List(ctx context.Context, filter ListFilter) (*ListResult, error) {
	f, err := ValidateListFilter(filter)
	if err != nil {
		return nil, err
	}
	args := append(c.globalArgs(), "list", "--json", "--sort", "id")
	if f.Status != "" {
		args = append(args, "--status", f.Status)
	}
	if f.Tag != "" {
		args = append(args, "--tag", f.Tag)
	}
	out, err := c.run(ctx, args)
	if err != nil {
		return nil, err
	}

	var docs []struct {
		ID        string   `json:"id"`
		Title     string   `json:"title"`
		Status    string   `json:"status"`
		Type      string   `json:"type"`
		Priority  string   `json:"priority"`
		Tags      []string `json:"tags"`
		BlockedBy []string `json:"blocked_by"`
		Etag      string   `json:"etag"`
	}
	trimmed := bytes.TrimSpace(out)
	if len(trimmed) == 0 {
		trimmed = []byte("[]")
	}
	if err := json.Unmarshal(trimmed, &docs); err != nil {
		return nil, newError(CodeBacklogBadData, "the backlog returned output this tool could not parse", err.Error())
	}

	result := &ListResult{Limit: f.Limit, Backlog: c.cfg.Name}
	result.Filter.Status = f.Status
	result.Filter.Tag = f.Tag
	for _, d := range docs {
		if len(result.Tasks) == f.Limit {
			result.Truncated = true
			break
		}
		result.Tasks = append(result.Tasks, TaskSummary{
			ID: d.ID, Title: d.Title, Status: d.Status, Type: d.Type,
			Priority: d.Priority, Tags: d.Tags, BlockedBy: d.BlockedBy, Version: d.Etag,
		})
	}
	result.Count = len(result.Tasks)
	return result, nil
}

// AddNote appends a note to a task. Presenter demonstration only: it refuses unless
// writes were enabled AND the data directory carries the disposable marker.
func (c *Client) AddNote(ctx context.Context, id, note string) (*Task, error) {
	c.noteMu.Lock()
	defer c.noteMu.Unlock()
	if !c.writeArmed {
		reason := c.writeReason
		if reason == "" {
			reason = "note writes are disabled"
		}
		return nil, newError(CodeWriteDisabled, "this backlog is read-only: "+reason, "")
	}
	if err := ValidateTaskID(id); err != nil {
		return nil, err
	}
	note = strings.TrimSpace(note)
	if note == "" {
		return nil, newError(CodeInvalidFilter, "note must not be empty", "")
	}
	if utf8.RuneCountInString(note) > 2000 {
		return nil, newError(CodeInvalidFilter, "note must be 2000 characters or fewer", "")
	}
	current, err := c.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if info, err := os.Stat(filepath.Join(c.cfg.DataDir, DisposableMarker)); err != nil || !info.Mode().IsRegular() {
		return nil, newError(CodeWriteDisabled, "this backlog is read-only: the disposable marker is missing", "")
	}
	body := "\n\n## Note (presenter demonstration)\n" + note + "\n"
	args := append(c.globalArgs(), "update", "--json", "--body-append", body, "--if-match", current.Version, "--", id)
	if _, err := c.run(ctx, args); err != nil {
		return nil, err
	}
	return c.Get(ctx, id)
}

// ExtractAcceptanceCriteria pulls the acceptance-criteria bullet list out of a task
// body. The body stays authoritative; this is a convenience for the agent.
func ExtractAcceptanceCriteria(body string) []string {
	lines := strings.Split(body, "\n")
	var out []string
	inSection := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(trimmed, "#") {
			if strings.TrimSpace(strings.TrimLeft(lower, "#")) == "acceptance criteria" {
				inSection = true
				continue
			}
			if inSection {
				break
			}
			continue
		}
		if !inSection {
			continue
		}
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			out = append(out, strings.TrimSpace(trimmed[2:]))
			continue
		}
		if trimmed == "" {
			continue
		}
		if len(out) > 0 {
			// A paragraph after the bullets ends the list.
			break
		}
	}
	return out
}

// limitedWriter stops writing after limit bytes and records the overflow.
type limitedWriter struct {
	w        *bytes.Buffer
	limit    int64
	written  int64
	overflow bool
}

func (l *limitedWriter) Write(p []byte) (int, error) {
	if l.written+int64(len(p)) > l.limit {
		l.overflow = true
		remaining := l.limit - l.written
		if remaining > 0 {
			l.w.Write(p[:remaining])
			l.written = l.limit
		}
		return len(p), nil
	}
	n, err := l.w.Write(p)
	l.written += int64(n)
	return n, err
}

// FormatTimeout renders a timeout for diagnostics.
func FormatTimeout(d time.Duration) string { return strconv.FormatFloat(d.Seconds(), 'f', 1, 64) + "s" }

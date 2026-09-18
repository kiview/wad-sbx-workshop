package backlog

import "testing"

func TestValidateTaskID(t *testing.T) {
	valid := []string{"wad-101", "wad-a1b2", "workshop-task-01", "a-b", "wad-101-extra-bits"}
	for _, id := range valid {
		if err := ValidateTaskID(id); err != nil {
			t.Errorf("ValidateTaskID(%q) = %v, want nil", id, err)
		}
	}
	invalid := []string{
		"", "wad", "WAD-101", "wad_101", "wad 101", "-wad-101", "wad-101-",
		"../etc/passwd", "wad-101;id", "$(id)", "wad-101|cat", "--beans-path",
		"a-b-c-d-e", "wad-101\n", "wad-101/../other",
	}
	for _, id := range invalid {
		if err := ValidateTaskID(id); err == nil {
			t.Errorf("ValidateTaskID(%q) = nil, want an error", id)
		}
	}
}

func TestValidateTaskIDLength(t *testing.T) {
	long := ""
	for len(long) < 70 {
		long += "ab-"
	}
	if err := ValidateTaskID(long); err == nil {
		t.Error("an over-long id was accepted")
	}
}

func TestValidateListFilter(t *testing.T) {
	got, err := ValidateListFilter(ListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if got.Limit != DefaultListLimit {
		t.Errorf("default limit = %d, want %d", got.Limit, DefaultListLimit)
	}

	if _, err := ValidateListFilter(ListFilter{Status: "todo", Tag: "workshop", Limit: 10}); err != nil {
		t.Errorf("a valid filter was rejected: %v", err)
	}
	for _, f := range []ListFilter{
		{Status: "sideways"},
		{Status: "--beans-path"},
		{Tag: "Workshop"},
		{Tag: "work shop"},
		{Tag: "-workshop"},
		{Limit: MaxListLimit + 1},
		{Limit: -3},
	} {
		if _, err := ValidateListFilter(f); err == nil {
			t.Errorf("filter %+v was accepted", f)
		}
	}
}

func TestExtractAcceptanceCriteria(t *testing.T) {
	body := `Some intro text.

## Acceptance criteria
- first rule
- second rule
* third rule

## Out of scope
- not this
`
	got := ExtractAcceptanceCriteria(body)
	want := []string{"first rule", "second rule", "third rule"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("criterion %d = %q, want %q", i, got[i], want[i])
		}
	}

	if got := ExtractAcceptanceCriteria("No sections here."); len(got) != 0 {
		t.Errorf("expected no criteria, got %v", got)
	}
	if got := ExtractAcceptanceCriteria("## Acceptance Criteria\n- case insensitive\n"); len(got) != 1 {
		t.Errorf("heading match should be case-insensitive, got %v", got)
	}
}

func TestErrorFormatting(t *testing.T) {
	err := newError(CodeInvalidID, "safe message", "host detail")
	if got := err.Error(); got != "invalid_id: safe message (host detail)" {
		t.Errorf("Error() = %q", got)
	}
	// The safe message must not carry the detail: the MCP layer sends Message only.
	if err.Message != "safe message" {
		t.Errorf("Message = %q", err.Message)
	}
}

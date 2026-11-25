// PostToolUse hook for edited Go files: gofmt the file, then run the forbidden-usage linters (depguard,
// forbidigo) on its package in api/.
//
// Problems are reported back to Claude with a "block" decision on stdout, so they are fixed in the same
// turn instead of surfacing only in CI. Anything else (golangci-lint missing, compile errors mid-edit) is
// ignored. The JSON output is used instead of exit code 2 because `go run` turns any non-zero exit into 1.
// Standard library only, so it runs with `go run` outside the api module.
package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type hookOutput struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

type hookInput struct {
	ToolInput struct {
		FilePath string `json:"file_path"`
	} `json:"tool_input"`
}

func main() {
	var in hookInput
	if err := json.NewDecoder(os.Stdin).Decode(&in); err != nil {
		return
	}

	path, err := filepath.Abs(in.ToolInput.FilePath)
	if err != nil || !strings.HasSuffix(path, ".go") {
		return
	}

	// A syntax error is expected while an edit is split across several steps; the compiler reports it later.
	if err := exec.Command("gofmt", "-w", path).Run(); err != nil {
		return
	}

	if violations := forbiddenUsages(path); len(violations) > 0 {
		block("Forbidden usages (see depguard and forbidigo in api/.golangci.yml):\n" + strings.Join(violations, "\n"))
	}
}

func block(reason string) {
	_ = json.NewEncoder(os.Stdout).Encode(hookOutput{Decision: "block", Reason: reason})
}

// forbiddenUsages runs only the fast linters that guard project conventions, to keep the hook under a second.
func forbiddenUsages(path string) []string {
	projectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	if projectDir == "" {
		return nil
	}

	apiDir := filepath.Join(projectDir, "api")
	rel, err := filepath.Rel(apiDir, filepath.Dir(path))
	if err != nil || strings.HasPrefix(rel, "..") {
		return nil
	}

	if _, err := exec.LookPath("golangci-lint"); err != nil {
		return nil
	}

	cmd := exec.Command("golangci-lint", "run", "--enable-only", "depguard,forbidigo", "./"+filepath.ToSlash(rel))
	cmd.Dir = apiDir
	out, _ := cmd.Output() // exits 1 when there are issues; the output is what matters

	var violations []string
	for line := range strings.Lines(string(out)) {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, "(depguard)") || strings.HasSuffix(line, "(forbidigo)") {
			violations = append(violations, line)
		}
	}

	return violations
}

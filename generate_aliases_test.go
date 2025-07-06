package main

import (
	"bufio"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

var (
	aliasPattern    = regexp.MustCompile(`(?m)^alias (\w+)=['"](kubectl.*?)['"]`)
	functionPattern = regexp.MustCompile(`(?m)^function (\w+)\(\) \{ (kubectl.*?) \$1; \}`)
)

type aliasEntry struct {
	Name    string
	Command string
	IsFunc  bool
}

func parseAliasFile(path string) ([]aliasEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []aliasEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if matches := aliasPattern.FindStringSubmatch(line); len(matches) == 3 {
			entries = append(entries, aliasEntry{Name: matches[1], Command: matches[2], IsFunc: false})
		} else if matches := functionPattern.FindStringSubmatch(line); len(matches) == 3 {
			entries = append(entries, aliasEntry{Name: matches[1], Command: matches[2], IsFunc: true})
		}
	}

	return entries, scanner.Err()
}

func TestAliasFileCommands(t *testing.T) {
	entries, err := parseAliasFile("kubectl_aliases")
	if err != nil {
		t.Fatalf("Failed to parse alias file: %v", err)
	}

	for _, entry := range entries {
		// Skip aliases that map to 'kubectl create' for now
		if strings.Contains(entry.Command, "kubectl create") {
			t.Logf("Skipping 'create' alias %q - TODO: Add support for dynamic parameters.", entry.Name)
			continue
		}

		// Skip aliases that map to 'kubectl edit' to avoid interactive terminal hang
		if strings.Contains(entry.Command, "kubectl edit") {
			t.Logf("Skipping 'edit' alias %q - command is interactive.", entry.Name)
			continue
		}

		t.Run(entry.Name, func(t *testing.T) {
			cmdStr := entry.Command
			if entry.IsFunc {
				cmdStr = strings.Replace(cmdStr, "$1", "nonexistent-test", -1)
			}

			// Handle known special cases
			switch entry.Name {
			case "kl": // logs
				cmdStr += " -n kube-system -l component=kube-apiserver"
			case "ke": // explain
				cmdStr += " pods"
			case "kg": // get
				cmdStr += " -n kube-system pods"
			}

			t.Logf("Running command: %s", cmdStr)

			parts := strings.Fields(cmdStr)
			if len(parts) == 0 {
				t.Skip("No command to run")
			}

			cmd := exec.Command(parts[0], parts[1:]...)
			out, err := cmd.CombinedOutput()
			stdout := string(out)

			t.Logf("Command output:\n%s", stdout)

			if err == nil {
				return // passed
			}

			if strings.Contains(stdout, "Unexpected args") ||
				strings.Contains(stdout, "unknown flag") ||
				strings.Contains(stdout, "unknown command") ||
				strings.Contains(stdout, "unknown shorthand flag") ||
				strings.Contains(stdout, "See 'kubectl") {
				t.Errorf("Command %q failed with bad args: %v\nOutput:\n%s", cmdStr, err, stdout)
			}
		})
	}
}

package main

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
)

const (
	aliasFile   = "kubectl_aliases"
	manifest    = "kubectl_test_objects.yaml"
	testName    = "test-pod"
	parallelism = 20
)

func TestMain(m *testing.M) {
	// Setup test resources
	if err := applyTestManifests(); err != nil {
		panic("Failed to apply test manifests: " + err.Error())
	}

	// Simplified wait logic for now; improve in future
	// time.Sleep(60 * time.Second)

	code := m.Run()
	os.Exit(code)
}

func applyTestManifests() error {
	cmd := exec.Command("kubectl", "apply", "-f", manifest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommandWithOutput(cmdStr string) (string, error) {
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	combined := outBuf.String() + errBuf.String()
	return combined, err
}

func cleanAliasCommand(cmdRaw string) string {
	// Remove inline comments and surrounding quotes
	cmd := strings.SplitN(cmdRaw, "#", 2)[0]
	cmd = strings.TrimSpace(cmd)
	cmd = strings.Trim(cmd, "'")
	return cmd
}

func TestAliasCommands(t *testing.T) {
	file, err := os.Open(aliasFile)
	if err != nil {
		t.Fatalf("Failed to open %s: %v", aliasFile, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sem := make(chan struct{}, parallelism)
	var wg sync.WaitGroup

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, " get ") && !strings.Contains(line, " describe ") {
			continue
		}

		var validCmd, invalidCmd string

		if strings.HasPrefix(line, "alias") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) < 2 {
				continue
			}
			cmd := cleanAliasCommand(parts[1])
			if strings.Contains(cmd, "pod") || strings.Contains(cmd, "svc") || strings.Contains(cmd, "configmap") {
				validCmd = cmd + " " + testName
				invalidCmd = cmd + " missing-pod"
			} else {
				validCmd = cmd
			}
		} else if strings.Contains(line, "() {") {
			start := strings.Index(line, "{")
			end := strings.LastIndex(line, "}")
			if start == -1 || end == -1 {
				continue
			}
			body := strings.TrimSpace(line[start+1 : end])
			body = cleanAliasCommand(body)
			validCmd = strings.ReplaceAll(body, "$1", testName)
			invalidCmd = strings.ReplaceAll(body, "$1", "missing-pod")
		}

		if validCmd != "" {
			sem <- struct{}{}
			wg.Add(1)
			go func(cmd string) {
				defer func() {
					<-sem
					wg.Done()
				}()
				t.Run("valid:"+cmd, func(t *testing.T) {
					t.Parallel()
					output, err := runCommandWithOutput(cmd)
					if err != nil {
						if strings.Contains(output, "unexpected EOF") || strings.Contains(output, "syntax error") {
							t.Fatalf("Invalid shell syntax in alias: %s\nOutput:\n%s", cmd, output)
						} else {
							t.Logf("Expected failure for unsupported or non-existent object: %s\nOutput:\n%s", cmd, output)
						}
					}
				})
			}(validCmd)
		}

		if invalidCmd != "" {
			sem <- struct{}{}
			wg.Add(1)
			go func(cmd string) {
				defer func() {
					<-sem
					wg.Done()
				}()
				t.Run("invalid:"+cmd, func(t *testing.T) {
					t.Parallel()
					output, err := runCommandWithOutput(cmd)
					if err == nil {
						t.Errorf("Expected failure, but command succeeded: %s\nOutput:\n%s", cmd, output)
					} else if strings.Contains(output, "unexpected EOF") || strings.Contains(output, "syntax error") {
						t.Errorf("Invalid shell syntax in alias: %s\nOutput:\n%s", cmd, output)
					}
				})
			}(invalidCmd)
		}
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading alias file: %v", err)
	}
}

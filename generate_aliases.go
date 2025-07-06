package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

// Modifier definitions for common CLI flags
var OutputModifiers = map[string]string{
	"j": "-o json",
	"y": "-o yaml",
	"w": "-o wide",
}

var InputModifiers = map[string]string{
	"f": "-f", // file input modifier
	"":  "",   // default input (e.g. stdin)
}

var ScopeModifiers = map[string]string{
	"n": "--all-namespaces",
}

// Common CRUD commands for Kubernetes resources
var CRUDCommands = []string{"get", "edit", "create", "delete", "describe"}

// Command-specific configuration including priority and supported modifiers
var CommandConfig = map[string]CommandSpec{
	"apply":    {Priority: true, Modifiers: map[string]map[string]string{"input": InputModifiers}},
	"create":   {Priority: true},
	"describe": {Priority: true, Modifiers: map[string]map[string]string{"scope": ScopeModifiers}},
	"edit":     {Priority: true},
	"get":      {Priority: true, Modifiers: map[string]map[string]string{"output": OutputModifiers, "scope": ScopeModifiers}},
	"logs":     {Priority: true},
	"top":      {Priority: true},
	"auth":     {Modifiers: map[string]map[string]string{"scope": ScopeModifiers}},
	"debug":    {Modifiers: map[string]map[string]string{"scope": ScopeModifiers}},
	"events":   {Modifiers: map[string]map[string]string{"output": OutputModifiers}},
	"delete":   {Modifiers: map[string]map[string]string{"input": InputModifiers}},
}

type CommandSpec struct {
	Priority  bool
	Modifiers map[string]map[string]string
}

// Loads the Kubernetes discovery client for dynamic resource introspection
func getDiscoveryClient() (*discovery.DiscoveryClient, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = clientcmd.RecommendedHomeFile
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return discovery.NewDiscoveryClientForConfig(config)
}

// Discovers all available API resources and their short names
func getAPIResources() map[string]string {
	resources := map[string]string{}
	client, err := getDiscoveryClient()
	if err != nil {
		return resources
	}
	apiLists, err := client.ServerPreferredResources()
	if err != nil {
		return resources
	}
	for _, list := range apiLists {
		for _, res := range list.APIResources {
			if len(res.ShortNames) > 0 {
				resources[res.ShortNames[0]] = res.Name
			}
		}
	}
	return resources
}

// Extracts available kubectl commands from the CLI help output
func getKubectlCommands() []string {
	out, err := exec.Command("kubectl", "--help").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing kubectl --help: %v\n", err)
		return nil
	}

	cmds := map[string]bool{}
	re := regexp.MustCompile(`^\s{2}([a-z0-9-]+)\s{2,}`)
	for _, line := range strings.Split(string(out), "\n") {
		if matches := re.FindStringSubmatch(line); matches != nil {
			cmds[matches[1]] = true
		}
	}

	var result []string
	for cmd := range cmds {
		result = append(result, cmd)
	}
	sort.Strings(result)
	return result
}

// Removes dashes from kubectl subcommands
func sanitizeCommand(cmd string) string {
	return strings.ReplaceAll(cmd, "-", "")
}

// Generates alias names for supported kubectl commands
func generateAliases(cmdConfig map[string]CommandSpec, dynamicCmds []string) map[string]string {
	used := map[string]string{}
	aliases := map[string]string{}

	var priorityList []string
	for cmd := range cmdConfig {
		if cmdConfig[cmd].Priority {
			priorityList = append(priorityList, cmd)
		}
	}
	sort.Strings(priorityList)
	for _, cmd := range priorityList {
		sanitized := sanitizeCommand(cmd)
		short := string(sanitized[0])
		if _, taken := used[short]; !taken {
			aliases[cmd] = short
			used[short] = cmd
		}
	}

	sorted := append([]string{}, dynamicCmds...)
	sort.Strings(sorted)
	for _, cmd := range sorted {
		if _, ok := aliases[cmd]; ok {
			continue
		}
		if _, valid := CommandConfig[cmd]; !valid {
			continue
		}
		sanitized := sanitizeCommand(cmd)
		alias := ""
		for i := 2; i <= len(sanitized); i++ {
			prefix := sanitized[:i]
			if _, taken := used[prefix]; !taken {
				alias = prefix
				break
			}
		}
		if alias == "" {
			alias = sanitized
		}
		aliases[cmd] = alias
		used[alias] = cmd
	}

	return aliases
}

// Writes sorted alias and function definitions to a file
func writeFunctionFile(aliases map[string]string, path string, resources map[string]string) error {
	var lines []string
	lines = append(lines, "# Auto-generated kubectl function alias file")

	for _, cmd := range sortedKeys(aliases) {
		alias := aliases[cmd]
		lines = append(lines, fmt.Sprintf("alias k%s='kubectl %s' # %s", alias, cmd, cmd))

		spec := CommandConfig[cmd]

		// Generate aliases for CRUD operations over resource short names
		if contains(CRUDCommands, cmd) {
			for _, pair := range sortedMap(resources) {
				short := pair[0]
				full := pair[1]
				lines = append(lines, fmt.Sprintf("alias k%s%s='kubectl %s %s' # %s %s", alias, short, cmd, short, cmd, full))

				// Generate function-style aliases for modifiers (e.g. output)
				if cmd == "get" {
					for suffix, flag := range OutputModifiers {
						lines = append(lines, fmt.Sprintf("function k%s%s%s() { kubectl %s %s \"$1\" %s; } # %s %s + output", alias, short, suffix, cmd, short, flag, cmd, full))
					}
				}
			}
		}

		// Append aliases for command-specific modifiers (e.g. -f, --all-namespaces)
		for modType, mods := range spec.Modifiers {
			for suffix, flag := range mods {
				lines = append(lines, fmt.Sprintf("function k%s%s() { kubectl %s %s \"$1\"; } # %s + %s modifier", alias, suffix, cmd, flag, cmd, modType))
			}
		}
	}

	sort.Strings(lines) // Ensure deterministic order
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

// Updates the README file with alias information and generated reference block
// TODO: Make managing \n chars in this func below easier
func updateReadme(aliases map[string]string, readmePath string) error {
	data, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}

	// Generate alias block
	var aliasLines []string
	aliasLines = append(aliasLines, "## All Available Aliases", "\n```")
	for _, cmd := range sortedKeys(aliases) {
		aliasLines = append(aliasLines, fmt.Sprintf("%s => k%s", cmd, aliases[cmd]))
	}
	aliasLines = append(aliasLines, "```")
	aliasBlock := strings.Join(aliasLines, "\n")

	// Read kubectl_aliases and trim *only* final newline if present
	fullContent, err := os.ReadFile("kubectl_aliases")
	if err != nil {
		return err
	}
	trimmedFullContent := strings.TrimSuffix(string(fullContent), "\n")
	refBlock := "\n\n## Full Alias Reference\n\n```bash\n" + trimmedFullContent + "\n```\n"

	// Build new README content
	split := strings.Split(string(data), "## All Available Aliases")
	var newReadme string
	if len(split) > 1 {
		newReadme = split[0] + aliasBlock + refBlock
	} else {
		newReadme = string(data) + "\n" + aliasBlock + refBlock
	}

	return os.WriteFile(readmePath, []byte(newReadme), 0644)
}

// sortedMap returns sorted key-value pairs from a map[string]string
func sortedMap(m map[string]string) [][2]string {
	var pairs [][2]string
	for k, v := range m {
		pairs = append(pairs, [2]string{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][0] < pairs[j][0]
	})
	return pairs
}

// Utility function to get sorted keys from a map
func sortedKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Checks if a string exists in a slice
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func main() {
	resources := getAPIResources()
	dynamicCommands := getKubectlCommands()
	aliases := generateAliases(CommandConfig, dynamicCommands)
	_ = writeFunctionFile(aliases, "kubectl_aliases", resources)
	_ = updateReadme(aliases, "README.md")
}

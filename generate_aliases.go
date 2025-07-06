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

var OutputModifiers = map[string]string{
	"j": "-o json",
	"y": "-o yaml",
	"l": "-o wide",
}

var InputModifiers = map[string]string{
	"f": "-f", // file input modifier
	"":  "",   // default input (e.g. stdin)
}

var ScopeModifiers = map[string]string{
	"n": "--all-namespaces",
}

var CRUDCommands = []string{"get", "edit", "create", "delete", "describe"}

var CommandConfig = map[string]CommandSpec{
	"apply":    {Priority: true, Modifiers: map[string]map[string]string{"input": InputModifiers}},
	"create":   {Priority: true},
	"describe": {Priority: true, Modifiers: map[string]map[string]string{"scope": ScopeModifiers}},
	"edit":     {Priority: true},
	"get":      {Priority: true, Modifiers: map[string]map[string]string{"output": OutputModifiers, "scope": ScopeModifiers}},
	"logs":     {Priority: true, Modifiers: map[string]map[string]string{"output": OutputModifiers, "scope": ScopeModifiers}},
	"top":      {Priority: true, Modifiers: map[string]map[string]string{"scope": ScopeModifiers}},
	"auth":     {Modifiers: map[string]map[string]string{"scope": ScopeModifiers}},
	"debug":    {Modifiers: map[string]map[string]string{"scope": ScopeModifiers}},
	"events":   {Modifiers: map[string]map[string]string{"output": OutputModifiers, "scope": ScopeModifiers}},
	"delete":   {Modifiers: map[string]map[string]string{"input": InputModifiers}},
}

type CommandSpec struct {
	Priority  bool
	Modifiers map[string]map[string]string
}

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

func sanitizeCommand(cmd string) string {
	return strings.ReplaceAll(cmd, "-", "")
}

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

func writeFunctionFile(aliases map[string]string, path string, resources map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("# Auto-generated kubectl function alias file\n")
	for _, cmd := range sortedKeys(aliases) {
		alias := aliases[cmd]
		f.WriteString(fmt.Sprintf("alias k%s='kubectl %s' # %s\n", alias, cmd, cmd))

		spec := CommandConfig[cmd]
		if contains(CRUDCommands, cmd) {
			for short, full := range resources {
				f.WriteString(fmt.Sprintf("alias k%s%s='kubectl %s %s' # %s %s\n", alias, short, cmd, short, cmd, full))
				if cmd == "get" {
					for suffix, flag := range OutputModifiers {
						f.WriteString(fmt.Sprintf("function k%s%s%s() { kubectl %s %s \"$1\" %s; } # %s %s + output\n", alias, short, suffix, cmd, short, flag, cmd, full))
					}
				}
			}
		}

		for modType, mods := range spec.Modifiers {
			for suffix, flag := range mods {
				f.WriteString(fmt.Sprintf("function k%s%s() { kubectl %s %s \"$1\"; } # %s + %s modifier\n", alias, suffix, cmd, flag, cmd, modType))
			}
		}
	}
	return nil
}

func updateReadme(aliases map[string]string, readmePath string) error {
	data, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}

	aliasLines := []string{"\n## All Available Aliases\n", "```"}
	for _, cmd := range sortedKeys(aliases) {
		alias := aliases[cmd]
		aliasLines = append(aliasLines, fmt.Sprintf("%s => k%s", cmd, alias))
	}
	aliasLines = append(aliasLines, "```")
	aliasBlock := strings.Join(aliasLines, "\n")

	fullContent, err := os.ReadFile("kubectl_aliases")
	if err != nil {
		return err
	}
	refBlock := "\n## Full Alias Reference\n```bash\n" + string(fullContent) + "\n```"

	split := strings.Split(string(data), "## All Available Aliases")
	var newReadme string
	if len(split) > 1 {
		newReadme = split[0] + aliasBlock + refBlock
	} else {
		newReadme = string(data) + "\n" + aliasBlock + refBlock + "\n"
	}

	return os.WriteFile(readmePath, []byte(newReadme), 0644)
}

func main() {
	resources := getAPIResources()
	dynamicCommands := getKubectlCommands()
	aliases := generateAliases(CommandConfig, dynamicCommands)
	_ = writeFunctionFile(aliases, "kubectl_aliases", resources)
	_ = updateReadme(aliases, "README.md")
}

func sortedKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	pluginName  = "find-app"
	NoticeColor = "\033[1;36m%s\033[0m"
)

// FindAppPlugin is the struct implementing the interface defined by the core CLI. It can
// be found at  "code.cloudfoundry.org/cli/plugin/plugin.go"
type FindAppPlugin struct{}
type CFApp struct {
	Name  string
	Guid  string
	State string
}
type CFApps struct {
	CFApps []*CFApp `json:"resources"`
}

// Run must be implemented by any plugin because it is part of the
// plugin interface defined by the core CLI.
//
// Run(....) is the entry point when the core CLI is invoking a command defined
// by a plugin. The first parameter, plugin.CliConnection, is a struct that can
// be used to invoke cli commands. The second parameter, args, is a slice of
// strings. args[0] will be the name of the command, and will be followed by
// any additional arguments a cli user typed in.
//
// Any error handling should be handled with the plugin itself (this means printing
// user facing errors). The CLI will exit 0 if the plugin exits 0 and will exit
// 1 should the plugin exits nonzero.
func (c *FindAppPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	startTime := time.Now().UnixMilli()
	c.validateInput(args)
	c.process(cliConnection, args[1])
	timeTaken := float32(time.Now().UnixMilli()-startTime) / 1000
	fmt.Printf("\nTook %vs\n", timeTaken)
}

func (c *FindAppPlugin) process(cliConnection plugin.CliConnection, searchString string) {
	space, err := cliConnection.GetCurrentSpace()
	if err != nil {
		c.logError("Unable to get current space due to ", err)
		return
	}
	url := fmt.Sprintf("/v3/apps?space_guids=%v&per_page=5000", space.Guid)
	lines, err := cliConnection.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		c.logError("Unable to execute cf curl "+url, err)
		return
	}
	result := strings.Join(lines, "")
	cfAppsWrapper, err := c.parseCFApps(result)
	if err != nil {
		c.logError("Json parse error for result "+result, err)
		return
	}
	for _, app := range cfAppsWrapper.CFApps {
		if strings.Contains(app.Name, searchString) {
			fmt.Printf("[Name:%q, State:%q]\n", app.Name, app.State)
		}
	}
}

func (c *FindAppPlugin) parseCFApps(inputString string) (*CFApps, error) {
	var cfApps CFApps
	err := json.Unmarshal([]byte(inputString), &cfApps)
	if err != nil {
		c.logError("Unable to parse response ", err)
		return &CFApps{}, err
	}
	return &cfApps, nil
}

func (c *FindAppPlugin) logError(prefix string, err error) {
	fmt.Println(prefix)
	fmt.Println(err)
}

func (c *FindAppPlugin) validateInput(args []string) {
	if len(args) != 2 {
		fmt.Println("Invalid usage. Please use run 'cf find-app -h' to get usage info")
		os.Exit(0)
	}
}

// GetMetadata must be implemented as part of the plugin interface
// defined by the core CLI.
//
// GetMetadata() returns a PluginMetadata struct. The first field, Name,
// determines the name of the plugin which should generally be without spaces.
// If there are spaces in the name a user will need to properly quote the name
// during uninstall otherwise the name will be treated as seperate arguments.
// The second value is a slice of Command structs. Our slice only contains one
// Command Struct, but could contain any number of them. The first field Name
// defines the command `cf basic-plugin-command` once installed into the CLI. The
// second field, HelpText, is used by the core CLI to display help information
// to the user in the core commands `cf help`, `cf`, or `cf -h`.
func (c *FindAppPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "FindAppPlugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name: pluginName,
				HelpText: "Searches for the applications in currently logged in space with name containing the search string.\n" +
					fmt.Sprintf(NoticeColor, "   Prerequisite: User should be logged into the required space."),

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "cf find-app <search_string>",
				},
			},
		},
	}
}

// Unlike most Go programs, the `Main()` function will not be used to run all of the
// commands provided in your plugin. Main will be used to initialize the plugin
// process, as well as any dependencies you might require for your
// plugin.
func main() {
	// Any initialization for your plugin can be handled here
	//
	// Note: to run the plugin.Start method, we pass in a pointer to the struct
	// implementing the interface defined at "code.cloudfoundry.org/cli/plugin/plugin.go"
	//
	// Note: The plugin's main() method is invoked at install time to collect
	// metadata. The plugin will exit 0 and the Run([]string) method will not be
	// invoked.
	plugin.Start(new(FindAppPlugin))
	// Plugin code should be written in the Run([]string) method,
	// ensuring the plugin environment is bootstrapped.
}

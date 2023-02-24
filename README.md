# Overview
This cf cli plugin allows users to search for application names that contain the given search string. 
This plugin is implemented following the basic plugin provided by [cloud foundry documentation](https://docs.cloudfoundry.org/cf-cli/develop-cli-plugins.html).

# Steps for installation
1. Install Go 
2. Clone plugin source code 
----
git clone https://github.com/jaisonjose89m/cf-cli-find-app-plugin.git
----
3. Build plugin 
   1. Move to `cf-cli-find-app-plugin` folder and execute `go build .`
4. Install plugin `cf install-plugin <plugin_binary_path> -f`

# Usage guide
```shell
cf find-app <search_string>
```
Note: `search_string` should be a simple string. The plugin is not able to perform pattern matching with regex yet.

Example:
```shell
cf find-app webservice
```
Output:
```console
[Name:"abcwebservice", State:"STARTED"]
[Name:"webserviceoq", State:"STARTED"]
```

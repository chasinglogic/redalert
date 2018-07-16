package checks

import "fmt"

var availableChecks = map[string]ArgableFunc{
	"file-exists": func(args map[string]interface{}) (Checker, error) {
		return FileChecker{}.FromArgs(args)
	},
	"file-does-not-exist": func(args map[string]interface{}) (Checker, error) {
		args["exists"] = false
		return FileChecker{}.FromArgs(args)
	},
	"run-script": func(args map[string]interface{}) (Checker, error) {
		return RunScript{}.FromArgs(args)
	},
	"run-bash-script": func(args map[string]interface{}) (Checker, error) {
		args["interpreter"] = "/bin/bash"
		return RunScript{}.FromArgs(args)
	},
}

// LoadCheck will return the appropriate Checker based on the test type name.
// As documented on the various checkers
func LoadCheck(name string, args map[string]interface{}) (Checker, error) {
	if argFunc, exists := availableChecks[name]; exists {
		return argFunc(args)
	}

	return nil, fmt.Errorf("%s is not a known check type", name)
}

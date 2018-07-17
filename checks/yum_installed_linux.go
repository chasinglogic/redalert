package checks

import (
	"fmt"
	"log"
	"os/exec"
)

func init() {
	availableChecks["yum-installed"] = func(args map[string]interface{}) (Checker, error) {
		return YumInstalled{}.FromArgs(args)
	}
}

// YumInstalled checks if an rpm package is installed on the system
//
// Type:
//   - yum-installed
//
// Support Platforms:
//   - Linux
//
// Arguments:
//   name (required): A string value that represents the rpm package
type YumInstalled struct {
	Name string
}

// Check if an rpm is installed on the system
func (yi YumInstalled) Check() error {
	out, err := exec.Command("rpm", "-qa", yi.Name).Output()
	if err != nil {
		log.Fatal(err)
	}

	if len(out) <= 0 {
		return fmt.Errorf("%s isn't installed and should be", yi.Name)
	}

	return nil
}

// FromArgs will populate the YumInstalled struct with the args given in the tests YAML
// config
func (yi YumInstalled) FromArgs(args map[string]interface{}) (Checker, error) {
	if err := requiredArgs(args, "name"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &yi); err != nil {
		return nil, err
	}

	return yi, nil
}

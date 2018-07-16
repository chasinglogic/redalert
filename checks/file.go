package checks

import (
	"fmt"
	"os"
)

// FileChecker checks if a file exists or does not
//
// Type:
//	 - file-exists
//   - file-does-not-exist
//
// Support Platforms:
//   - Mac
//   - Linux
//   - Windows
//
// Arguments:
//   name (required): A string value that points to a path on the filesystem.
//   exists: An optional boolean indicating whether the file should exist or not.
//			 For file-does-not-exist type tests this is always set to false, for
//			 the normal file-exists type tests this value defaults to true, the
//			 file should exist, but can be set to false if desired.
//
// Notes:
//   For Unix systems no `~` expansion is done. So ~/.bashrc is not a valid path,
//   or at least will not do what you think it will. Additionally, when checking
//   paths on Windows provide windows style paths (i.e. C:\My\File\Path.txt).
type FileChecker struct {
	Exists bool
	Name   string
}

// Check if a file exists or does not
func (fc FileChecker) Check() error {
	_, err := os.Stat(fc.Name)

	switch os.IsNotExist(err) {
	case true && fc.Exists:
		return fmt.Errorf("%s doesn't exist and should", fc.Name)
	case false && !fc.Exists:
		return fmt.Errorf("%s does exist and shouldn't", fc.Name)
	default:
		return nil
	}
}

// FromArgs will populate the FileChecker with the args given in the tests YAML
// config
func (fc FileChecker) FromArgs(args map[string]interface{}) (Checker, error) {
	if err := requiredArgs(args, "name"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &fc); err != nil {
		return nil, err
	}

	if _, existsGiven := args["exists"]; !fc.Exists && !existsGiven {
		fc.Exists = true
	}

	return fc, nil
}

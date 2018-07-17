package checks

import (
	"fmt"
	"syscall"
)

// UlimitChecker checks if current process resource limits are above a given minimum
//
// Type:
//   - ulimit
//
// Supported Platforms:
//   - MacOS
//   - Linux
//
// Arguments:
//   - item (required): A string value representing the type of limit to check
//   - limit (required): Numerical value representing the minimum value to be tested
//   - type (optional): "hard" or "soft" with a default of "hard"
//
// Notes:
//   - Windows is not supported
//   - "item" strings are from http://www.linux-pam.org/Linux-PAM-html/sag-pam_limits.html
//     - Not all are supported. See `limitsByName` map below for full list
//   - "limit" can be '-1' to represent that the resource limit should be unlimited
//
type UlimitChecker struct {
	Item   string
	Limit  uint64
	IsHard bool
	Type   string
}

// Map symbolic limit names to rlimit constants
//
var limitsByName = map[string]int{
	"core":   syscall.RLIMIT_CORE,
	"data":   syscall.RLIMIT_DATA,
	"fsize":  syscall.RLIMIT_FSIZE,
	"nofile": syscall.RLIMIT_NOFILE,
	"stack":  syscall.RLIMIT_STACK,
	"cpu":    syscall.RLIMIT_CPU,
	"as":     syscall.RLIMIT_AS,
}

// Check if a ulimit is high enough
//
func (uc UlimitChecker) Check() error {

	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(limitsByName[uc.Item], &rLimit)

	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}

	var LimitToCheck uint64
	if uc.IsHard {
		uc.Type = "hard"
		LimitToCheck = rLimit.Max
	} else {
		uc.Type = "soft"
		LimitToCheck = rLimit.Cur
	}

	if uc.Limit == syscall.RLIM_INFINITY && LimitToCheck != syscall.RLIM_INFINITY {
		return fmt.Errorf("Process %s ulimit (%d) of type \"%s\" is lower than required (unlimited)", uc.Type, LimitToCheck, uc.Item)
	} else if LimitToCheck < uc.Limit {
		return fmt.Errorf("Process %s ulimit (%d) of type \"%s\" is lower than required (%d)", uc.Type, LimitToCheck, uc.Item, uc.Limit)
	}

	return nil
}

// FromArgs will populate the UlimitChecker with the args given in the tests YAML
// config
//
// yaml inputs:
// item (required)
// limit (required)
// type ("soft"/"hard" - optional, default hard)
//
// Checker members:
// Item string
// Limit uint64
// IsHard bool
// Type string
func (uc UlimitChecker) FromArgs(args map[string]interface{}) (Checker, error) {
	if err := requiredArgs(args, "item"); err != nil {
		return nil, err
	}

	if err := requiredArgs(args, "limit"); err != nil {
		return nil, err
	}

	if args["limit"] == -1 {
		args["limit"] = syscall.RLIM_INFINITY
	}

	if err := decodeFromArgs(args, &uc); err != nil {
		return nil, err
	}

	limitType := uc.Type
	if limitType == "soft" || limitType == "" {
		uc.IsHard = false
	} else {
		uc.IsHard = true
	}

	return uc, nil
}

func init() {
	availableChecks["ulimit"] = func(args map[string]interface{}) (Checker, error) {
		return UlimitChecker{}.FromArgs(args)
	}

}

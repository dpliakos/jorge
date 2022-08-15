package jorge

import (
	"fmt"
)

type ErrorCode string
type SolutionMessage string

const (
	// System level errors
	E000 ErrorCode = "The active directory does not exist"
	E001 ErrorCode = "Directory already exist" // Already exists
	E002 ErrorCode = "Path is not a directory"
	E003 ErrorCode = "Could not create directory"
	E004 ErrorCode = "Target file is not a regular file"
	E005 ErrorCode = "Unknown error"
	E006 ErrorCode = "Could not open target file"
	E007 ErrorCode = "Could not write to file"
	E008 ErrorCode = "Error constructing path"

	// Application level errors. >= 100
	E100 = "Active directory does not belong in a jorge project"
	E101 = "Path to .jorge is not a directory"
	E102 = "Envs directory already exist"
	E103 = "Corrupted internal jorge directory (.jorge)"
	E104 = "Could not open the jorge configuration"
	E105 = "Could not read jorge config yaml file"
	E106 = "Could not write jorge config yaml file"
	E107 = "Could not read configuration file"
	E108 = "Could not Create internal directory for the configuration environment"
	E109 = "Environment already exist"
	E110 = "Could not store configuration file"
	E111 = "Environment does not exist"
)

const (
	// system level messages
	S000 SolutionMessage = "Please make sure the active directory is a valid path"
	S001 SolutionMessage = "Make sure %s is a directory"
	S002 SolutionMessage = "Make sure user %s has write access to the active directory"
	S003 SolutionMessage = "Make sure %s file exist and is readable by the user"

	// application level messages
	S100 = "Please run `jorge init` to initialize a jorge project"
	S101 = "Do not try to initialize a directory which is already a jorge project"
	S102 = "Backup configuration files manually by copying files under ./.jorge/envs and initialize the jorge project again using `jorge init`"
	S103 = "Make sure user %s has write privileges to .jorge directory"
	S104 = "Please use an environment name that is not already in use. You can use `jorge ls` to see the list of current environments"
	S105 = "You can create environment by running `jorge use -n %s`"
)

func (e ErrorCode) Str() string {
	return fmt.Sprint(e)
}

func (e ErrorCode) Err() error {
	return fmt.Errorf(fmt.Sprint(e))
}

func (e ErrorCode) Is(targetError ErrorCode) bool {
	return e == targetError
}

func (s SolutionMessage) Str(args ...interface{}) string {
	return fmt.Sprintf(fmt.Sprint(s), args...)
}

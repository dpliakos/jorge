package jorge

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const jorgeConfigDir = ".jorge"
const configFileName = "config.yml"

type JorgeConfig struct {
	CurrentEnv     string `yaml:"currentEnv"`
	ConfigFilePath string `yaml:"configFilePath"`
}

// hasJorgeDir
// Determines whether the path pass as parameter is has a .jorge dir
func hasJorgeDir(path string) bool {
	if _, err := os.Stat(filepath.Join(path, ".jorge")); err != nil {
		return false
	} else {
		return true
	}
}

// resolveJorgeDir
// The jorge command uses the current active directory as it's root. This function
// determines if the current directory has the .jorge dir in it, or if this
// directory belongs to a jorge project by searching the .jorge dir in it's
// ancestors
func resolveJorgeDir() (string, *EncapsulatedError) {
	currentAbsPath, err := filepath.Abs(filepath.Clean(filepath.Join(".")))

	if err != nil {
		encErr := EncapsulatedError{
			OriginalErr: err,
			Message:     ErrorCode.Str(E000),
			Solution:    SolutionMessage.Str(S000),
			Code:        1,
		}
		return "", &encErr
	}

	parts := strings.Split(currentAbsPath, string(os.PathSeparator))
	var currentActivePath string
	currentActivePath = currentAbsPath

	for i := 0; i < len(parts); i++ {
		if hasJorgeDir(currentActivePath) {
			log.Debug("Resolve jorge dir at ", currentActivePath)
			return currentActivePath, nil
		} else {
			currentActivePath = filepath.Join(currentActivePath, "..")
		}
	}

	encErr := EncapsulatedError{
		OriginalErr: ErrorCode.Err(E100),
		Message:     ErrorCode.Str(E100),
		Solution:    SolutionMessage.Str(S100),
		Code:        100,
	}

	return "", &encErr
}

// getJorgeDir
// Resolves the path of the .jorge directory and returns an absolute path that
// leads to it
func getJorgeDir() (string, *EncapsulatedError) {
	basePath, err := resolveJorgeDir()
	jorgeDir := filepath.Join(basePath, ".jorge")

	if err != nil {
		return "", err
	}

	if configDir, err := os.Stat(jorgeDir); err == nil {

		if !configDir.Mode().IsDir() {
			encErr := EncapsulatedError{
				OriginalErr: ErrorCode.Err(E001),
				Message:     ErrorCode.Str(E101),
				Solution:    SolutionMessage.Str(S001, jorgeDir),
				Code:        001,
			}

			return "", &encErr
		}

		log.Debug("Found .jorge dir")
		return jorgeDir, nil
	} else {
		encErr := EncapsulatedError{
			OriginalErr: err,
			Message:     ErrorCode.Str(E100),
			Solution:    SolutionMessage.Str(S000),
			Code:        100,
		}
		return "", &encErr
	}
}

// createJorgeDir
// Creates the .jorge directory in the current active directory of the process
func createJorgeDir() (string, *EncapsulatedError) {

	configDirPath := filepath.Join(jorgeConfigDir)

	if _, err := resolveJorgeDir(); err == nil {
		return "", err
	} else if ErrorCode(E100).Is(ErrorCode(err.Message)) {
		if err := os.Mkdir(configDirPath, 0700); err != nil {
			encErr := EncapsulatedError{
				OriginalErr: err,
				Message:     ErrorCode(E003).Str(),
				Solution:    SolutionMessage.Str(S002, GetUser()),
				Code:        3,
			}
			return "", &encErr
		}

		if basePath, err := resolveJorgeDir(); err == nil {
			log.Debug("Created .jorge dir")
			return filepath.Join(basePath, configDirPath), nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

// removeJorgeDir
// Deletes the .jorge directory and all of it's contents
func removeJorgeDir() *EncapsulatedError {
	if jorgeDir, err := resolveJorgeDir(); err != nil {
		return err
	} else {
		fmt.Printf("Jorge init failed. Please remove %s\n", filepath.Join(jorgeDir, ".jorge"))
		return nil
	}
}

// createJorgeEnvsDir
// Creates the .jorge/envs directory
func createJorgeEnvsDir() (string, *EncapsulatedError) {
	if jorgeDir, err := getJorgeDir(); err != nil {
		return "", err
	} else {
		if _, err := os.Stat(filepath.Join(jorgeDir, "envs")); err == nil {
			encErr := EncapsulatedError{
				OriginalErr: ErrorCode.Err(E001),
				Message:     ErrorCode.Str(E102),
				Solution:    SolutionMessage.Str(S101),
				Code:        1,
			}
			return "", &encErr
		} else if errors.Is(err, os.ErrNotExist) {
			if mkdirErr := os.Mkdir(filepath.Join(jorgeDir, "envs"), 0700); mkdirErr != nil {
				encError := EncapsulatedError{
					OriginalErr: mkdirErr,
					Message:     ErrorCode.Str(E003),
					Solution:    SolutionMessage.Str(S002, GetUser()),
					Code:        3,
				}
				return "", &encError
			} else {
				envsDirPath := filepath.Join(jorgeDir, "envs")
				log.Debug(fmt.Sprintf("Created the envs dir %s", envsDirPath))
				return envsDirPath, nil
			}
		} else {
			encError := EncapsulatedError{
				OriginalErr: err,
				Message:     ErrorCode.Str(E002),
				Solution:    SolutionMessage.Str(S100),
				Code:        2,
			}
			return "", &encError
		}
	}
}

// getEnvsDirPath
// Returns the absolute path to the envs directory
func getEnvsDirPath() (string, *EncapsulatedError) {
	jorgeDir, err := getJorgeDir()
	if err != nil {
		return "", err
	}

	envsDirPath := filepath.Join(jorgeDir, "envs")

	_, dirErr := os.Stat(envsDirPath)

	if errors.Is(dirErr, os.ErrNotExist) {
		return createJorgeEnvsDir()
	} else {
		log.Debug(fmt.Sprintf("Found the envs dir at %s", envsDirPath))
		return envsDirPath, nil
	}
}

// getInternalConfig
// Returns the current active configuration for the Jorge command.
// The internal configuration is stored a yml file found under the .jorge dir
func getInternalConfig() (JorgeConfig, *EncapsulatedError) {
	jorgeDir, encErr := getJorgeDir()
	if encErr != nil {
		return JorgeConfig{}, encErr
	}

	configFilePath := filepath.Join(jorgeDir, configFileName)
	configFileMeta, err := os.Stat(configFilePath)
	log.Debug(fmt.Sprintf("Using configuration file %s", configFilePath))

	if err != nil {
		encErr := EncapsulatedError{
			OriginalErr: err,
			Message:     ErrorCode.Str(E103),
			Solution:    SolutionMessage.Str(S102),
			Code:        103,
		}
		return JorgeConfig{}, &encErr
	}

	if !configFileMeta.Mode().IsRegular() {
		encErr := EncapsulatedError{
			OriginalErr: ErrorCode.Err(E004),
			Message:     ErrorCode.Str(E104),
			Solution:    SolutionMessage.Str(S102),
			Code:        4,
		}
		return JorgeConfig{}, &encErr
	} else {
		log.Debug("Jorge config file is regular file")
	}

	fileData, err := ioutil.ReadFile(configFilePath)

	data := make(map[string]string)
	ymlError := yaml.Unmarshal(fileData, &data)

	if ymlError != nil {
		encErr := EncapsulatedError{
			OriginalErr: ymlError,
			Message:     ErrorCode.Str(E105),
			Solution:    SolutionMessage.Str(S102),
			Code:        105,
		}
		return JorgeConfig{}, &encErr
	}

	var config = JorgeConfig{
		CurrentEnv:     data["currentEnv"],
		ConfigFilePath: data["configFilePath"],
	}

	return config, nil
}

// setInternalConfig
// Given a JorgeConfig struct, it updates the jorge configuration with the
// values found in the parameter struct
func setInternalConfig(configUpdates JorgeConfig) (JorgeConfig, *EncapsulatedError) {
	jorgeDir, err := getJorgeDir()
	if err != nil {
		return JorgeConfig{}, err
	}

	configFilePath := filepath.Join(jorgeDir, configFileName)
	currentConfig, err := getInternalConfig()

	if err != nil {
		if errors.Is(err.OriginalErr, os.ErrNotExist) {
			currentConfig = JorgeConfig{}
		} else {
			return JorgeConfig{}, err
		}
	}

	newConfig := currentConfig

	numUpdates := 0
	if currentConfig.CurrentEnv != configUpdates.CurrentEnv {
		newConfig.CurrentEnv = configUpdates.CurrentEnv
		log.Debug(fmt.Sprintf("Found updated config key 'CurrentEnv' (from '%s' to '%s')", currentConfig.CurrentEnv, configUpdates.CurrentEnv))
		numUpdates++
	}

	if len(configUpdates.ConfigFilePath) > 0 && currentConfig.ConfigFilePath != configUpdates.ConfigFilePath {
		newConfig.ConfigFilePath = configUpdates.ConfigFilePath
		log.Debug(fmt.Sprintf("Found updated config key 'ConfigFilePath' (from '%s' to '%s')", currentConfig.ConfigFilePath, configUpdates.ConfigFilePath))
		numUpdates++
	}

	if numUpdates == 0 {
		log.Debug("Called setConfig without updates")
	}

	log.Debug(newConfig)
	data, yamlErr := yaml.Marshal(&newConfig)

	if yamlErr != nil {
		return JorgeConfig{}, err
	}

	if writeError := ioutil.WriteFile(configFilePath, data, 0700); writeError != nil {
		log.Debug(fmt.Sprintf("Error while writing file"))

		encError := EncapsulatedError{
			OriginalErr: writeError,
			Message:     ErrorCode.Str(E106),
			Solution:    SolutionMessage.Str(S102),
			Code:        106,
		}
		return JorgeConfig{}, &encError
	} else {
		log.Debug(fmt.Sprintf("Wrote updated config file %s", configFilePath))
	}

	return newConfig, nil
}

// getEnvs
// It returns a list with the available environments
func getEnvs() ([]string, *EncapsulatedError) {
	envsDirPath, err := getEnvsDirPath()

	if err != nil {
		return []string{}, err
	}

	log.Debug(fmt.Sprintf("Using envs dir %s", envsDirPath))
	files, readEnvsErr := ioutil.ReadDir(envsDirPath)

	log.Debug(fmt.Sprintf("Found %d files under the envs dir.", len(files)))

	if readEnvsErr != nil {
		encErr := EncapsulatedError{
			OriginalErr: readEnvsErr,
			Message:     ErrorCode.Str(E103),
			Solution:    SolutionMessage.Str(S102),
			Code:        103,
		}
		return []string{}, &encErr
	}

	envs := make([]string, len(files))

	for i := range files {
		envs[i] = files[i].Name()
	}

	return envs, nil
}

// setConfigAsMain
// Given an existing environment, it replaces the user configuration file, with
// the one that is stored for the jorge environment
func setConfigAsMain(target string, envName string) (int64, *EncapsulatedError) {
	envsDir, err := getEnvsDirPath()

	if err != nil {
		return -1, err
	}

	_, targeFileName := filepath.Split(target)
	log.Debug(fmt.Sprintf("Target file path %v. Found name %v", target, targeFileName))
	storedFile, storedFileErr := os.Stat(filepath.Join(envsDir, envName, targeFileName))

	if storedFileErr != nil {
		if errors.Is(storedFileErr, os.ErrNotExist) {
			encErr := EncapsulatedError{
				OriginalErr: storedFileErr,
				Message:     ErrorCode.Str(E111),
				Solution:    SolutionMessage.Str(S105, envName),
				Code:        111,
			}
			return -1, &encErr
		} else {
			encErr := EncapsulatedError{
				OriginalErr: storedFileErr,
				Message:     ErrorCode.Str(E005),
				Solution:    SolutionMessage.Str(S102),
				Code:        5,
			}
			return -1, &encErr
		}
	}

	if !storedFile.Mode().IsRegular() {
		return -1, err
	}
	log.Debug(fmt.Sprintf("File %v is regular.", target))

	// ----
	sourceFilePath := filepath.Join(envsDir, envName, targeFileName)
	sourceFile, sourceFileErr := os.Open(sourceFilePath)
	defer sourceFile.Close()

	if sourceFileErr != nil {
		encErr := EncapsulatedError{
			OriginalErr: sourceFileErr,
			Message:     ErrorCode.Str(E006),
			Solution:    SolutionMessage.Str(S003, sourceFilePath),
			Code:        6,
		}

		return -1, &encErr
	}

	log.Debug(fmt.Sprintf("Source file %v found", sourceFilePath))

	_, fileName := filepath.Split(target)

	destination, destinationErr := os.Create(filepath.Join(target))
	log.Debug(fmt.Sprintf("Target file path %v", target))
	if destinationErr != nil {
		encErr := EncapsulatedError{
			OriginalErr: destinationErr,
			Message:     ErrorCode.Str(E007),
			Solution:    SolutionMessage.Str(S002, GetUser()),
			Code:        7,
		}
		return -1, &encErr
	}
	defer destination.Close()

	nBytes, copyErr := io.Copy(destination, sourceFile)
	log.Debug(fmt.Sprintf("Wrote %d bytes to %v", nBytes, fileName))

	if copyErr != nil {
		encError := EncapsulatedError{
			OriginalErr: copyErr,
			Message:     ErrorCode.Str(E007),
			Solution:    SolutionMessage.Str(S002, GetUser()),
			Code:        7,
		}
		return -1, &encError
	}
	return nBytes, nil
}

// requestConfigFileFromUser
// It shows a cli prompt that accepts a string. The string is expected to be the
// filepath to the user's configuration file
func requestConfigFileFromUser() (string, *EncapsulatedError) {
	fmt.Print("Config file path: ")
	var filePath string

	fmt.Scanln(&filePath)

	configFileRelativePath := filepath.Clean(filePath)

	if _, err := os.Stat(configFileRelativePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			encErr := EncapsulatedError{
				OriginalErr: err,
				Message:     ErrorCode.Str(E006),
				Solution:    SolutionMessage.Str(S003, configFileRelativePath),
				Code:        6,
			}
			return "", &encErr
		} else {
			encErr := EncapsulatedError{
				OriginalErr: err,
				Message:     ErrorCode.Str(E005),
				Solution:    SolutionMessage.Str(S003, configFileRelativePath),
				Code:        5,
			}
			return "", &encErr
		}
	} else {
		return configFileRelativePath, nil
	}
}

func initializeJorgeProject(configFilePathFlag string) *EncapsulatedError {
	jorgeDir, err := createJorgeDir()
	if err != nil {
		return err
	}

	var configFileName string

	if len(configFilePathFlag) > 0 {
		configFileName = configFilePathFlag
	} else {
		configFileName, err = requestConfigFileFromUser()
	}

	if err != nil {
		return err
	}

	absConfigFileName, absConfigFileNameErr := filepath.Abs(configFileName)

	if absConfigFileNameErr != nil {
		encErr := EncapsulatedError{
			OriginalErr: absConfigFileNameErr,
			Message:     ErrorCode.Str(E008),
			Solution:    SolutionMessage.Str(S003, absConfigFileName),
			Code:        8,
		}

		return &encErr
	}

	if _, err := os.Stat(configFileName); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			encErr := EncapsulatedError{
				OriginalErr: err,
				Message:     ErrorCode.Str(E107),
				Solution:    SolutionMessage.Str(S003, configFileName),
				Code:        107,
			}
			return &encErr
		}
	}

	var relativePathToConfig string
	var relativePathToConfigErr error
	if relativePathToConfig, relativePathToConfigErr = filepath.Rel(filepath.Join(jorgeDir, ".."), absConfigFileName); relativePathToConfigErr != nil {
		encErr := EncapsulatedError{
			OriginalErr: relativePathToConfigErr,
			Message:     ErrorCode.Str(E107),
			Solution:    SolutionMessage.Str(S003, absConfigFileName),
			Code:        107,
		}
		return &encErr
	}

	freshJorgeConfig := JorgeConfig{
		CurrentEnv:     "default",
		ConfigFilePath: relativePathToConfig,
	}

	if _, err := setInternalConfig(freshJorgeConfig); err != nil {
		return err
	}

	_, storeFileErr := StoreConfigFile(freshJorgeConfig.ConfigFilePath, "default")

	if storeFileErr != nil {
		return storeFileErr
	}

	jorgeRecordExist, _ := ExistsInFile(".gitignore", ".jorge")

	if !jorgeRecordExist {
		AppendToFile(".gitignore", ".jorge")
	}

	return nil
}

func deleteJorgeEnv(env string) *EncapsulatedError {
	jorgeDir, err := getJorgeDir()
	if err != nil {
		fmt.Println("Error")
		return nil
	}

	target := filepath.Join(jorgeDir, "envs", env)

	if _, err := os.Stat(target); err != nil {
		encError := EncapsulatedError{
			OriginalErr: err,
			Message:     ErrorCode.Str(E002),
			Solution:    SolutionMessage.Str(S001, target),
			Code:        2,
		}

		return &encError
	}

	removeErr := os.RemoveAll(target)
	if removeErr != nil {
		fmt.Println("could not remove ", target)
		encError := EncapsulatedError{
			OriginalErr: removeErr,
			Message:     ErrorCode.Str(E009),
			Solution:    SolutionMessage.Str(S001, target),
			Code:        9,
		}

		return &encError
	}

	return nil
}

// StoreConfigFile
// It stores the current active user config file under an jorge environment name
func StoreConfigFile(path string, envName string) (int64, *EncapsulatedError) {
	envsDir, err := getEnvsDirPath()

	if err != nil {
		return -1, err
	}

	activeConfigFileMeta, activeConfigFileMetaErr := os.Stat(path)

	if activeConfigFileMetaErr != nil {
		encErr := EncapsulatedError{
			OriginalErr: activeConfigFileMetaErr,
			Message:     ErrorCode.Str(E107),
			Solution:    SolutionMessage.Str(S003, path),
			Code:        107,
		}
		return -1, &encErr
	}

	if !activeConfigFileMeta.Mode().IsRegular() {
		encErr := EncapsulatedError{
			OriginalErr: ErrorCode.Err(E004),
			Message:     ErrorCode.Str(E004),
			Solution:    SolutionMessage.Str(S003, path),
			Code:        4,
		}
		return -2, &encErr
	} else {
		log.Debug("Active configuration file " + path + " is a regular file")
	}

	targetEnvDirName := filepath.Join(envsDir, envName)
	if _, err := os.Stat(targetEnvDirName); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			mkdirErr := os.Mkdir(targetEnvDirName, 0700)
			if mkdirErr != nil {
				encErr := EncapsulatedError{
					OriginalErr: mkdirErr,
					Message:     ErrorCode.Str(E108),
					Solution:    SolutionMessage.Str(S103, GetUser()),
					Code:        108,
				}
				return -3, &encErr
			}
			log.Debug(fmt.Sprintf("Created env dir %s", targetEnvDirName))
		} else {
			encError := EncapsulatedError{
				OriginalErr: err,
				Message:     ErrorCode.Str(E005),
				Solution:    SolutionMessage.Str(S103, GetUser()),
				Code:        5,
			}
			return -4, &encError
		}
	}

	log.Debug(fmt.Sprintf("Source path is %v", path))
	sourceFile, sourceFileErr := os.Open(path)
	log.Debug("Source file opened")
	defer sourceFile.Close()

	if sourceFileErr != nil {
		encErr := EncapsulatedError{
			OriginalErr: sourceFileErr,
			Message:     ErrorCode.Str(E107),
			Solution:    SolutionMessage.Str(S103, GetUser()),
			Code:        107,
		}

		return -1, &encErr
	}

	_, fileName := filepath.Split(path)

	destinationPath := filepath.Join(targetEnvDirName, fileName)
	destination, destinationErr := os.Create(destinationPath)

	if destinationErr != nil {
		encErr := EncapsulatedError{
			OriginalErr: destinationErr,
			Message:     ErrorCode.Str(E110),
			Solution:    SolutionMessage.Str(S103, GetUser()),
			Code:        110,
		}
		return -1, &encErr
	}

	log.Debug(fmt.Sprintf("Stored file %s", destinationPath))
	log.Debug("Destination file created")
	defer destination.Close()

	nBytes, copyErr := io.Copy(destination, sourceFile)
	log.Debug(fmt.Sprintf("Wrote %d bytes from %v to %v", nBytes, path, destinationPath))

	if copyErr != nil {
		encErr := EncapsulatedError{
			OriginalErr: copyErr,
			Message:     ErrorCode.Str(E110),
			Solution:    SolutionMessage.Str(S103, GetUser()),
			Code:        110,
		}

		return -1, &encErr
	}

	return nBytes, nil
}

// UseConfigFile
// It replaces the current active user configuration file with the one that is
// stored under the jorge environment
func UseConfigFile(envName string, createEnv bool) (int64, *EncapsulatedError) {

	config, err := getInternalConfig()

	if err != nil {
		return -1, err
	}

	target := config.ConfigFilePath

	if createEnv {
		if existingEnvs, err := getEnvs(); err == nil {
			if Contains(existingEnvs, envName) {
				encErr := EncapsulatedError{
					OriginalErr: ErrorCode.Err(E109),
					Message:     ErrorCode.Str(E109),
					Solution:    SolutionMessage.Str(S104),
					Code:        109,
				}
				return -1, &encErr
			} else {
				if _, err := StoreConfigFile(target, envName); err != nil {
					return -1, err
				} else {
					log.Debug(fmt.Sprintf("Created new file for env %s", envName))
				}
			}
		} else {
			return -1, err
		}
	}

	if _, err := setConfigAsMain(target, envName); err != nil {
		return -1, err
	} else {
		log.Debug(fmt.Sprintf("Used %s as main config file", envName))

		newConfig := JorgeConfig{
			CurrentEnv: envName,
		}

		if _, err := setInternalConfig(newConfig); err != nil {
			return -1, err
		} else {
			return 1, nil
		}
	}
}

func SelectEnvironment(envName string) *EncapsulatedError {
	if _, err := setInternalConfig(JorgeConfig{
		CurrentEnv: envName,
	}); err != nil {
		return err
	} else {
		fmt.Println(fmt.Sprintf("Using environment %s", envName))
		return nil
	}
}

// ListEnvironments
// Shows a list with all the available environment for the user
func ListEnvironments() *EncapsulatedError {
	envs, err := getEnvs()

	if err != nil {
		return err
	}

	config, err := getInternalConfig()
	if err != nil {
		return err
	}

	currentEnvFound := false
	for _, fileName := range envs {
		if fileName == config.CurrentEnv {
			currentEnvFound = true
			fmt.Printf("* %s\n", fileName)
		} else {
			fmt.Println(fileName)
		}
	}

	if !currentEnvFound {
		fmt.Printf("* %s (uncommitted)\n", config.CurrentEnv)
	}

	return nil
}

// Init
// Initializes a jorge project by creating the .jorge directory, setting the
// first environment and requesting the config file path for the user
func Init(configFilePathFlag string) *EncapsulatedError {
	if err := initializeJorgeProject(configFilePathFlag); err != nil {
		removeJorgeDir()
		return err
	} else {
		return nil
	}
}

// CommitCurrentEnv
// It stores the current active user configuration file under the current selected
// jorge environment (found at the ./jorge/config.yml file)
func CommitCurrentEnv() *EncapsulatedError {
	config, err := getInternalConfig()

	if err != nil {
		return err
	}
	jorgeDir, err := resolveJorgeDir()

	if err != nil {
		return err
	}

	activeUserConfig := filepath.Join(jorgeDir, config.ConfigFilePath)
	_, storeConfigErr := StoreConfigFile(activeUserConfig, config.CurrentEnv)

	if storeConfigErr != nil {
		return storeConfigErr
	} else {
		return nil
	}
}

// RestoreEnv
// It replaces the current active configuration file, with the one that is
// stored under the current environment (found under the ./.jorge/config.yml)
func RestoreEnv() *EncapsulatedError {
	config, err := getInternalConfig()

	if err != nil {
		return err
	}

	jorgeDir, err := resolveJorgeDir()
	if err != nil {
		return err
	}

	activeUserConfig := filepath.Join(jorgeDir, config.ConfigFilePath)
	_, restoreError := setConfigAsMain(activeUserConfig, config.CurrentEnv)

	if restoreError != nil {
		return restoreError
	} else {
		return nil
	}
}

func RemoveEnv(envName string) *EncapsulatedError {
	envs, err := getEnvs()

	if err != nil {
		return err
	}

	config, err := getInternalConfig()
	if err != nil {
		return err
	}

	if envName == config.CurrentEnv {
		encErr := EncapsulatedError{
			OriginalErr: ErrorCode.Err(E112),
			Message:     ErrorCode.Str(E112),
			Solution:    SolutionMessage.Str(S106),
			Code:        112,
		}

		return &encErr
	}

	targetEnvFound := false
	for _, fileName := range envs {
		if fileName == envName {
			targetEnvFound = true
		}
	}

	if !targetEnvFound {
		encErr := EncapsulatedError{
			OriginalErr: ErrorCode.Err(E111),
			Message:     ErrorCode.Str(E111),
			Code:        111,
		}

		return &encErr
	}

	if err = deleteJorgeEnv(envName); err != nil {
		return err
	}

	return nil
}

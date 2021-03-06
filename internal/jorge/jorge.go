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
func resolveJorgeDir() (string, error) {
	currentAbsPath, err := filepath.Abs(filepath.Clean(filepath.Join(".")))

	if err != nil {
		return "", err
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

	return "", fmt.Errorf("Current directory does not belong to a jorge project")
}

// getJorgeDir
// Resolves the path of the .jorge directory and returns an absolute path that
// leads to it
func getJorgeDir() (string, error) {
	basePath, err := resolveJorgeDir()
	jorgeDir := filepath.Join(basePath, ".jorge")

	if err != nil {
		return "", err
	}

	if configDir, err := os.Stat(jorgeDir); err == nil {

		if !configDir.Mode().IsDir() {
			return "", fmt.Errorf(fmt.Sprintf("path %s exist and it's not a dir", jorgeDir))
		}

		log.Debug("Found .jorge dir")
		return jorgeDir, nil
	} else {
		return "", err
	}
}

// createJorgeDir
// Creates the .jorge directory in the current active directory of the process
func createJorgeDir() (string, error) {

	configDirPath := filepath.Join(jorgeConfigDir)

	if _, err := resolveJorgeDir(); err == nil {
		return "", fmt.Errorf("This directory belongs to a jorge project")
	} else if err.Error() == "Current directory does not belong to a jorge project" {
		if err := os.Mkdir(configDirPath, 0700); err != nil {
			return "", err
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
func removeJorgeDir() error {
	if jorgeDir, err := resolveJorgeDir(); err != nil {
		return err
	} else {
		fmt.Printf("Jorge init failed. Please remove %s\n", filepath.Join(jorgeDir, ".jorge"))
		return nil
	}
}

// createJorgeEnvsDir
// Creates the .jorge/envs directory
func createJorgeEnvsDir() (string, error) {
	if jorgeDir, err := getJorgeDir(); err != nil {
		return "", err
	} else {
		if _, err := os.Stat(filepath.Join(jorgeDir, "envs")); err == nil {
			return "", fmt.Errorf("Envs directory already created")
		} else if errors.Is(err, os.ErrNotExist) {
			if mkdirErr := os.Mkdir(filepath.Join(jorgeDir, "envs"), 0700); mkdirErr != nil {
				return "", err
			} else {
				envsDirPath := filepath.Join(jorgeDir, "envs")
				log.Debug(fmt.Sprintf("Created the envs dir %s", envsDirPath))
				return envsDirPath, nil
			}
		} else {
			return "", err
		}
	}
}

// getEnvsDirPath
// Returns the absolute path to the envs directory
func getEnvsDirPath() (string, error) {
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
func getInternalConfig() (JorgeConfig, error) {
	jorgeDir, err := getJorgeDir()
	if err != nil {
		return JorgeConfig{}, err
	}

	configFilePath := filepath.Join(jorgeDir, configFileName)
	configFileMeta, err := os.Stat(configFilePath)
	log.Debug(fmt.Sprintf("Using configuration file %s", configFilePath))

	if err != nil {
		return JorgeConfig{}, err
	}

	if !configFileMeta.Mode().IsRegular() {
		return JorgeConfig{}, fmt.Errorf("jorge config is not a regular file")
	} else {
		log.Debug("Jorge config file is regular file")
	}

	fileData, err := ioutil.ReadFile(configFilePath)

	data := make(map[string]string)
	err2 := yaml.Unmarshal(fileData, &data)

	if err2 != nil {
		log.Fatal(err2)
		return JorgeConfig{}, err2
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
func setInternalConfig(configUpdates JorgeConfig) (JorgeConfig, error) {
	jorgeDir, err := getJorgeDir()
	if err != nil {
		return JorgeConfig{}, err
	}

	configFilePath := filepath.Join(jorgeDir, configFileName)
	currentConfig, err := getInternalConfig()

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
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
		return JorgeConfig{}, writeError
	} else {
		log.Debug(fmt.Sprintf("Wrote updated config file %s", configFilePath))
	}

	return newConfig, nil
}

// getEnvs
// It returns a list with the available environments
func getEnvs() ([]string, error) {
	envsDirPath, err := getEnvsDirPath()

	if err != nil {
		return []string{}, err
	}

	log.Debug(fmt.Sprintf("Using envs dir %s", envsDirPath))
	files, err := ioutil.ReadDir(envsDirPath)

	log.Debug(fmt.Sprintf("Found %d files under the envs dir.", len(files)))

	if err != nil {
		return []string{}, err
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
func setConfigAsMain(target string, envName string) (int64, error) {
	envsDir, err := getEnvsDirPath()

	if err != nil {
		return -1, err
	}

	_, targeFileName := filepath.Split(target)
	log.Debug(fmt.Sprintf("Target file path %v. Found name %v", target, targeFileName))
	storedFile, err := os.Stat(filepath.Join(envsDir, envName, targeFileName))

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return -6, err
		} else {
			return -7, err
		}
	}

	if !storedFile.Mode().IsRegular() {
		return -8, err
	}
	log.Debug(fmt.Sprintf("File %v is regular.", target))

	// ----
	sourceFilePath := filepath.Join(envsDir, envName, targeFileName)
	sourceFile, err := os.Open(sourceFilePath)
	defer sourceFile.Close()
	log.Debug(fmt.Sprintf("Source file %v found", sourceFilePath))

	_, fileName := filepath.Split(target)

	destination, err := os.Create(filepath.Join(target))
	log.Debug(fmt.Sprintf("Target file path %v", target))
	if err != nil {
		return -9, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, sourceFile)
	log.Debug(fmt.Sprintf("Wrote %d bytes to %v", nBytes, fileName))
	return nBytes, err
}

// requestConfigFileFromUser
// It shows a cli prompt that accepts a string. The string is expected to be the
// filepath to the user's configuration file
func requestConfigFileFromUser() (string, error) {
	fmt.Print("Config file path: ")
	var filePath string

	fmt.Scanln(&filePath)

	configFileRelativePath := filepath.Clean(filePath)

	if _, err := os.Stat(configFileRelativePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", err
		} else {
			return "", err
		}
	} else {
		return configFileRelativePath, nil
	}
}

// StoreConfigFile
// It stores the current active user config file under an jorge environment name
func StoreConfigFile(path string, envName string) (int64, error) {
	envsDir, err := getEnvsDirPath()

	if err != nil {
		return -1, err
	}

	activeConfigFileMeta, err := os.Stat(path)

	if err != nil {
		return -1, err
	}

	if !activeConfigFileMeta.Mode().IsRegular() {
		return -2, fmt.Errorf("%s is not a regular file", path)
	} else {
		log.Debug("Active configuration file " + path + " is a regular file")
	}

	targetEnvDirName := filepath.Join(envsDir, envName)
	if _, err := os.Stat(targetEnvDirName); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			mkdirErr := os.Mkdir(targetEnvDirName, 0700)
			if mkdirErr != nil {
				return -3, mkdirErr
			}
			log.Debug(fmt.Sprintf("Created env dir %s", targetEnvDirName))
		} else {
			return -4, err
		}
	}

	log.Debug(fmt.Sprintf("Source path is %v", path))
	sourceFile, err := os.Open(path)
	log.Debug("Source file opened")
	defer sourceFile.Close()

	_, fileName := filepath.Split(path)

	destinationPath := filepath.Join(targetEnvDirName, fileName)
	destination, err := os.Create(destinationPath)
	log.Debug(fmt.Sprintf("Stored file %s", destinationPath))

	if err != nil {
		return -5, err
	}
	log.Debug("Destination file created")
	defer destination.Close()

	nBytes, err := io.Copy(destination, sourceFile)
	log.Debug(fmt.Sprintf("Wrote %d bytes from %v to %v", nBytes, path, destinationPath))
	return nBytes, err
}

// UseConfigFile
// It replaces the current active user configuration file with the one that is
// stored under the jorge environment
func UseConfigFile(envName string, createEnv bool) (int64, error) {

	config, err := getInternalConfig()

	if err != nil {
		return -1, err
	}

	target := config.ConfigFilePath

	if createEnv {
		if existingEnvs, err := getEnvs(); err == nil {
			if Contains(existingEnvs, envName) {
				return -1, fmt.Errorf(fmt.Sprintf("Environment %s already exist", envName))
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
			return -2, err
		} else {
			return 1, nil
		}
	}
}

func SelectEnvironment(envName string) error {
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
func ListEnvironments() error {
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

func cleanupInit() {
	if a := recover(); a != nil {
		fmt.Println(a)
		removeJorgeDir()
		panic(a)
	}
}

// Init
// Initializes a jorge project by creating the .jorge directory, setting the
// first environment and requesting the config file path for the user
func Init(configFilePathFlag string) error {
	jorgeDir, err := createJorgeDir()
	if err != nil {
		return err
	}

	defer cleanupInit()

	var configFileName string

	if len(configFilePathFlag) > 0 {
		configFileName = configFilePathFlag
	} else {
		configFileName, err = requestConfigFileFromUser()
	}

	if err != nil {
		panic(err)
	}

	absConfigFileName, err := filepath.Abs(configFileName)

	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(configFileName); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			panic(fmt.Errorf("File %s does not exist", configFileName))
		}
	}

	var relativePathToConfig string
	if relativePathToConfig, err = filepath.Rel(filepath.Join(jorgeDir, ".."), absConfigFileName); err != nil {
		panic(err)
	}

	freshJorgeConfig := JorgeConfig{
		CurrentEnv:     "default",
		ConfigFilePath: relativePathToConfig,
	}

	if _, err := setInternalConfig(freshJorgeConfig); err != nil {
		panic(err)
	}

	_, storeFileErr := StoreConfigFile(freshJorgeConfig.ConfigFilePath, "default")

	if storeFileErr != nil {
		panic(err)
	}

	jorgeRecordExist, err := ExistsInFile(".gitignore", ".jorge")

	if !jorgeRecordExist {
		AppendToFile(".gitignore", ".jorge")
	}

	return nil
}

// CommitCurrentEnv
// It stores the current active user configuration file under the current selected
// jorge environment (found at the ./jorge/config.yml file)
func CommitCurrentEnv() error {
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
func RestoreEnv() error {
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

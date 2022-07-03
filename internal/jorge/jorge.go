package jorge

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const jorgeConfigDir = ".jorge"
const configFileName = "config.yml"

type JorgeConfig struct {
	CurrentEnv     string `yaml:"currentEnv"`
	ConfigFilePath string `yaml:"configFilePath"`
}

func getJorgeDir() (string, error) {
	// TODO: check recursively for the parent dir
	configFilePath := filepath.Join(jorgeConfigDir)

	if configDir, err := os.Stat(configFilePath); err == nil {

		if !configDir.Mode().IsDir() {
			return "", fmt.Errorf(fmt.Sprintf("path %s exist and it's not a dir", configFilePath))
		}

		log.Debug("Found .jorge dir")
		return configFilePath, nil
	} else if errors.Is(err, os.ErrNotExist) {
		mkdirErr := os.Mkdir(jorgeConfigDir, 0700) // TODO: change permissions

		if mkdirErr != nil {
			return "", mkdirErr
		} else {
			log.Debug("Created .jorge dir")
			return configFilePath, nil
		}
	} else {
		return "", err
	}
}

func getEnvsDirPath() (string, error) {
	jorgeDir, err := getJorgeDir()
	if err != nil {
		return "", err
	}

	envsDirPath := filepath.Join(jorgeDir, "envs")

	_, dirErr := os.Stat(envsDirPath)

	if errors.Is(dirErr, os.ErrNotExist) {
		err := os.Mkdir(envsDirPath, 0700)

		if err != nil {
			return "", err
		} else {
			log.Debug(fmt.Sprintf("Created the envs dir %s", envsDirPath))
		}
	} else {
		log.Debug(fmt.Sprintf("Found the envs dir at %s", envsDirPath))
	}

	return envsDirPath, nil
}

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

func Init() error {
	jorgePath, err := getJorgeDir()

	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(jorgePath, configFileName)); err == nil {
		fmt.Println(filepath.Join(jorgePath, configFileName))
		return fmt.Errorf("Jorge is already initialized")
	}

	configFileName, err := requestConfigFileFromUser()

	if err != nil {
		return err
	}

	freshJorgeConfig := JorgeConfig{
		CurrentEnv:     "default",
		ConfigFilePath: configFileName,
	}

	if _, err := setInternalConfig(freshJorgeConfig); err != nil {
		return err
	}

	_, storeFileErr := StoreConfigFile(freshJorgeConfig.ConfigFilePath, "default")

	if storeFileErr != nil {
		return err
	}

	jorgeRecordExist, err := ExistsInFile(".gitignore", ".jorge")

	if !jorgeRecordExist {
		AppendToFile(".gitignore", ".jorge")
	}

	return nil
}

func CommitCurrentEnv() error {
	config, err := getInternalConfig()

	if err != nil {
		return err
	}

	_, storeConfigErr := StoreConfigFile(config.ConfigFilePath, config.CurrentEnv)

	if storeConfigErr != nil {
		return storeConfigErr
	} else {
		return nil
	}
}

func RestoreEnv() error {
	config, err := getInternalConfig()

	if err != nil {
		return err
	}

	_, restoreError := setConfigAsMain(config.ConfigFilePath, config.CurrentEnv)

	if restoreError != nil {
		return restoreError
	} else {
		return nil
	}
}

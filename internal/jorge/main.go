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
	currentEnv string
}

func getJorgeDir() (string, error) {
	// TODO: check recursively for the parent dir
	configFilePath := filepath.Join(jorgeConfigDir)

	if _, err := os.Stat(configFilePath); err == nil {
		log.Debug("Found .jorge dir")
		return configFilePath, nil
	} else if errors.Is(err, os.ErrNotExist) {
		mkdirErr := os.Mkdir(jorgeConfigDir, 0777) // TODO: change permissions

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
		err := os.Mkdir(envsDirPath, 0777)

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

func getConfig() (JorgeConfig, error) {
	jorgeDir, err := getJorgeDir()
	if err != nil {
		return JorgeConfig{}, err
	}

	configFilePath := filepath.Join(jorgeDir, configFileName)
	configFileMeta, err := os.Stat(configFilePath)
	log.Debug(fmt.Sprintf("Using configuration file %s", configFilePath))

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// TODO create
		} else {
			return JorgeConfig{}, err
		}
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
		currentEnv: data["currentEnv"],
	}

	return config, nil
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
			mkdirErr := os.Mkdir(targetEnvDirName, 0777)
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
	fmt.Println(nBytes)
	return nBytes, err
}

func UseConfigFile(target string, envName string) (int64, error) {
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

	destination, err := os.Create(filepath.Join(fileName))
	log.Debug(fmt.Sprintf("Target file path %v", fileName))
	if err != nil {
		return -9, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, sourceFile)
	log.Debug(fmt.Sprintf("Wrote %d bytes to %v", nBytes, fileName))
	return nBytes, err
}

func ListEnvironments() error {
	envs, err := getEnvs()

	if err != nil {
		return err
	}

	config, err := getConfig()
	if err != nil {
		return err
	}

	currentEnvFound := false
	for _, fileName := range envs {
		if fileName == config.currentEnv {
			currentEnvFound = true
			fmt.Printf("* %s\n", fileName)
		} else {
			fmt.Println(fileName)
		}
	}

	if !currentEnvFound {
		fmt.Printf("* %s (uncommitted)\n", config.currentEnv)
	}

	return nil
}

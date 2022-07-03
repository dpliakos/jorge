package jorge

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetJorgeDir(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	path, err := getJorgeDir()
	defer os.Remove(".jorge")

	if err != nil {
		t.Fail()
	}

	if path != ".jorge" {
		t.Fail()
	}
}

func TestGetJorgeDirCreatesJorgeDir(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	_, err := getJorgeDir()
	defer os.Remove(".jorge")

	if err != nil {
		t.Fail()
	}

	if _, err := os.Stat(filepath.Join(testingRoot, ".jorge")); err != nil {
		t.Fail()
	}
}

func TestGetJorgeDirWithError(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	os.Create(filepath.Join(testingRoot, ".jorge"))
	defer os.Remove(filepath.Join(testingRoot, ".jorge"))

	path, err := getJorgeDir()

	if err == nil || path != "" {
		t.Fail()
	}
}

func TestGetEnvsPathDir(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	path, err := getEnvsDirPath()
	defer os.Remove(filepath.Join(testingRoot, ".jorge"))
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs"))

	t.Log(path)
	t.Log(err)

	if err != nil || path != filepath.Join(".jorge", "envs") {
		t.Fail()
	}
}

func TestGetInternalConfig(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	os.Mkdir(filepath.Join(testingRoot, ".jorge"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge"))

	if err := os.WriteFile(filepath.Join(testingRoot, ".jorge", "config.yml"), []byte("currentEnv: default\nconfigFilePath: .env"), 0600); err != nil {
		t.Fail()
	}
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "config.yml"))

	if config, err := getInternalConfig(); err != nil {
		t.Fail()
	} else if config.ConfigFilePath != ".env" {
		t.Fail()
	} else if config.CurrentEnv != "default" {
		t.Fail()
	}
}

func TestSetInternalConfig(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	os.Mkdir(filepath.Join(testingRoot, ".jorge"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge"))

	newConfig := JorgeConfig{
		CurrentEnv:     "mockEnv",
		ConfigFilePath: "mockFilePath",
	}

	setInternalConfig(newConfig)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "config.yml"))

	if data, err := os.ReadFile(filepath.Join(testingRoot, ".jorge", "config.yml")); err != nil {
		t.Log(err)
		t.Fail()
	} else if strings.Index(string(data), "currentEnv: mockEnv") != 0 {
		t.Fail()
	} else if strings.Index(string(data), "configFilePath: mockFilePath") < 0 {
		t.Fail()
	}
}

func TestGetEnvs(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	os.Mkdir(filepath.Join(testingRoot, ".jorge"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge"))

	os.Mkdir(filepath.Join(testingRoot, ".jorge", "envs"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs"))

	os.Mkdir(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv"))

	os.WriteFile(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv", "mockConfig"), []byte("mock config contents"), 0600)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv", "mockConfig"))

	if envs, err := getEnvs(); err != nil {
		t.Fail()
	} else if len(envs) != 1 {
		t.Fail()
	} else if envs[0] != "mockEnv" {
		t.Fail()
	}
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs"))
}

func TestGetConfigAsMain(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	os.Mkdir(filepath.Join(testingRoot, ".jorge"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge"))

	os.Mkdir(filepath.Join(testingRoot, ".jorge", "envs"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs"))

	os.Mkdir(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv"))

	os.WriteFile(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv", "mainTestConfig"), []byte("mock config contents"), 0600)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv", "mainTestConfig"))

	if _, err := setConfigAsMain("mainTestConfig", "mockEnv"); err != nil {
		t.Log(err)
		t.Fail()
	}
	defer os.Remove(filepath.Join(testingRoot, "mainTestConfig"))

	if _, err := os.Stat(filepath.Join(testingRoot, "mainTestConfig")); err != nil {
		t.Log("Config file mainTestConfig was not found")
		t.FailNow()
	}

	if data, err := os.ReadFile(filepath.Join(testingRoot, "mainTestConfig")); err != nil {
		t.Log("Error reading mainTestConfig")
		t.Fail()
	} else if string(data) != "mock config contents" {
		t.Log("mainTestConfigContents does not match")
		t.Fail()
	}
}

func TestStoreConfigFile(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	os.Mkdir(filepath.Join(testingRoot, ".jorge"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge"))

	os.Mkdir(filepath.Join(testingRoot, ".jorge", "envs"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs"))

	os.WriteFile(filepath.Join(testingRoot, "mainTestConfig"), []byte("mock config contents"), 0600)
	defer os.Remove(filepath.Join(testingRoot, "mainTestConfig"))

	if _, err := StoreConfigFile(filepath.Join(testingRoot, "mainTestConfig"), "newMockEnv"); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer os.RemoveAll(filepath.Join(testingRoot, ".jorge", "envs"))

	if _, err := os.Stat(filepath.Join(testingRoot, ".jorge", "envs", "newMockEnv")); err != nil {
		t.Log(err)
		t.Fail()
	}

	if data, err := os.ReadFile(filepath.Join(testingRoot, ".jorge", "envs", "newMockEnv", "mainTestConfig")); err != nil {
		t.Log(err)
		t.FailNow()
	} else if string(data) != "mock config contents" {
		t.Logf("Expected %s, but found %s", "mock config contents", string(data))
		t.FailNow()
	}
}

func TestUseConfigFileWhenEnvExists(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	os.Mkdir(filepath.Join(testingRoot, ".jorge"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge"))

	if err := os.WriteFile(filepath.Join(testingRoot, ".jorge", "config.yml"), []byte("currentEnv: default\nconfigFilePath: mainTestConfig"), 0600); err != nil {
		t.Fail()
	}
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "config.yml"))

	os.Mkdir(filepath.Join(testingRoot, ".jorge", "envs"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs"))

	os.Mkdir(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv"))

	os.WriteFile(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv", "mainTestConfig"), []byte("mock config contents"), 0600)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv", "mainTestConfig"))

	if _, err := UseConfigFile("mockEnv", false); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if _, err := os.Stat(filepath.Join(testingRoot, "mainTestConfig")); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if data, err := os.ReadFile(filepath.Join(testingRoot, "mainTestConfig")); err != nil {
		t.Log(err)
		t.FailNow()
	} else if string(data) != "mock config contents" {
		t.Logf("Expected %s, but found %s", "mock config contents", string(data))
		t.FailNow()
	}

	if data, err := os.ReadFile(filepath.Join(testingRoot, ".jorge", "config.yml")); err != nil {
		t.Log(err)
		t.FailNow()
	} else if strings.Index(string(data), "currentEnv: mockEnv") < 0 {
		t.Log("Current env is not mockEnv")
		t.Fail()
	}
}

func TestUseConfigFileWhenEnvDoesNotExist(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	os.Mkdir(filepath.Join(testingRoot, ".jorge"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge"))

	if err := os.WriteFile(filepath.Join(testingRoot, ".jorge", "config.yml"), []byte("currentEnv: default\nconfigFilePath: mainTestConfig"), 0600); err != nil {
		t.Fail()
	}
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "config.yml"))

	if _, err := UseConfigFile("mockEnv", true); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer os.RemoveAll(filepath.Join(testingRoot, ".jorge", "envs"))

	if _, err := os.Stat(filepath.Join(testingRoot, "mainTestConfig")); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if data, err := os.ReadFile(filepath.Join(testingRoot, "mainTestConfig")); err != nil {
		t.Log(err)
		t.FailNow()
	} else if string(data) != "mock config contents" {
		t.Logf("Expected %s, but found %s", "mock config contents", string(data))
		t.FailNow()
	}

	if data, err := os.ReadFile(filepath.Join(testingRoot, ".jorge", "config.yml")); err != nil {
		t.Log(err)
		t.FailNow()
	} else if strings.Index(string(data), "currentEnv: mockEnv") < 0 {
		t.Log("Current env is not mockEnv")
		t.Fail()
	}
}

func TestCommitCurrentEnv(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	os.Mkdir(filepath.Join(testingRoot, ".jorge"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge"))

	os.Mkdir(filepath.Join(testingRoot, ".jorge", "envs"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs"))

	os.Mkdir(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv"))

	os.WriteFile(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv", "mainTestConfig"), []byte("mock config contents"), 0600)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv", "mainTestConfig"))

	if err := os.WriteFile(filepath.Join(testingRoot, "mainTestConfig"), []byte("updated mock config contents"), 0600); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer os.Remove(filepath.Join(testingRoot, "mainTestConfig"))

	if err := os.WriteFile(filepath.Join(testingRoot, ".jorge", "config.yml"), []byte("currentEnv: mockEnv\nconfigFilePath: mainTestConfig"), 0600); err != nil {
		t.FailNow()
	}
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "config.yml"))

	if err := CommitCurrentEnv(); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if data, err := os.ReadFile(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv", "mainTestConfig")); err != nil {
		t.Log(err)
		t.FailNow()
	} else if strings.Index(string(data), "updated mock config contents") != 0 {
		t.Log("Config contents were not updated")
		t.Fail()
	}
}

func TestRestoreEnv(t *testing.T) {
	testingRoot := filepath.Join(os.TempDir(), "jorge-testing")
	os.Mkdir(testingRoot, 0700)
	defer os.Remove(testingRoot)
	os.Chdir(testingRoot)

	os.Mkdir(filepath.Join(testingRoot, ".jorge"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge"))

	os.Mkdir(filepath.Join(testingRoot, ".jorge", "envs"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs"))

	os.Mkdir(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv"), 0700)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv"))

	os.WriteFile(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv", "mainTestConfig"), []byte("mock config contents"), 0600)
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv", "mainTestConfig"))

	if err := os.WriteFile(filepath.Join(testingRoot, "mainTestConfig"), []byte("updated mock config contents"), 0600); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer os.Remove(filepath.Join(testingRoot, "mainTestConfig"))

	if err := os.WriteFile(filepath.Join(testingRoot, ".jorge", "config.yml"), []byte("currentEnv: mockEnv\nconfigFilePath: mainTestConfig"), 0600); err != nil {
		t.FailNow()
	}
	defer os.Remove(filepath.Join(testingRoot, ".jorge", "config.yml"))

	if err := RestoreEnv(); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if data, err := os.ReadFile(filepath.Join(testingRoot, ".jorge", "envs", "mockEnv", "mainTestConfig")); err != nil {
		t.Log(err)
		t.FailNow()
	} else if strings.Index(string(data), "mock config contents") != 0 {
		t.Log("Contents of active config were not updated")
		t.Fail()
	}
}

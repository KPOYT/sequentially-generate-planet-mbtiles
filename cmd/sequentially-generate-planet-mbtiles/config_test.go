package sequentiallygenerateplanetmbtiles

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func init() {

	tmpDir := filepath.Join(os.TempDir(), "sequentially-generate-planet-mbtiles")

	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		panic("Could not find os.TempDir. Permissions issue? TempleOS? Expect many tests to fail.")
	}

	cfg = &configuration{
		PlanetFile:       "", // CHANGE TO TEST FILE ONCE SET U
		DataDir:          tmpDir,
		OutDir:           filepath.Join(tmpDir, "out"),
		IncludeOcean:     true,
		IncludeLanduse:   true,
		TilemakerConfig:  "",
		TilemakerProcess: "",
		MaxRamMb:         1,
	}

	fmt.Printf("cfg initialised: %+v\n", cfg)
}

func TestSetConfigByJSON(t *testing.T) {
	// Test run in separate processes so as not to pollute other tests

	invalidConfig := "../../configs/test/TEST_INVALID.json"

	if _, err := os.Stat(invalidConfig); os.IsNotExist(err) {
		t.Errorf("TestSetConfigByJSON test config file does not exist")
	}

	if os.Getenv("TEST_SET_CONFIG_BY_JSON") == "1" {
		os.Args = append(os.Args, "-c", invalidConfig)
		flag.Parse()
		setConfigByJSON()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestSetConfigByJSON")
	cmd.Env = append(os.Environ(), "TEST_SET_CONFIG_BY_JSON=1")
	defer os.Unsetenv("TEST_SET_CONFIG_BY_JSON")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		if e.ExitCode() == exitInvalidJSON {
			return
		}
	}

	t.Fatalf("process ran with err %v, want exit status %v", err, exitInvalidJSON)

	// Add test for valid json
}

func TestSetConfigByFlags(t *testing.T) {

}

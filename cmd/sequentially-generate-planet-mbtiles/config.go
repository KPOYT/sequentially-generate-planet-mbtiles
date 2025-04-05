package sequentiallygenerateplanetmbtiles

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/KPOYT/lj_go/pkg/lj_json"
)

type configuration struct {
	srcFileProvided  bool
	PbfFile          string `json:"pbfFile"`
	WorkingDir       string `json:"workingDir"`
	OutDir           string `json:"outDir"`
	ExcludeOcean     bool   `json:"excludeOcean"`
	ExcludeLanduse   bool   `json:"excludeLanduse"`
	TilemakerConfig  string `json:"TilemakerConfig"`
	TilemakerProcess string `json:"TilemakerProcess"`
	MaxRamMb         uint64 `json:"maxRamMb"`
	OutAsDir         bool   `json:"outAsDir"`
	SkipSlicing      bool   `json:"skipSlicing"`
	MergeOnly        bool   `json:"mergeOnly"`
	SkipDownload     bool   `json:"skipDownload"`
}

func initConfig() {
	if fl.config == "" {
		setConfigByFlags()
	} else {
		setConfigByJSON()
	}

	verifyPaths()

	cfg.PbfFile = convertAbs(cfg.PbfFile)
	cfg.WorkingDir = convertAbs(cfg.WorkingDir)
	cfg.OutDir = convertAbs(cfg.OutDir)
}

func setConfigByJSON() {
	err := lj_json.DecodeTo(cfg, fl.config, 1000)
	if err != nil {
		log.Printf("Unable to decode config: %s", err)
		os.Exit(exitInvalidJSON)
	}
	if cfg.MaxRamMb == 0 {
		cfg.MaxRamMb = getRam()
	}
}

func setConfigByFlags() {
	cfg.PbfFile = fl.pbfFile
	cfg.WorkingDir = fl.workingDir
	cfg.OutDir = fl.outDir
	cfg.ExcludeOcean = fl.excludeOcean
	cfg.ExcludeLanduse = fl.excludeLanduse
	cfg.TilemakerConfig = fl.tilemakerConfig
	cfg.TilemakerProcess = fl.tilemakerProcess
	cfg.MaxRamMb = getRam()
	cfg.OutAsDir = fl.outAsDir
	cfg.SkipSlicing = fl.skipSlicing
	cfg.MergeOnly = fl.mergeOnly
	cfg.SkipDownload = fl.skipDownload
}

func setTmPaths() {
	if cfg.TilemakerProcess == "" || strings.ToLower(cfg.TilemakerProcess) == "tileserver-gl-basic" {
		cfg.TilemakerProcess = filepath.Join(pth.temp, "tilemaker", "resources", "process-openmaptiles.lua")
		log.Println("using tileserver-gl-basic style target")
	}

	if cfg.TilemakerProcess == "sgpm-bright" {
		cfg.TilemakerProcess = filepath.Join(pth.temp, "tilemaker", "resources", "process-bright.lua")
		log.Println("using sgpm-bright style target")
	}

	if cfg.TilemakerConfig == "" {
		cfg.TilemakerConfig = filepath.Join(pth.temp, "tilemaker", "resources", "config-openmaptiles.json")
	}

	cfg.TilemakerConfig = convertAbs(cfg.TilemakerConfig)
	cfg.TilemakerProcess = convertAbs(cfg.TilemakerProcess)
}

func verifyPaths() {
	if cfg.PbfFile != "" {
		if _, err := os.Stat(cfg.PbfFile); os.IsNotExist(err) {
			log.Fatalf("planet file does not exist: %s", cfg.PbfFile)
		}
		cfg.srcFileProvided = true
	} else {
		if fl.test {
			cfg.PbfFile = filepath.Join(cfg.WorkingDir, "pbf", "morocco-latest.osm.pbf")
		} else {
			cfg.PbfFile = filepath.Join(cfg.WorkingDir, "pbf", "planet-latest.osm.pbf")
		}
	}

	if cfg.TilemakerConfig != "" {
		if _, err := os.Stat(cfg.TilemakerConfig); os.IsNotExist(err) {
			log.Fatalf("tilemaker config does not exist: %s", cfg.TilemakerConfig)
		}
	}

	if cfg.TilemakerProcess != "" && cfg.TilemakerProcess != "sgpm-bright" && cfg.TilemakerProcess != "tileserver-gl-basic" {
		if _, err := os.Stat(cfg.TilemakerProcess); os.IsNotExist(err) {
			log.Fatalf("tilemaker process does not exist: %s", cfg.TilemakerProcess)
		}
	}
}

func getRam() uint64 {
	if fl.maxRamMb != 0 {
		return fl.maxRamMb
	}

	var memTotalKb uint64

	f, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Printf("Unable to open /proc/meminfo: %s", err)
		memTotalKb = 65536 * 1024
	} else {
		defer f.Close()

		scanner := bufio.NewScanner(f)

		memTotalStr := ""

		for scanner.Scan() {
			if strings.HasPrefix(scanner.Text(), "MemTotal:") {
				memTotalStr = regexp.MustCompile(`[0-9]+`).FindString(scanner.Text())
			}
		}

		memTotalKb, err = strconv.ParseUint(memTotalStr, 10, 64)
		if err != nil {
			log.Printf("Unable to parse MemTotal: %s", err)
			memTotalKb = 4096 * 1024
		}
	}

	if (memTotalKb / 1024) < 2048 {
		log.Printf("system ram is less than 2 GB; this is dangerously low, but we will do our best!")
		return 1024
	}

	memTotalMb := memTotalKb / 1024

	totalRamToUse := memTotalMb - 2048

	if totalRamToUse < 1024 {
		totalRamToUse = 1024
	}

	log.Printf("attempting to use not more than %d MB of ram", totalRamToUse)

	return totalRamToUse
}

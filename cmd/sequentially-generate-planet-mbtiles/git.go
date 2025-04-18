package sequentiallygenerateplanetmbtiles

import (
	"os"
	"path/filepath"

	"github.com/lambdajack/sequentially-generate-planet-mbtiles/internal/git"
	"github.com/lambdajack/sequentially-generate-planet-mbtiles/internal/system"
)

type repos struct {
	gdal       git.Repo
	osmiumTool git.Repo
	libosmium  git.Repo
	tilemaker  git.Repo
	tippecanoe git.Repo
}

var gh repos

func cloneRepos() {
	gh = repos{
		gdal: git.Repo{
			Url: "https://github.com/KPOYT/gdal",
			Tag: "v3.6.5",
			Dst: filepath.Join(pth.temp, "gdal"),
		},
		osmiumTool: git.Repo{
			Url: "https://github.com/KPOYT/osmium-tool",
			Tag: "v1.14.0",
			Dst: filepath.Join(pth.temp, "osmium", "osmium-tool"),
		},
		libosmium: git.Repo{
			Url: "https://github.com/KPOYT/libosmium",
			Tag: "v2.18.0",
			Dst: filepath.Join(pth.temp, "osmium", "libosmium"),
		},
		tilemaker: git.Repo{
			Url: "https://github.com/KPOYT/tilemaker",
			Tag: "v3.0.1",
			Dst: filepath.Join(pth.temp, "tilemaker"),
		},
		tippecanoe: git.Repo{
			Url: "https://github.com/KPOYT/tippecanoe",
			Tag: "1.37.0",
			Dst: filepath.Join(pth.temp, "tippecanoe"),
		},
	}

	var f []string

	err := gh.gdal.Clone()
	if err != nil {
		f = append(f, "gdal")
	}

	err = gh.osmiumTool.Clone()
	if err != nil {
		f = append(f, "osmium-tool")
	}

	err = gh.libosmium.Clone()
	if err != nil {
		f = append(f, "libosmium")
	}

	err = gh.tilemaker.Clone()
	if err != nil {
		f = append(f, "tilemaker")
	}

	err = gh.tippecanoe.Clone()
	if err != nil {
		f = append(f, "tippecanoe")
	}

	filepath.Walk(cfg.WorkingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		system.SetUserOwner(path)
		return nil
	})

	for _, e := range f {
		lg.err.Fatalf("error cloning %s: %v", e, err)
	}
}

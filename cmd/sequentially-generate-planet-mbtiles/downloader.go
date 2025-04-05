package sequentiallygenerateplanetmbtiles

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/KPOYT/lj_go/pkg/lj_http"
)

type downloadInformation struct {
	url, destFileName, destDir string
}

func downloadOsmData() {
	var downloads = []downloadInformation{
		{
			url:          "https://planet.openstreetmap.org/pbf/planet-latest.osm.pbf",
			destFileName: "planet-latest.osm.pbf",
			destDir:      pth.pbfDir,
		},
		{
			url:          "https://osmdata.openstreetmap.de/download/water-polygons-split-4326.zip",
			destFileName: "water-polygons-split-4326.zip",
			destDir:      pth.workingDir,
		},
		{
			url:          "https://naciscdn.org/naturalearth/10m/cultural/ne_10m_urban_areas.zip",
			destFileName: "ne_10m_urban_areas.zip",
			destDir:      pth.workingDir,
		},
		{
			url:          "https://naciscdn.org/naturalearth/10m/physical/ne_10m_antarctic_ice_shelves_polys.zip",
			destFileName: "ne_10m_antarctic_ice_shelves_polys.zip",
			destDir:      pth.workingDir,
		},
		{
			url:          "https://naciscdn.org/naturalearth/10m/physical/ne_10m_glaciated_areas.zip",
			destFileName: "ne_10m_glaciated_areas.zip",
			destDir:      pth.workingDir,
		},
	}

	if fl.test {
		downloads[0].url = "https://download.geofabrik.de/africa/morocco-latest.osm.pbf"
		downloads[0].destFileName = "morocco-latest.osm.pbf"
		lg.rep.Printf("test flag provided; downloading test data %s - %s", downloads[0].destFileName, downloads[0].url)
	}

	for _, dl := range downloads {
		if _, err := os.Stat(filepath.Join(dl.destDir, dl.destFileName)); os.IsNotExist(err) {
			if dl.destFileName == "planet-latest.osm.pbf" {
				if cfg.srcFileProvided {
					lg.rep.Printf("source file provided - skipping planet download %s", dl.url)
					continue
				}
				if cfg.SkipDownload {
					lg.rep.Printf("skip download flag provided - skipping planet download %s", dl.url)
					continue
				}
			}

			if cfg.ExcludeOcean && strings.Contains(dl.destFileName, "water") {
				lg.rep.Printf("skipping download of %s", dl.url)
				continue
			}

			if cfg.ExcludeLanduse && strings.Contains(dl.destFileName, "ne_") {
				lg.rep.Printf("skipping download of %s", dl.url)
				continue
			}

			err := lj_http.Download(dl.url, dl.destDir, dl.destFileName)
			if err != nil {
				lg.err.Printf("error downloading %s: %s", dl.url, err)
				os.Exit(exitDownloadURL)
			}
			lg.rep.Printf("Download success: %v\n", dl.destFileName)
		} else {
			lg.rep.Printf("%v already exists; skipping download.\n", dl.destFileName)
		}
	}
}

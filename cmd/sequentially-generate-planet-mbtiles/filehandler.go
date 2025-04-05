package sequentiallygenerateplanetmbtiles

import (
	"io/fs"
	"fmt"
	"os"
	"path/filepath"
)

func moveOcean() {
	if !cfg.ExcludeOcean {
		filepath.Walk(pth.coastlineDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("moveOcean() rename error: %s\n", err)
				return err
			}
			if !info.IsDir() {
				err := os.Rename(path, filepath.Join(pth.coastlineDir, filepath.Base(path)))
				if err != nil {
					fmt.Printf("moveOcean() rename error: %s\n", err)
					return err
				}
			}
			return nil
		})

		// Remove empty directories after Walk finishes
		filepath.Walk(pth.coastlineDir, func(path string, info fs.FileInfo, err error) error {
			if path == pth.coastlineDir {
				return nil
			}
			if err != nil {
				fmt.Printf("moveOcean() remove error: %s\n", err)
				return err
			}
			if info.IsDir() {
				err := os.Remove(path)
				if err != nil && !os.IsNotExist(err) {
					fmt.Printf("moveOcean() remove error: %s\n", err)
					return err
				}
			}
			return nil
		})
	}
}

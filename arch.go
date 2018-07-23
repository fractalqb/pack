package pack

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// ZipDist creates a ZIP file with filenam zipname. Extracting the package will
// create a folder with name dist. The contents of the distribution is taken
// from distDir.
func ZipDist(zipname string, dist, distDir string) error {
	zf, err := os.Create(zipname)
	if err != nil {
		return err
	}
	defer zf.Close()
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err = os.Chdir(distDir); err != nil {
		return err
	}
	defer os.Chdir(wd)
	zipw := zip.NewWriter(zf)
	defer zipw.Close()
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}
		zhdr, err := zip.FileInfoHeader(stat)
		if err != nil {
			return err
		}
		zhdr.Name = filepath.Join(dist, path)
		zwr, err := zipw.CreateHeader(zhdr)
		if err != nil {
			return err
		}
		rd, err := os.Open(path)
		if err != nil {
			return err
		}
		defer rd.Close()
		io.Copy(zwr, rd)
		return nil
	})
	return err
}

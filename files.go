package pack

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type OsDepNames = map[string]string

var OsDepExe = OsDepNames{
	"windows": "%s.exe",
}

// CopyFile copies file with path src to file with path dst. It also tarnsfers
// the file mode from the src file to the dst file.
func CopyFile(dst, src string, osdn OsDepNames) error {
	if osdn != nil {
		if pat, ok := osdn[runtime.GOOS]; ok {
			src = fmt.Sprintf(pat, src)
		}
	}
	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	log.Printf("copy '%s' → '%s'", src, dst)
	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, stat.Mode())
	return err
}

// CopyToDir copies a list of files to a single destination directory using
// CopyFile on each source file.
func CopyToDir(dst string, osdn OsDepNames, files ...string) error {
	for _, f := range files {
		b := filepath.Base(f)
		dst := filepath.Join(dst, b)
		err := CopyFile(dst, f, osdn)
		if err != nil {
			return err
		}
	}
	return nil
}

// CopyRecursive copies the content of the src directory into the dst directory.
// If filter is not nil only those files or subdirectories are copied for which
// the filter returns true.
func CopyRecursive(
	dst, src string,
	osdn OsDepNames,
	filter func(dir string, info os.FileInfo) bool,
) error {
	rddir, err := os.Open(src)
	if err != nil {
		log.Println(err)
	}
	defer rddir.Close()
	infos, err := rddir.Readdir(1)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		log.Println(err)
	}
	for len(infos) > 0 {
		info := infos[0]
		if filter == nil || filter(src, info) {
			switch {
			case info.IsDir():
				ddir := filepath.Join(dst, info.Name())
				err := os.Mkdir(ddir, 0777)
				if err != nil {
					return err
				}
				err = CopyRecursive(ddir, filepath.Join(src, info.Name()), osdn, filter)
				if err != nil {
					return err
				}
			case info.Mode().IsRegular():
				src := filepath.Join(src, info.Name())
				err := CopyToDir(dst, osdn, src)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("copy irregular file '%s/%s'", src, info.Name())
			}
		}
		if infos, err = rddir.Readdir(1); err != nil {
			if err == io.EOF {
				return nil
			}
			log.Println(err)
		}
	}
	return nil
}

// CopyTree copies the directory tree src into the dst directory, i.e. on succes
// the directory dst will contain one, perhaps new, directory src.
func CopyTree(
	dst, src string,
	osdn OsDepNames,
	filter func(dir string, info os.FileInfo) bool,
) error {
	dst = filepath.Join(dst, filepath.Base(src))
	err := os.Mkdir(dst, 0777)
	if err != nil {
		return err
	}
	err = CopyRecursive(dst, src, osdn, filter)
	return err
}

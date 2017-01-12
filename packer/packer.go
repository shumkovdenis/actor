package packer

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"io/ioutil"

	"github.com/emirpasic/gods/maps/treemap"
)

func Pack(oldPath, newPath, packPath string) error {
	oldFiles := treemap.NewWithStringComparator()
	if err := content(oldPath, oldFiles); err != nil {
		return err
	}

	newFiles := treemap.NewWithStringComparator()
	if err := content(newPath, newFiles); err != nil {
		return err
	}

	diffFiles, err := diff(oldPath, newPath, oldFiles, newFiles)
	if err != nil {
		return err
	}

	if err := pack(newPath, packPath, diffFiles); err != nil {
		return err
	}

	return nil
}

func Unpack(packPath, targetPath string) error {
	return unpack(packPath, targetPath)
}

func content(rootPath string, out *treemap.Map) error {
	root, err := os.Open(rootPath)
	if err != nil {
		return err
	}
	defer root.Close()

	files, err := root.Readdir(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(rootPath, file.Name())

		if file.Name() != ".DS_Store" {
			out.Put(path, file)
		}

		if file.IsDir() {
			if err := content(path, out); err != nil {
				return err
			}
		}
	}

	return nil
}

func diff(oldBasePath, newBasePath string, oldFiles, newFiles *treemap.Map) (*treemap.Map, error) {
	out := treemap.NewWithStringComparator()

	it := newFiles.Iterator()
	for it.Next() {
		newFilePath := it.Key().(string)
		newFileInfo := it.Value().(os.FileInfo)

		rel, err := filepath.Rel(newBasePath, newFilePath)
		if err != nil {
			return nil, err
		}

		oldFilePath := filepath.Join(oldBasePath, rel)

		found := oldFiles.Any(func(key, value interface{}) bool {
			return key == oldFilePath
		})

		if !found {
			out.Put(newFilePath, newFileInfo)
		} else {
			if newFileInfo.Mode()&os.ModeSymlink != 0 {
				equal, err := compareLinks(newFilePath, oldFilePath)
				if err != nil {
					return nil, err
				}
				if !equal {
					out.Put(newFilePath, newFileInfo)
				}
			} else if !newFileInfo.IsDir() {
				equal, err := compareFiles(newFilePath, oldFilePath)
				if err != nil {
					return nil, err
				}
				if !equal {
					out.Put(newFilePath, newFileInfo)
				}
			}
		}
	}

	return out, nil
}

func hash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func compareFiles(aFilePath, bFilePath string) (bool, error) {
	aFileHash, err := hash(aFilePath)
	if err != nil {
		return false, err
	}

	bFileHash, err := hash(bFilePath)
	if err != nil {
		return false, err
	}

	return (aFileHash == bFileHash), nil
}

func compareLinks(aLinkPath, bLinkPath string) (bool, error) {
	aLinkDest, err := os.Readlink(aLinkPath)
	if err != nil {
		return false, err
	}

	bLinkDest, err := os.Readlink(bLinkPath)
	if err != nil {
		return false, err
	}

	return aLinkDest == bLinkDest, nil
}

func pack(sourcePath, packPath string, files *treemap.Map) error {
	zipFile, err := os.Create(packPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	it := files.Iterator()
	for it.Next() {
		filePath := it.Key().(string)
		fileInfo := it.Value().(os.FileInfo)

		path, err := filepath.Rel(sourcePath, filePath)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}

		header.Name = path

		if fileInfo.IsDir() {
			header.Name += string(os.PathSeparator)
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			continue
		}

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			destLink, err := os.Readlink(filePath)

			if err != nil {
				return err
			}

			if _, err := writer.Write([]byte(destLink)); err != nil {
				return err
			}

			continue
		}

		if err := packFile(writer, filePath); err != nil {
			return err
		}
	}

	return nil
}

func packFile(writer io.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(writer, file); err != nil {
		return err
	}

	return nil
}

func unpack(packPath, targetPath string) error {
	zipReader, err := zip.OpenReader(packPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, zipFile := range zipReader.File {
		if err := unpackFile(zipFile, targetPath); err != nil {
			return err
		}
	}

	return nil
}

func unpackFile(zipFile *zip.File, targetPath string) error {
	path := filepath.Join(targetPath, zipFile.Name)

	if zipFile.FileInfo().IsDir() {
		os.MkdirAll(path, zipFile.Mode())
		return nil
	}

	reader, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	if zipFile.Mode()&os.ModeSymlink != 0 {
		destLink, err := ioutil.ReadAll(reader)

		if err != nil {
			return err
		}

		if err := os.Symlink(string(destLink), path); err != nil {
			return err
		}

		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return err
	}

	return nil
}

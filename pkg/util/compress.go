package util

import (
	"os"

	"github.com/mholt/archiver"
)

func Archive(dirPath, zipFileName string) error {
	err := archiver.Archive([]string{dirPath}, zipFileName)
	if err != nil {
		return err
	}

	return nil
}

func Unarchive(zipFileName, outputPath string) error {
	err := archiver.Unarchive(zipFileName, outputPath)
	if err != nil {
		return err
	}
	err = os.Remove(zipFileName)
	if err != nil {
		return err
	}

	return nil
}

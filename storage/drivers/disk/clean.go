package storage

import (
	"os"
	"strings"
	"time"

	godirwalk "github.com/karrick/godirwalk"
)

var (
	countProcessed = 0
)

func (s *Storage) cleanRoutine() error {
	options := &godirwalk.Options{
		Unsorted: true,
		Callback: cleanRoutine,
	}
	return godirwalk.Walk(s.path, options)
}

func cleanRoutine(osPathname string, de *godirwalk.Dirent) error {

	// stoper
	countProcessed++
	if countProcessed%100 == 0 {
		time.Sleep(10 * time.Millisecond)
	}

	if err := cleanEmptyDir(osPathname, de); err != nil {
		return err
	}
	return cleanAddled(osPathname, de)
}

func cleanEmptyDir(osPathname string, de *godirwalk.Dirent) error {
	// try to delete, if directory is empty it's can be deleted
	// consciously not doing slow Readdir to check if directory empty
	if de.IsDir() {
		os.Remove(osPathname)
	}
	return nil
}

func cleanAddled(osPathname string, de *godirwalk.Dirent) error {
	if de.IsDir() {
		return nil
	}
	if strings.HasSuffix(osPathname, headerExt) {
		// read and check
		header, err := parseHeader(osPathname)
		if err != nil {
			return nil
		}
		if !header.hasValidTTL() {
			os.RemoveAll(osPathname)
			os.RemoveAll(strings.TrimSuffix(osPathname, headerExt))
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil
}

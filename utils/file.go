package utils

import (
	"os"
	"path/filepath"
)

func CheckNumFolder(dir string) (int, error) {
	fileCount := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Kiểm tra nếu đó là file và không phải là thư mục
		if !info.IsDir() {
			fileCount++
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return fileCount, nil
}

package config

import (
	"log"
	"os"
)

func createFolder() {
	newFolderPending := "file/pending_file"
	if err := os.Mkdir(newFolderPending, 0075); err != nil {
		log.Println(err)
	}
	newFolderAddModel := "file/file_add_model"
	if err := os.Mkdir(newFolderAddModel, 0075); err != nil {
		log.Println(err)
	}
}

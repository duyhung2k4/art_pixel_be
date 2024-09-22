package config

import (
	"log"
	"os"
)

func createFolder() {
	if err := os.Mkdir("file/pending_file", 0075); err != nil {
		log.Println(err)
	}
	if err := os.Mkdir("file/file_add_model", 0075); err != nil {
		log.Println(err)
	}
	if err := os.Mkdir("file/auth_face", 0075); err != nil {
		log.Println(err)
	}
}

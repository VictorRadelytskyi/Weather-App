package utils

import (
	"mime"
)

func Init() {
	mime.AddExtensionType(".css", "text/css")
}

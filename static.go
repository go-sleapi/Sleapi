package sleapi

import (
	"fmt"
	"net/http"
	"os"
	"path"
	//"strings"
)

func Static(directory string) func(w ResponseWriter, req *http.Request) {
	//fmt.Println("Setup Static")

	//staticDir := http.Dir(directory)
	//prefix := directory

	return func(w ResponseWriter, req *http.Request) {
		//fmt.Println("Serving Static files")
		//fmt.Println("Serving Path: ", req.URL.Path)
		//fmt.Println("StaticDir: ", staticDir)
		pwd, _ := os.Getwd()
		//fmt.Println("Current Directory: ", pwd)

		reqPath := req.URL.Path

		if reqPath == "/" {
			reqPath = "static/index.html"
		}

		//trimmedPath := strings.TrimPrefix(path, prefix)
		trimmedPath := path.Join(pwd, reqPath)

		//fmt.Println("File to open: ", trimmedPath)
		//httpFile, error := staticDir.Open(trimmedPath)
		httpFile, error := os.Open(trimmedPath)
		if error != nil {
			fmt.Println("Error serving Static file: ", error.Error())
		}

		if httpFile != nil {
			//fmt.Println("HttpFile: ", httpFile)

			stat, err := httpFile.Stat()
			if err != nil {
				fmt.Println("Error getting Stat() on file: ", err.Error())
			}

			defer httpFile.Close()

			http.ServeContent(w, req, reqPath, stat.ModTime(), httpFile)
		}
	}
}

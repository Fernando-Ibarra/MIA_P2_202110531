package Comands

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func MakeRep() string {

	// current Path
	// currentPath, _ := os.Getwd()
	// diskPath := currentPath + "/MIA/P2/"
	path := "/home/fernando/Escritorio/"

	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Error al leer la carpeta ")
		return ""
	}

	jsonResponse := "["
	for _, report := range files {
		if strings.Contains(report.Name(), ".dot") || strings.Contains(report.Name(), ".txt") {

			pathFile := path + report.Name()
			bytesContent, errFile := ioutil.ReadFile(pathFile)
			if errFile != nil {
				log.Fatal(err)
			}

			res := string(bytesContent)

			filter1 := strings.ReplaceAll(res, "\n", "\\n")
			content := strings.ReplaceAll(filter1, `"`, `\"`)

			jsonResponse += "{"
			jsonResponse += "\"name\": " + "\"" + report.Name() + "\","
			jsonResponse += "\"content\": " + "\"" + content + "\""
			jsonResponse += "},"
		}
	}
	res := jsonResponse[:len(jsonResponse)-1]
	res += "]"

	fmt.Println(res)
	return res
}

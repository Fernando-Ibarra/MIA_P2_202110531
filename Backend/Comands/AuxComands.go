package Comands

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"os/exec"
	"strings"
	"unsafe"
)

func Compare(a string, b string) bool {
	if strings.ToUpper(a) == strings.ToUpper(b) {
		return true
	}
	return false
}

var responComand = ""

func Error(op string, message string, responseString string) {
	responseString += ""
	responseString += "\tERROR: " + op + "\n\tTIPO: " + message
	responComand += ""
	responComand += "\tERROR: " + op + "\n\tTIPO: " + message + "\n"
}

func ExistedFile(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func WritingBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func Message(op string, message string, responseString string) {
	responseString += ""
	responseString += "\tCOMANDO: " + op + "\n\tTIPO: " + message + "\n"
	responseString += ""
	responComand += ""
	responComand += "\tCOMANDO: " + op + "\n\tTIPO: " + message + "\n"
	responComand += ""
}

func readBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func readDisk(path string, responseString string) *Structs.MBR {
	mbr := Structs.MBR{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		Error("FDISK", "Error al abrir el archivo", responseString)
		return nil
	}
	file.Seek(0, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &mbr)
	if err_ != nil {
		Error("FDISK", "Error al leer el archivo", responseString)
		return nil
	}
	var mDir *Structs.MBR = &mbr
	return mDir
}

func CreateFile(nameFile string) {
	var _, err = os.Stat(nameFile)

	if os.IsNotExist(err) {
		var file, err_2 = os.Create(nameFile)
		if err_2 != nil {
			return
		}
		defer file.Close()
	}
}

func WriteFile(content string, nameFile string) {
	var file, err = os.OpenFile(nameFile, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return
	}
	err = file.Sync()
	if err != nil {
		return
	}
}

func Execute(nameFile string, file string, extension string, responseString string) {
	path, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(path, "-T"+extension, file).Output()
	mode := 0777
	_ = os.WriteFile(nameFile, cmd, os.FileMode(mode))
	Message("REP", "Archivo "+nameFile+", se ha generado correctamente", responseString)
}

func SetStringtoRes(text string) {
	responComand += text
}

func GetStirngRes() string {
	return responComand
}

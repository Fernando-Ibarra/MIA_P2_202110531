package Comands

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type ActiveUser struct {
	User     string
	Password string
	Id       string
	Uid      int
	Gid      int
}

var Logged ActiveUser

func DataUserLogin(context []string, responseString string) bool {
	id := ""
	user := ""
	pass := ""
	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "id") {
			id = tk[1]
		} else if Compare(tk[0], "user") {
			user = tk[1]
		} else if Compare(tk[0], "pass") {
			pass = tk[1]
		}
	}
	if id == "" || user == "" || pass == "" {
		Error("LOGIN", "Se necesitan parámetros obligatorios para el comando login", responseString)
		return false
	}
	return activeSession(user, pass, id, responseString)
}

func activeSession(u string, p string, id string, responseString string) bool {
	var path string
	partition := GetMount("LOGIN", id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("LOGIN", "No se encontro la partición montada con el id: "+id, responseString)
		return false
	}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("LOGIN", "No se ha encontrado el disco", responseString)
		return false
	}
	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("LOGIN", "Error al leer el archivo", responseString)
		return false
	}
	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("LOGIN", "Error al leer el archivo", responseString)
		return false
	}
	var fb Structs.FilesBlocks
	txt := ""
	for block := 1; block < 16; block++ {
		if inode.I_block[block-1] == -1 {
			break
		}
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(block-1), 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.FilesBlocks{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &fb)
		if err_ != nil {
			Error("LOGIN", "Error al leer el archivo", responseString)
			return false
		}
		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	vctr := strings.Split(txt, "\n")
	for i := 0; i < len(vctr)-1; i++ {
		line := vctr[i]
		if line[2] == 'U' || line[2] == 'u' {
			in := strings.Split(line, ",")
			if Compare(in[3], u) && Compare(in[4], p) && in[0] != "0" {
				idGroup := "0"
				exists := false
				for j := 0; j < len(vctr)-1; j++ {
					line2 := vctr[j]
					if (line2[2] == 'G' || line2[2] == 'g') && line2[0] != '0' {
						inG := strings.Split(line2, ",")
						if inG[2] == in[2] {
							idGroup = inG[0]
							exists = true
							break
						}
					}
				}
				if !exists {
					Error("LOGIN", "No se encontre el grupo \""+in[2]+"\".", responseString)
					return false
				}
				Message("LOGIN", "Logueado correctamente", responseString)
				responseString += "\t\t--------------------BIENVENIDO " + u + "--------------------"
				Logged.Id = id
				Logged.User = u
				Logged.Password = p
				Logged.Uid, _ = strconv.Atoi(in[0])
				Logged.Gid, _ = strconv.Atoi(idGroup)
				return true
			}
		}
	}
	Error("LOGIN", "No se encontró el usuario "+u, responseString)
	return false
}

func LogOut(responseString string) bool {
	Message("LOGOUT", "¡ADIOS "+Logged.User+" , espero volver a verte!", responseString)
	Logged = ActiveUser{}
	return false
}

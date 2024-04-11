package Comands

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func DataUser(context []string, action string, responseString string) {
	user := ""
	pass := ""
	grp := ""
	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "user") {
			user = tk[1]
		} else if Compare(tk[0], "pass") {
			pass = tk[1]
		} else if Compare(tk[0], "grp") {
			grp = tk[1]
		}
	}
	if Compare(action, "MK") {
		if user == "" || pass == "" || grp == "" {
			Error(action+"USER", "Se necesitan parámetros obligatorios para crear un usuario", responseString)
			return
		}
		if len(user) > 10 || len(pass) > 10 || len(grp) > 10 {
			Error(action+"USER", "La cantidad maxima de caracteres que se pueden usar son 10", responseString)
			return
		}
		mkuser(user, pass, grp, responseString)
	} else if Compare(action, "RM") {
		if user == "" {
			Error(action+"USER", "Se necesitan parametros obligatorios para eliminar un usuario", responseString)
			return
		}
		rmuser(user, responseString)
	} else if Compare(action, "CH") {
		if user == "" || grp == "" {
			Error(action+"GRP", "Se necesitan parametros obligatorios para cambiar el grupo de un usuario", responseString)
			return
		}
		chgrp(user, grp, responseString)
	} else {
		Error(action+"USER", "No se reconoce este comando", responseString)
		return
	}

}

func mkuser(user string, pass string, grp string, responseString string) {
	if !Compare(Logged.User, "root") {
		Error("MKUSR", "Solo el usuario \"root\" puede acceder a estos comandos", responseString)
		return
	}

	var path string
	partition := GetMount("MKGRP", Logged.Id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("MKUSR", "No se encontró la partición montada con el id: "+Logged.Id, responseString)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKUSR", "No se ha encontrado el disco", responseString)
		return
	}

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("MKUSR", "Error al leer el archivo", responseString)
		return
	}

	jour := Structs.NewJournaling()
	jourW := Structs.NewJournaling()
	var posJour int64
	if Compare(strconv.Itoa(int(super.S_filesystem_type)), "3") {
		for i := 0; i < int(super.S_inodes_count); i++ {
			file.Seek(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBlock{}))+int64(unsafe.Sizeof(Structs.Journaling{}))*int64(i), 0)
			posJour = partition.Part_start + int64(unsafe.Sizeof(Structs.SuperBlock{})) + int64(unsafe.Sizeof(Structs.Journaling{}))*int64(i)
			data = readBytes(file, int(unsafe.Sizeof(Structs.Journaling{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &jour)
			if err_ != nil {
				Error("MKUSR", "Error al leer el archivo", responseString)
				return
			}

			pathJournaling := ""
			for k := 0; k < len(jour.Path); k++ {
				if jour.Path[k] != 0 {
					pathJournaling += string(jour.Path[k])
				}
			}

			if Compare(pathJournaling, "-") {
				contentU := user + " - " + pass
				dateU := time.Now().String()
				copy(jourW.Operation[:], "mkusr")
				copy(jourW.Path[:], "users.txt")
				copy(jourW.Content[:], contentU)
				copy(jourW.Date[:], dateU)
				file.Close()

				file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
				if err != nil {
					Error("MKUSR", "No se ha encontrado el disco", responseString)
					return
				}
				file.Seek(posJour, 0)
				var binJu bytes.Buffer
				binary.Write(&binJu, binary.BigEndian, jourW)
				WritingBytes(file, binJu.Bytes())
				file.Close()
				break
			} else {
				continue
			}
		}
	}

	file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKUSR", "No se ha encontrado el disco", responseString)
		return
	}

	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("MKUSER", "Error al leer el archivo", responseString)
		return
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
			Error("MKUSER", "Error al leer el archivo", responseString)
			return
		}
		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	vctr := strings.Split(txt, "\n")
	exists := false
	for i := 0; i < len(vctr)-1; i++ {
		line := vctr[i]
		if (line[2] == 'G' || line[2] == 'g') && line[0] != '0' {
			in := strings.Split(line, ",")
			if in[2] == grp {
				exists = true
				break
			}
		}
	}

	if !exists {
		Error("MKUSER", "No se encontro el grupo \""+grp+"\".", responseString)
		return
	}

	c := 0
	for i := 0; i < len(vctr)-1; i++ {
		line := vctr[i]
		if line[2] == 'U' || line[2] == 'u' {
			c++
			in := strings.Split(line, ",")
			if in[3] == user {
				if line[0] != '0' {
					Error("MKUSER", "El nombre "+user+", ya esta en uso", responseString)
					return
				}
			}
		}
	}

	txt += strconv.Itoa(c+1) + ",U," + grp + "," + user + "," + pass + "\n"
	tam := len(txt)
	var cadenaS []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadenaS = append(cadenaS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadenaS = append(cadenaS, txt)
		}
	} else {
		cadenaS = append(cadenaS, txt)
	}

	if len(cadenaS) > 16 {
		Error("MKUSER", "Se ha llenado la cantidad de archivos posibles y no se puede generar más", responseString)
		return
	}

	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("MKUSER", "No se ha encontrado el disco", responseString)
		return
	}

	for i := 0; i < len(cadenaS); i++ {
		var fbAux Structs.FilesBlocks
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			WritingBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadenaS[i])
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
		var bin1 bytes.Buffer
		binary.Write(&bin1, binary.BigEndian, fbAux)
		WritingBytes(file, bin1.Bytes())
	}
	for i := 0; i < len(cadenaS); i++ {
		inode.I_block[i] = int64(i)
	}
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	WritingBytes(file, inodos.Bytes())

	Message("MKUSER", "Usuario "+user+", creado correctamente", responseString)
	file.Close()
}

func rmuser(user string, responseString string) {
	if !Compare(Logged.User, "root") {
		Error("RMUSR", "Solo el usuario \"root\" puede acceder a estos comandos", responseString)
		return
	}

	var path string
	partition := GetMount("MKGRP", Logged.Id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("RMUSR", "No se encontró la partición montada con el id: "+Logged.Id, responseString)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("RMUSR", "No se ha encontrado el disco", responseString)
		return
	}

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("RMUSR", "Error al leer el archivo", responseString)
		return
	}

	jour := Structs.NewJournaling()
	jourW := Structs.NewJournaling()
	var posJour int64
	if Compare(strconv.Itoa(int(super.S_filesystem_type)), "3") {
		for i := 0; i < int(super.S_inodes_count); i++ {
			file.Seek(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBlock{}))+int64(unsafe.Sizeof(Structs.Journaling{}))*int64(i), 0)
			posJour = partition.Part_start + int64(unsafe.Sizeof(Structs.SuperBlock{})) + int64(unsafe.Sizeof(Structs.Journaling{}))*int64(i)
			data = readBytes(file, int(unsafe.Sizeof(Structs.Journaling{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &jour)
			if err_ != nil {
				Error("RMUSR", "Error al leer el archivo", responseString)
				return
			}

			pathJournaling := ""
			for k := 0; k < len(jour.Path); k++ {
				if jour.Path[k] != 0 {
					pathJournaling += string(jour.Path[k])
				}
			}

			if Compare(pathJournaling, "-") {
				contentU := user
				dateU := time.Now().String()
				copy(jourW.Operation[:], "mkusr")
				copy(jourW.Path[:], "users.txt")
				copy(jourW.Content[:], contentU)
				copy(jourW.Date[:], dateU)
				file.Close()

				file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
				if err != nil {
					Error("RMUSR", "No se ha encontrado el disco", responseString)
					return
				}
				file.Seek(posJour, 0)
				var binJu bytes.Buffer
				binary.Write(&binJu, binary.BigEndian, jourW)
				WritingBytes(file, binJu.Bytes())
				file.Close()
				break
			} else {
				continue
			}
		}
	}

	file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKGRP", "No se ha encontrado el disco", responseString)
		return
	}

	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("RMUSER", "Error al leer el archivo", responseString)
		return
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
			Error("RMUSER", "Error al leer el archivo", responseString)
			return
		}
		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	aux := ""

	vctr := strings.Split(txt, "\n")
	exists := false
	for i := 0; i < len(vctr)-1; i++ {
		line := vctr[i]
		if (line[2] == 'U' || line[2] == 'u') && line[0] != '0' {
			in := strings.Split(line, ",")
			if in[3] == user {
				exists = true
				aux += strconv.Itoa(0) + ",U," + in[2] + "," + in[3] + "," + in[4] + "\n"
				continue
			}
		}
		aux += line + "\n"
	}

	if !exists {
		Error("MKUSER", "No se encontro el usuario \""+user+"\".", responseString)
		return
	}

	txt = aux
	tam := len(txt)
	var cadenaS []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadenaS = append(cadenaS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadenaS = append(cadenaS, txt)
		}
	} else {
		cadenaS = append(cadenaS, txt)
	}

	if len(cadenaS) > 16 {
		Error("RMUSER", "Se ha llenado la cantidad de archivos posibles y no se puede generar más", responseString)
		return
	}

	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("RMUSER", "No se ha encontrado el disco", responseString)
		return
	}

	for i := 0; i < len(cadenaS); i++ {
		var fbAux Structs.FilesBlocks
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			WritingBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadenaS[i])
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
		var bin1 bytes.Buffer
		binary.Write(&bin1, binary.BigEndian, fbAux)
		WritingBytes(file, bin1.Bytes())
	}
	for i := 0; i < len(cadenaS); i++ {
		inode.I_block[i] = int64(i)
	}
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	WritingBytes(file, inodos.Bytes())

	Message("RMUSER", "Usuario "+user+", eliminado correctamente", responseString)
	file.Close()
}

func chgrp(user string, grp string, responseString string) {
	if !Compare(Logged.User, "root") {
		Error("CHGRP", "Solo el usuario \"root\" puede acceder a estos comandos", responseString)
		return
	}

	var path string
	partition := GetMount("MKGRP", Logged.Id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("CHGRP", "No se encontró la partición montada con el id: "+Logged.Id, responseString)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("CHGRP", "No se ha encontrado el disco", responseString)
		return
	}

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("CHGRP", "Error al leer el archivo", responseString)
		return
	}

	jour := Structs.NewJournaling()
	jourW := Structs.NewJournaling()
	var posJour int64
	if Compare(strconv.Itoa(int(super.S_filesystem_type)), "3") {
		for i := 0; i < int(super.S_inodes_count); i++ {
			file.Seek(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBlock{}))+int64(unsafe.Sizeof(Structs.Journaling{}))*int64(i), 0)
			posJour = partition.Part_start + int64(unsafe.Sizeof(Structs.SuperBlock{})) + int64(unsafe.Sizeof(Structs.Journaling{}))*int64(i)
			data = readBytes(file, int(unsafe.Sizeof(Structs.Journaling{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &jour)
			if err_ != nil {
				Error("CHGRP", "Error al leer el archivo", responseString)
				return
			}

			pathJournaling := ""
			for k := 0; k < len(jour.Path); k++ {
				if jour.Path[k] != 0 {
					pathJournaling += string(jour.Path[k])
				}
			}

			if Compare(pathJournaling, "-") {
				operation := "chgrp"
				pathU := "users.txt"
				contentU := grp
				dateU := time.Now().String()
				copy(jourW.Operation[:], operation)
				copy(jourW.Path[:], pathU)
				copy(jourW.Content[:], contentU)
				copy(jourW.Date[:], dateU)
				file.Close()

				file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
				if err != nil {
					Error("CHGRP", "No se ha encontrado el disco", responseString)
					return
				}
				file.Seek(posJour, 0)
				var binJu bytes.Buffer
				binary.Write(&binJu, binary.BigEndian, jourW)
				WritingBytes(file, binJu.Bytes())
				file.Close()
				break
			} else {
				continue
			}
		}
	}

	file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("CHGRP", "No se ha encontrado el disco", responseString)
		return
	}

	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("CHGRP", "Error al leer el archivo", responseString)
		return
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
			Error("CHGRP", "Error al leer el archivo", responseString)
			return
		}
		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	aux := ""
	vctr := strings.Split(txt, "\n")
	exists := false
	for i := 0; i < len(vctr)-1; i++ {
		line := vctr[i]
		if (line[2] == 'G' || line[2] == 'g') && line[0] != 0 {
			in := strings.Split(line, ",")
			if in[2] == grp {
				exists = true
				continue
			}
		}
		aux += line + "\n"
	}
	if !exists {
		Error("CHGRP", "No se encontró \""+grp+"\".", responseString)
		return
	}

	aux = ""
	vctr = strings.Split(txt, "\n")
	exists = false
	for i := 0; i < len(vctr)-1; i++ {
		line := vctr[i]
		if (line[2] == 'U' || line[2] == 'u') && line[0] != '0' {
			in := strings.Split(line, ",")
			if in[3] == user {
				exists = true
				aux += in[0] + ",U," + grp + "," + in[3] + "," + in[4] + "\n"
				continue
			}
		}
		aux += line + "\n"
	}

	if !exists {
		Error("CHGRP", "No se encontro el usuario \""+user+"\".", responseString)
		return
	}

	txt = aux
	tam := len(txt)
	var cadenaS []string
	if tam > 64 {
		for tam > 64 {
			aux = ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadenaS = append(cadenaS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadenaS = append(cadenaS, txt)
		}
	} else {
		cadenaS = append(cadenaS, txt)
	}

	if len(cadenaS) > 16 {
		Error("CHGRP", "Se ha llenado la cantidad de archivos posibles y no se puede generar más", responseString)
		return
	}

	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("CHGRP", "No se ha encontrado el disco", responseString)
		return
	}

	for i := 0; i < len(cadenaS); i++ {
		var fbAux Structs.FilesBlocks
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			WritingBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadenaS[i])
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
		var bin1 bytes.Buffer
		binary.Write(&bin1, binary.BigEndian, fbAux)
		WritingBytes(file, bin1.Bytes())
	}
	for i := 0; i < len(cadenaS); i++ {
		inode.I_block[i] = int64(i)
	}
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	WritingBytes(file, inodos.Bytes())

	Message("CHGRP", "Usuario "+user+", se ha cambiado al grupo "+grp+" correctamente", responseString)
	file.Close()
}

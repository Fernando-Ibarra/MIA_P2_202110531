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

func DataGroup(context []string, action string, responseString string) {
	name := ""
	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "name") {
			name = tk[1]
		}
	}
	if name == "" {
		Error("MKGROUP", "No se encontro el parámetro name en el comando", responseString)
		return
	}
	if Compare(action, "MK") {
		mkgrp(name, responseString)
	} else if Compare(action, "RM") {
		rmgrp(name, responseString)
	} else {
		Error(action+"GRP", "No se reconoce este comando", responseString)
	}
}

func mkgrp(n string, responseString string) {
	if !Compare(Logged.User, "root") {
		Error("MKGRP", "Solo el usuario \"root\" puede acceder a estos comandos", responseString)
		return
	}

	var path string
	partition := GetMount("MKGRP", Logged.Id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("MKGRP", "No se encontró la partición montada con el id: "+Logged.Id, responseString)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKGRP", "No se ha encontrado el disco", responseString)
		return
	}

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("MKGRP", "Error al leer el archivo", responseString)
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
				Error("MKDIR", "Error al leer el archivo", responseString)
				return
			}

			pathJournaling := ""
			for k := 0; k < len(jour.Path); k++ {
				if jour.Path[k] != 0 {
					pathJournaling += string(jour.Path[k])
				}
			}

			if Compare(pathJournaling, "-") {
				operation := "mkgrp"
				pathU := "users.txt"
				contentU := n
				dateU := time.Now().String()
				copy(jourW.Operation[:], operation)
				copy(jourW.Path[:], pathU)
				copy(jourW.Content[:], contentU)
				copy(jourW.Date[:], dateU)
				file.Close()

				file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
				if err != nil {
					Error("MKDIR", "No se ha encontrado el disco", responseString)
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
		Error("MKGRP", "Error al leer el archivo", responseString)
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
			Error("MKGRP", "Error al leer el archivo", responseString)
			return
		}
		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	vctr := strings.Split(txt, "\n")
	c := 0
	for i := 0; i < len(vctr)-1; i++ {
		line := vctr[i]
		if line[2] == 'G' || line[2] == 'g' {
			c++
			in := strings.Split(line, ",")
			if in[2] == n {
				if line[0] != 0 {
					Error("MKGRP", "EL nombre "+n+", ya esta en uso", responseString)
				}
			}
		}
	}

	txt += strconv.Itoa(c+1) + ",G," + n + "\n"

	tam := len(txt)
	var cadS []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadS = append(cadS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadS = append(cadS, txt)
		}
	} else {
		cadS = append(cadS, txt)
	}
	if len(cadS) > 16 {
		Error("MKGRP", "Se ha llenado la cantidad de archivo sposibles y no se puede generar más", responseString)
		return
	}
	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("MKGRP", "No se ha encontrado el disco", responseString)
		return
	}

	for i := 0; i < len(cadS); i++ {
		var fbAux Structs.FilesBlocks
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			WritingBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadS[i])

		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
		var bin1 bytes.Buffer
		binary.Write(&bin1, binary.BigEndian, fbAux)
		WritingBytes(file, bin1.Bytes())
	}

	for i := 0; i < len(cadS); i++ {
		inode.I_block[i] = int64(i)
	}
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var ino bytes.Buffer
	binary.Write(&ino, binary.BigEndian, inode)
	WritingBytes(file, ino.Bytes())

	Message("MKGRP", "Grupo "+n+", creado correctamente", responseString)
	file.Close()
}

func rmgrp(n string, responseString string) {
	if !Compare(Logged.User, "root") {
		Error("RMGRP", "Solo el usuario \"root\" puede acceder a estos comandos", responseString)
		return
	}

	var path string
	partition := GetMount("MKGRP", Logged.Id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("RMGRP", "No se encontró la partición montada con el id: "+Logged.Id, responseString)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("RMGRP", "No se ha encontrado el disco", responseString)
		return
	}

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("RMGRP", "Error al leer el archivo", responseString)
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
				Error("RMGRP", "Error al leer el archivo", responseString)
				return
			}

			pathJournaling := ""
			for k := 0; k < len(jour.Path); k++ {
				if jour.Path[k] != 0 {
					pathJournaling += string(jour.Path[k])
				}
			}

			if Compare(pathJournaling, "-") {
				operation := "rmgrp"
				pathU := "users.txt"
				contentU := n
				dateU := time.Now().String()
				copy(jourW.Operation[:], operation)
				copy(jourW.Path[:], pathU)
				copy(jourW.Content[:], contentU)
				copy(jourW.Date[:], dateU)
				file.Close()

				file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
				if err != nil {
					Error("RMGRP", "No se ha encontrado el disco", responseString)
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
		Error("RMGRP", "No se ha encontrado el disco", responseString)
		return
	}

	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("RMGRP", "Error al leer el archivo", responseString)
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
			Error("MKGRP", "Error al leer el archivo", responseString)
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
			if in[2] == n {
				exists = true
				aux += strconv.Itoa(0) + ",G," + in[2] + "\n"
				continue
			}
		}
		aux += line + "\n"
	}
	if !exists {
		Error("RMGRP", "No se encontró \""+n+"\".", responseString)
		return
	}
	txt = aux

	tam := len(txt)
	var cadS []string
	if tam > 64 {
		for tam > 64 {
			aux = ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadS = append(cadS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadS = append(cadS, txt)
		}
	} else {
		cadS = append(cadS, txt)
	}
	if len(cadS) > 16 {
		Error("RMGRP", "Se ha llenado la cantidad de archivo posibles y no se puede generar más", responseString)
		return
	}
	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("MKGRP", "No se ha encontrado el disco", responseString)
		return
	}

	for i := 0; i < len(cadS); i++ {
		var fbAux Structs.FilesBlocks
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			WritingBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadS[i])

		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(i), 0)
		var bin1 bytes.Buffer
		binary.Write(&bin1, binary.BigEndian, fbAux)
		WritingBytes(file, bin1.Bytes())
	}

	for i := 0; i < len(cadS); i++ {
		inode.I_block[i] = int64(i)
	}
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var ino bytes.Buffer
	binary.Write(&ino, binary.BigEndian, inode)
	WritingBytes(file, ino.Bytes())

	Message("RMGRP", "Grupo "+n+", eliminado correctamente", responseString)
	file.Close()
}

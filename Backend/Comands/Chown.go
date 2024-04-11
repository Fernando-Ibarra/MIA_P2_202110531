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

func DataChown(context []string, partition Structs.Partition, pth string, responseString string) {
	rBoolean := false
	user := ""
	path := ""
	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "path") {
			path = tk[1]
		} else if Compare(tk[0], "r") {
			rBoolean = true
		} else if Compare(tk[0], "user") {
			user = tk[1]
		}
	}
	if path == "" || user == "" {
		Error("CHMOD", "Se necesitan parametros obligatorios para crear un directorio", responseString)
		return
	}
	tmp := GetPath(path)
	chown(tmp, rBoolean, partition, pth, user, responseString)
}

func chown(path []string, r bool, partition Structs.Partition, pth string, u string, responseString string) {

	gid, uid := getUser(u, Logged.Id, responseString)
	copyPath := path
	spr := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))
	if err != nil {
		Error("CHOWN", "No se ha encontrado el disco", responseString)
		return
	}

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &spr)
	if err_ != nil {
		Error("CHOWN", "Error al leer el archivo", responseString)
		return
	}

	file.Seek(spr.S_inode_start, 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("CHOWN", "Error al leer el archivo", responseString)
		return
	}

	var fileToFound string
	if len(path) == 0 {
		Error("CHOWN", "No se ha brindado un path valido", responseString)
		return
	}

	var bi int64
	fnd := false
	inodetmp := Structs.NewInodos()

	fileToFound = path[len(path)-1]
	var aux []string
	for i := 0; i < len(path); i++ {
		aux = append(aux, path[i])
	}
	path = aux
	var stack string

	for v := 0; v < len(path)-1; v++ {
		fnd = false
		for i := 0; i < 16; i++ {
			if i < 16 {
				if inode.I_block[i] != -1 {
					folder = Structs.NewDirectoriesBlocks()
					file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
					data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
					buffer = bytes.NewBuffer(data)
					err_ = binary.Read(buffer, binary.BigEndian, &folder)
					if err_ != nil {
						Error("CHOWN", "Error al leer el archivo", responseString)
						return
					}

					for j := 0; j < 4; j++ {
						nameFolder := ""
						for nam := 0; nam < len(folder.B_content[j].B_name); nam++ {
							if folder.B_content[j].B_name[nam] == 0 {
								continue
							}
							nameFolder += string(folder.B_content[j].B_name[nam])
						}
						if Compare(nameFolder, path[v]) {
							stack += "/" + path[v]
							fnd = true
							inode = Structs.NewInodos()
							// inodo padre de la carpeta donde esta el archivo
							file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*folder.B_content[j].B_inodo, 0)
							data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
							buffer = bytes.NewBuffer(data)
							err_ = binary.Read(buffer, binary.BigEndian, &inode)
							if err_ != nil {
								Error("CHOWN", "Error al leer el archivo", responseString)
								return
							}

							if inode.I_uid != int64(Logged.Uid) {
								Error("CHOWN", "No tiene permisos para crear carpetas en este directorio", responseString)
								return
							}

							break
						}
					}

				} else {
					break
				}
			}
		}
		if !fnd {
			address := ""
			for i := 0; i < len(path); i++ {
				address += "/" + path[i]
			}
			Error("CHOWN", "No se pudo crear el directorio "+address+", no existen directorios", responseString)
			return
		}
	}

	fnd = false
	for i := 0; i < 16; i++ {
		if inode.I_block[i] != -1 {
			if i < 16 {
				// carpeta donde va tendria que estar el archivo
				file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
				data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
				buffer = bytes.NewBuffer(data)
				err_ = binary.Read(buffer, binary.BigEndian, &folder)
				if err_ != nil {
					Error("CHOWN", "Error al leer el archivo", responseString)
					return
				}

				if folder.B_content[2].B_inodo != -1 {
					nameAux1 := ""
					for nam := 0; nam < len(folder.B_content[2].B_name); nam++ {
						if folder.B_content[2].B_name[nam] == 0 {
							continue
						}
						nameAux1 += string(folder.B_content[2].B_name[nam])
					}
					if Compare(nameAux1, fileToFound) {
						bi = folder.B_content[2].B_inodo
						file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
						data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &inodetmp)
						if err_ != nil {
							Error("CHOWN", "Error al leer el archivo", responseString)
							return
						}

						uidInt, _ := strconv.Atoi(uid)
						gidInt, _ := strconv.Atoi(gid)
						inodetmp.I_uid = int64(uidInt)
						inodetmp.I_gid = int64(gidInt)

						file.Close()
						file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
						if err != nil {
							Error("CHOWN", "No se ha encontrado el disco", responseString)
							return
						}

						file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
						var binInodeTmp bytes.Buffer
						binary.Write(&binInodeTmp, binary.BigEndian, inodetmp)
						WritingBytes(file, binInodeTmp.Bytes())

						file.Close()
						file, err = os.Open(strings.ReplaceAll(pth, "\"", ""))
						if err != nil {
							Error("CHOWN", "No se ha encontrado el disco", responseString)
							return
						}

						if r {
							for w := 0; w < 16; w++ {
								if inodetmp.I_block[w] != -1 {
									control := inodetmp.I_block[w]
									recurChown(control, partition, pth, gid, uid, responseString)
								} else {
									break
								}
							}
						}
						ruta := ""
						for p := 0; p < len(copyPath); p++ {
							ruta += "/" + copyPath[p]
						}
						Message("CHOWN", "Se ha cambiado el propietario de la carpeta "+ruta, responseString)
						break
					}
				}

				if folder.B_content[3].B_inodo != -1 {
					nameAux2 := ""
					for nam := 0; nam < len(folder.B_content[3].B_name); nam++ {
						if folder.B_content[3].B_name[nam] == 0 {
							continue
						}
						nameAux2 += string(folder.B_content[3].B_name[nam])
					}
					if Compare(nameAux2, fileToFound) {
						bi = folder.B_content[3].B_inodo
						file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
						data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &inodetmp)
						if err_ != nil {
							Error("CHOWN", "Error al leer el archivo", responseString)
							return
						}

						uidInt, _ := strconv.Atoi(uid)
						gidInt, _ := strconv.Atoi(gid)
						inodetmp.I_uid = int64(uidInt)
						inodetmp.I_gid = int64(gidInt)

						file.Close()
						file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
						if err != nil {
							Error("CHOWN", "No se ha encontrado el disco", responseString)
							return
						}

						file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
						var binInodeTmp bytes.Buffer
						binary.Write(&binInodeTmp, binary.BigEndian, inodetmp)
						WritingBytes(file, binInodeTmp.Bytes())

						file.Close()
						file, err = os.Open(strings.ReplaceAll(pth, "\"", ""))
						if err != nil {
							Error("CHOWN", "No se ha encontrado el disco", responseString)
							return
						}

						if r {
							for w := 0; w < 16; w++ {
								if inodetmp.I_block[w] != -1 {
									control := inodetmp.I_block[w]
									recurChown(control, partition, pth, gid, uid, responseString)
								} else {
									break
								}
							}
						}
						ruta := ""
						for p := 0; p < len(copyPath); p++ {
							ruta += "/" + copyPath[p]
						}
						file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
						Message("CHOWN", "Se ha cambiado el propietario de la carpeta "+ruta, responseString)
						break
					}
				}
			}
		}
	}
}

func getUser(u string, id string, responseString string) (string, string) {
	var path string
	partition := GetMount("LOGIN", id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("CHOWN", "No se encontro la partición montada con el id: "+id, responseString)
		return "", ""
	}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("CHOWN", "No se ha encontrado el disco", responseString)
		return "", ""
	}
	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("CHOWN", "Error al leer el archivo", responseString)
		return "", ""
	}
	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("CHOWN", "Error al leer el archivo", responseString)
		return "", ""
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
			Error("CHOWN", "Error al leer el archivo", responseString)
			return "", ""
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
			if Compare(in[3], u) && in[0] != "0" {
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
					Error("CHOWN", "No se encontre el grupo \""+in[2]+"\".", responseString)
					return "", ""
				}
				return idGroup, in[0]
			}
		}
	}
	Error("CHOWN", "No se encontró el usuario "+u, responseString)
	return "", ""
}

func recurChown(control int64, partition Structs.Partition, pth string, gid string, uid string, responseString string) {
	spr := Structs.NewSuperBlock()
	inode := Structs.NewInodos()

	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))
	if err != nil {
		Error("CHMOD", "No se ha encontrado el disco", responseString)
		return
	}

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &spr)
	if err_ != nil {
		Error("CHMOD", "Error al leer el archivo", responseString)
		return
	}

	directory := Structs.NewDirectoriesBlocks()
	file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*control+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*control, 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &directory)

	if directory.B_content[2].B_inodo != -1 {
		nameAux1 := ""
		for nam := 0; nam < len(directory.B_content[2].B_name); nam++ {
			if directory.B_content[2].B_name[nam] == 0 {
				continue
			}
			nameAux1 += string(directory.B_content[2].B_name[nam])
		}
		inode = Structs.NewInodos()
		bi := directory.B_content[2].B_inodo
		file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("CHMOD", "Error al leer el archivo", responseString)
			return
		}

		uidInt, _ := strconv.Atoi(uid)
		gidInt, _ := strconv.Atoi(gid)
		inode.I_uid = int64(uidInt)
		inode.I_gid = int64(gidInt)

		file.Close()
		file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
		if err != nil {
			Error("CHMOD", "No se ha encontrado el disco", responseString)
			return
		}

		file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
		var binInodeTmp bytes.Buffer
		binary.Write(&binInodeTmp, binary.BigEndian, inode)
		WritingBytes(file, binInodeTmp.Bytes())

		file.Close()
		file, err = os.Open(strings.ReplaceAll(pth, "\"", ""))
		if err != nil {
			Error("CHMOD", "No se ha encontrado el disco", responseString)
			return
		}

		if strings.Contains(nameAux1, ".") {
		} else {
			for w := 0; w < 16; w++ {
				if inode.I_block[w] != -1 {
					controlInt := inode.I_block[w]
					recurChown(controlInt, partition, pth, gid, uid, responseString)
				} else {
					break
				}
			}
		}

	}

	if directory.B_content[3].B_inodo != -1 {
		bi := directory.B_content[3].B_inodo
		nameAux1 := ""
		for nam := 0; nam < len(directory.B_content[3].B_name); nam++ {
			if directory.B_content[3].B_name[nam] == 0 {
				continue
			}
			nameAux1 += string(directory.B_content[3].B_name[nam])
		}
		inode = Structs.NewInodos()
		file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("CHOWN", "Error al leer el archivo", responseString)
			return
		}

		uidInt, _ := strconv.Atoi(uid)
		gidInt, _ := strconv.Atoi(gid)
		inode.I_uid = int64(uidInt)
		inode.I_gid = int64(gidInt)

		file.Close()
		file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
		if err != nil {
			Error("CHMOD", "No se ha encontrado el disco", responseString)
			return
		}
		file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
		var binInodeTmp bytes.Buffer
		binary.Write(&binInodeTmp, binary.BigEndian, inode)
		WritingBytes(file, binInodeTmp.Bytes())
		file.Close()
		file, err = os.Open(strings.ReplaceAll(pth, "\"", ""))
		if err != nil {
			Error("CHMOD", "No se ha encontrado el disco", responseString)
			return
		}

		if strings.Contains(nameAux1, ".") {
		} else {
			for w := 0; w < 16; w++ {
				if inode.I_block[w] != -1 {
					controlInt := inode.I_block[w]
					recurChown(controlInt, partition, pth, gid, uid, responseString)
				} else {
					break
				}
			}
		}
	}
	return
}

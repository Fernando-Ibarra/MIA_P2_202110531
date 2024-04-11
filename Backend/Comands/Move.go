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

func DataMove(context []string, partition Structs.Partition, pth string, responseString string) {
	destino := ""
	path := ""

	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "path") {
			path = tk[1]
		} else if Compare(tk[0], "destino") {
			destino = tk[1]
		}
	}
	if path == "" || destino == "" {
		Error("MKDIR", "Se necesitan parametros obligatorios para crear un directorio", responseString)
		return
	}
	tmp := GetPath(path)
	tmp2 := GetPath(destino)
	move(tmp, tmp2, partition, pth, responseString)
}

func move(path []string, dest []string, partition Structs.Partition, pth string, responseString string) {
	bi, nameFolderGet := getBi(path, partition, pth, responseString)
	copyPath := dest
	spr := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))
	if err != nil {
		Error("MOVE", "No se ha encontrado el disco", responseString)
		return
	}

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &spr)
	if err_ != nil {
		Error("MOVE", "Error al leer el archivo", responseString)
		return
	}

	jour := Structs.NewJournaling()
	jourW := Structs.NewJournaling()
	var posJour int64

	if Compare(strconv.Itoa(int(spr.S_filesystem_type)), "3") {
		rutaFs3_2 := ""
		for i := 0; i < len(copyPath); i++ {
			rutaFs3_2 += "/" + copyPath[i]
		}

		rutaFs3_1 := ""
		for i := 0; i < len(path); i++ {
			rutaFs3_1 += "/" + copyPath[i]
		}

		for i := 0; i < int(spr.S_inodes_count); i++ {
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
				operation := "move"
				pathU := rutaFs3_1
				contentU := rutaFs3_2
				dateU := time.Now().String()
				copy(jourW.Operation[:], operation)
				copy(jourW.Path[:], pathU)
				copy(jourW.Content[:], contentU)
				copy(jourW.Date[:], dateU)
				file.Close()

				file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
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

	file, err = os.Open(strings.ReplaceAll(pth, "\"", ""))
	if err != nil {
		Error("MOVE", "No se ha encontrado el disco", responseString)
		return
	}

	file.Seek(spr.S_inode_start, 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("MOVE", "Error al leer el archivo", responseString)
		return
	}

	if len(dest) == 0 {
		Error("MOVE", "No se ha brindado un path valido", responseString)
		return
	}

	fnd := false
	var father int64
	var past int64

	var aux []string
	for i := 0; i < len(dest); i++ {
		aux = append(aux, dest[i])
	}
	dest = aux
	var stack string

	for v := 0; v < len(dest); v++ {
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
						Error("MOVE", "Error al leer el archivo", responseString)
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
						if Compare(nameFolder, dest[v]) {
							stack += "/" + dest[v]
							fnd = true
							father = folder.B_content[j].B_inodo
							inode = Structs.NewInodos()
							file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*folder.B_content[j].B_inodo, 0)
							data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
							buffer = bytes.NewBuffer(data)
							err_ = binary.Read(buffer, binary.BigEndian, &inode)
							if err_ != nil {
								Error("MOVE", "Error al leer el archivo", responseString)
								return
							}

							if inode.I_uid != int64(Logged.Uid) {
								Error("MOVE", "No tiene permisos para crear carpetas en este directorio", responseString)
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
			for i := 0; i < len(dest); i++ {
				address += "/" + dest[i]
			}
			Error("MOVE", "No se pudo crear el directorio "+address+", no existen directorios", responseString)
			return
		}
	}

	/*
		Por si el padre tiene una carpeta donde hay un espacio libre para
	*/
	fnd = false
	for i := 0; i < 16; i++ {
		if inode.I_block[i] != -1 {
			if i < 16 {
				folderAux := folder
				file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
				data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
				buffer = bytes.NewBuffer(data)
				err_ = binary.Read(buffer, binary.BigEndian, &folder)
				if err_ != nil {
					Error("MOVE", "Error al leer el archivo", responseString)
					return
				}
				nameAux1 := ""
				for nam := 0; nam < len(folder.B_content[2].B_name); nam++ {
					if folder.B_content[2].B_name[nam] == 0 {
						continue
					}
					nameAux1 += string(folder.B_content[2].B_name[nam])
				}

				nameAux2 := ""
				for nam := 0; nam < len(folderAux.B_content[2].B_name); nam++ {
					if folderAux.B_content[2].B_name[nam] == 0 {
						continue
					}
					nameAux2 += string(folderAux.B_content[2].B_name[nam])
				}
				padre := ""
				for k := 0; k < len(path); k++ {
					if k >= 1 {
						padre = path[k-1]
					}
				}
				if padre == nameAux1 {
					continue
				}
				for j := 0; j < 4; j++ {
					if folder.B_content[j].B_inodo == -1 {

						folder.B_content[j].B_inodo = bi
						copy(folder.B_content[j].B_name[:], nameFolderGet)
						file.Close()

						file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
						if err != nil {
							Error("MKDIR", "No se ha encontrado el disco", responseString)
							return
						}
						file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
						var binFolder bytes.Buffer
						binary.Write(&binFolder, binary.BigEndian, folder)
						WritingBytes(file, binFolder.Bytes())

						fnd = true
						i = 20
						address := ""
						for p := 0; p < len(path); p++ {
							address += "/" + path[p]
						}

						address2 := ""
						for p := 0; p < len(dest); p++ {
							address2 += "/" + dest[p]
						}

						Message("MOVE", "Se ha movido "+address2+" a "+address, responseString)
						break
					}
				}
			}
		} else {
			break
		}
	}

	if !fnd {
		for i := 0; i < 16; i++ {
			if inode.I_block[i] == -1 {
				if i < 16 {
					past = GetFree(spr, pth, "BB", responseString)
					if past == -1 {
						Error("MOVE", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad", responseString)
						return
					}

					inode.I_block[i] = past
					folder = Structs.NewDirectoriesBlocks()
					copy(folder.B_content[0].B_name[:], ".")
					folder.B_content[0].B_inodo = past
					copy(folder.B_content[1].B_name[:], "..")
					folder.B_content[1].B_inodo = father
					folder.B_content[2].B_inodo = bi
					copy(folder.B_content[2].B_name[:], nameFolderGet)
					copy(folder.B_content[3].B_name[:], "-")

					file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
					if err != nil {
						Error("MOVE", "No se ha encontrado el disco", responseString)
						return
					}
					file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*father, 0)
					var binInodo bytes.Buffer
					binary.Write(&binInodo, binary.BigEndian, inode)
					WritingBytes(file, binInodo.Bytes())

					file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*past+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*past, 0)
					var binFolder bytes.Buffer
					binary.Write(&binFolder, binary.BigEndian, folder)
					WritingBytes(file, binFolder.Bytes())

					file.Close()

					address := ""
					for p := 0; p < len(path); p++ {
						address += "/" + path[p]
					}

					address2 := ""
					for p := 0; p < len(dest); p++ {
						address2 += "/" + dest[p]
					}

					Message("MOVE", "Se ha movido "+address2+" a "+address, responseString)
					break
				}
			}
		}
	}

}

func getBi(path []string, partition Structs.Partition, pth string, responseString string) (int64, string) {
	spr := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))
	if err != nil {
		Error("MOVE", "No se ha encontrado el disco", responseString)
		return -1, ""
	}

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &spr)
	if err_ != nil {
		Error("MOVE", "Error al leer el archivo", responseString)
		return -1, ""
	}

	file.Seek(spr.S_inode_start, 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("MOVE", "Error al leer el archivo", responseString)
		return -1, ""
	}

	var fileToFound string
	if len(path) == 0 {
		Error("MOVE", "No se ha brindado un path valido", responseString)
		return -1, ""
	}

	var bi int64
	fnd := false
	nameFolderReturn := ""

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
						Error("MOVE", "Error al leer el archivo", responseString)
						return -1, ""
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
								Error("MOVE", "Error al leer el archivo", responseString)
								return -1, ""
							}

							if inode.I_uid != int64(Logged.Uid) {
								Error("MOVE", "No tiene permisos para crear carpetas en este directorio", responseString)
								return -1, ""
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
			Error("MOVE", "No se pudo crear el directorio "+address+", no existen directorios", responseString)
			return -1, ""
		}
	}

	fnd = false
	for i := 0; i < 16; i++ {
		if inode.I_block[i] != -1 {
			if i < 16 {
				file, err = os.Open(strings.ReplaceAll(pth, "\"", ""))
				if err != nil {
					Error("MOVE", "No se ha encontrado el disco", responseString)
					return -1, ""
				}
				// carpeta donde va tendria que estar el archivo
				file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
				data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
				buffer = bytes.NewBuffer(data)
				err_ = binary.Read(buffer, binary.BigEndian, &folder)
				if err_ != nil {
					Error("MOVE", "Error al leer el archivo", responseString)
					return -1, ""
				}
				for j := 0; j < 4; j++ {
					if folder.B_content[j].B_inodo != -1 {

						nameAux1 := ""
						for nam := 0; nam < len(folder.B_content[2].B_name); nam++ {
							if folder.B_content[2].B_name[nam] == 0 {
								continue
							}
							nameAux1 += string(folder.B_content[2].B_name[nam])
						}

						nameAux2 := ""
						for nam := 0; nam < len(folder.B_content[3].B_name); nam++ {
							if folder.B_content[3].B_name[nam] == 0 {
								continue
							}
							nameAux2 += string(folder.B_content[3].B_name[nam])
						}

						if Compare(nameAux1, fileToFound) {
							// Get inodo id
							bi = folder.B_content[2].B_inodo
							nameFolderReturn = nameAux1

							folder.B_content[2].B_inodo = -1
							for nam := 0; nam < len(folder.B_content[2].B_name); nam++ {
								if folder.B_content[2].B_name[nam] == 0 {
									continue
								}
								folder.B_content[2].B_name[nam] = 0
							}

							copy(folder.B_content[2].B_name[:], "-")

							file.Close()
							file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
							if err != nil {
								Error("MOVE", "No se ha encontrado el disco", responseString)
								return -1, ""
							}

							file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
							var binFile bytes.Buffer
							binary.Write(&binFile, binary.BigEndian, folder)
							WritingBytes(file, binFile.Bytes())
							file.Close()

							break
						} else if Compare(nameAux2, fileToFound) {
							bi = folder.B_content[3].B_inodo
							nameFolderReturn = nameAux2

							folder.B_content[3].B_inodo = -1
							for nam := 0; nam < len(folder.B_content[3].B_name); nam++ {
								if folder.B_content[3].B_name[nam] == 0 {
									continue
								}
								folder.B_content[3].B_name[nam] = 0
							}
							copy(folder.B_content[3].B_name[:], "-")

							file.Close()
							file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
							if err != nil {
								Error("MOVE", "No se ha encontrado el disco", responseString)
								return -1, ""
							}

							file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
							var binFile bytes.Buffer
							binary.Write(&binFile, binary.BigEndian, folder)
							WritingBytes(file, binFile.Bytes())
							file.Close()

							break
						}
					} else {
						break
					}
				}
			}
		}
	}
	return bi, nameFolderReturn
}

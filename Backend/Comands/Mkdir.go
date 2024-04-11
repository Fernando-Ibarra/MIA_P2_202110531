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

func DataDir(context []string, partition Structs.Partition, pth string, responseString string) {
	rBoolean := false
	path := ""

	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "path") {
			path = tk[1]
		} else if Compare(tk[0], "r") {
			rBoolean = true
		}
	}
	if path == "" {
		Error("MKDIR", "Se necesitan parametros obligatorios para crear un directorio", responseString)
		return
	}
	tmp := GetPath(path)
	mkdir(tmp, rBoolean, partition, pth, responseString)
}

func mkdir(path []string, r bool, partition Structs.Partition, pth string, responseString string) {
	copyPath := path
	spr := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))
	if err != nil {
		Error("MKDIR", "No se ha encontrado el disco", responseString)
		return
	}

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &spr)
	if err_ != nil {
		Error("MKDIR", "Error al leer el archivo", responseString)
		return
	}

	jour := Structs.NewJournaling()
	jourW := Structs.NewJournaling()
	var posJour int64
	if Compare(strconv.Itoa(int(spr.S_filesystem_type)), "3") {
		rutaFs3 := ""
		for i := 0; i < len(copyPath); i++ {
			rutaFs3 += "/" + copyPath[i]
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
				operation := "mkdir"
				pathU := rutaFs3
				contentU := "-"
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
		Error("MKDIR", "No se ha encontrado el disco", responseString)
		return
	}

	file.Seek(spr.S_inode_start, 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("MKDIR", "Error al leer el archivo", responseString)
		return
	}

	var newf string
	if len(path) == 0 {
		Error("MKDIR", "No se ha brindado un path valido", responseString)
		return
	}

	var past int64
	var bi int64
	var bb int64
	fnd := false
	inodetmp := Structs.NewInodos()
	foldertmp := Structs.NewDirectoriesBlocks()

	newf = path[len(path)-1]
	var father int64

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
						Error("MKDIR", "Error al leer el archivo", responseString)
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
							father = folder.B_content[j].B_inodo
							inode = Structs.NewInodos()
							file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*folder.B_content[j].B_inodo, 0)

							data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
							buffer = bytes.NewBuffer(data)
							err_ = binary.Read(buffer, binary.BigEndian, &inode)
							if err_ != nil {
								Error("MKDIR", "Error al leer el archivo", responseString)
								return
							}

							if inode.I_uid != int64(Logged.Uid) {
								Error("MKDIR", "No tiene permisos para crear carpetas en este directorio", responseString)
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
			if r {
				stack += "/" + path[v]
				mkdir(GetPath(stack), false, partition, pth, responseString)
				file.Seek(spr.S_inode_start, 0)

				data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
				buffer = bytes.NewBuffer(data)
				err_ = binary.Read(buffer, binary.BigEndian, &inode)
				if err_ != nil {
					Error("MKDIR", "Error al leer el archivo", responseString)
					return
				}

				if v == len(path)-2 {
					stack += "/" + path[v+1]
					mkdir(GetPath(stack), false, partition, pth, responseString)
					file.Seek(spr.S_inode_start, 0)
					data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
					buffer = bytes.NewBuffer(data)
					err_ = binary.Read(buffer, binary.BigEndian, &inode)
					if err_ != nil {
						Error("MKDIR", "Error al leer el archivo", responseString)
						return
					}
					return
				}
			} else {
				address := ""
				for i := 0; i < len(path); i++ {
					address += "/" + path[i]
				}
				Error("MKDIR", "No se pudo crear el directorio "+address+", no existen directorios", responseString)
				return
			}
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
					Error("MKDIR", "Error al leer el archivo", responseString)
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
						past = inode.I_block[i]
						bi = GetFree(spr, pth, "BI", responseString)
						if bi == -1 {
							Error("MKDIR", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad", responseString)
							return
						}
						bb = GetFree(spr, pth, "BB", responseString)
						if bb == -1 {
							Error("MKDIR", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad", responseString)
							return
						}

						inodetmp.I_uid = int64(Logged.Uid)
						inodetmp.I_gid = int64(Logged.Gid)
						inodetmp.I_s = int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))

						dateNow := time.Now().String()
						copy(inodetmp.I_atime[:], spr.S_mtime[:])
						copy(inodetmp.I_ctime[:], dateNow)
						copy(inodetmp.I_mtime[:], dateNow)
						inodetmp.I_type = 0
						inodetmp.I_perm = 664
						inodetmp.I_block[0] = bb

						copy(foldertmp.B_content[0].B_name[:], ".")
						foldertmp.B_content[0].B_inodo = bi
						copy(foldertmp.B_content[1].B_name[:], "..")
						foldertmp.B_content[1].B_inodo = father
						copy(foldertmp.B_content[2].B_name[:], "-")
						copy(foldertmp.B_content[3].B_name[:], "-")

						folder.B_content[j].B_inodo = bi
						copy(folder.B_content[j].B_name[:], newf)
						fnd = true
						i = 20
						break
					}
				}
			}
		} else {
			break
		}
	}
	/*
	 Encontrando un espacio donde se puede escribir el bloque/carpeta si hay espacio libre en un carpeta padre. Se usa un nuevo inodo
	*/
	if !fnd {
		for i := 0; i < 16; i++ {
			if inode.I_block[i] == -1 {
				if i < 16 {
					bi = GetFree(spr, pth, "BI", responseString)
					if bi == -1 {
						Error("MKDIR", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad", responseString)
						return
					}
					past = GetFree(spr, pth, "BB", responseString)
					if past == -1 {
						Error("MKDIR", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad", responseString)
						return
					}

					bb = GetFree(spr, pth, "BB", responseString)

					folder = Structs.NewDirectoriesBlocks()
					copy(folder.B_content[0].B_name[:], ".")
					folder.B_content[0].B_inodo = bi
					copy(folder.B_content[1].B_name[:], "..")
					folder.B_content[1].B_inodo = father
					folder.B_content[2].B_inodo = bi
					copy(folder.B_content[2].B_name[:], newf)
					copy(folder.B_content[3].B_name[:], "-")

					inodetmp.I_uid = int64(Logged.Uid)
					inodetmp.I_gid = int64(Logged.Gid)
					inodetmp.I_s = int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))

					dateNow := time.Now().String()
					copy(inodetmp.I_atime[:], spr.S_mtime[:])
					copy(inodetmp.I_ctime[:], dateNow)
					copy(inodetmp.I_mtime[:], dateNow)
					inodetmp.I_type = 0
					inodetmp.I_perm = 664
					inodetmp.I_block[0] = bb

					copy(foldertmp.B_content[0].B_name[:], ".")
					foldertmp.B_content[0].B_inodo = bi
					copy(foldertmp.B_content[1].B_name[:], ".")
					foldertmp.B_content[1].B_inodo = father
					copy(foldertmp.B_content[2].B_name[:], "-")
					copy(foldertmp.B_content[3].B_name[:], "-")
					file.Close()

					copy(folder.B_content[2].B_name[:], newf)
					inode.I_block[i] = past
					file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
					if err != nil {
						Error("MKDIR", "No se ha encontrado el disco", responseString)
						return
					}
					file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*father, 0)
					var binInodo bytes.Buffer
					binary.Write(&binInodo, binary.BigEndian, inode)
					WritingBytes(file, binInodo.Bytes())
					file.Close()
					break
				}
			}
		}
	}

	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("MKDIR", "No se ha encontrado el disco", responseString)
		return
	}

	file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
	var binInodeTmp bytes.Buffer
	binary.Write(&binInodeTmp, binary.BigEndian, inodetmp)
	WritingBytes(file, binInodeTmp.Bytes())

	file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*bb+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*bb, 0)
	var binFolderTmp bytes.Buffer
	binary.Write(&binFolderTmp, binary.BigEndian, foldertmp)
	WritingBytes(file, binFolderTmp.Bytes())

	file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*past+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*past, 0)
	var binFolder bytes.Buffer
	binary.Write(&binFolder, binary.BigEndian, folder)
	WritingBytes(file, binFolder.Bytes())

	updateBm(spr, pth, "BI", responseString)
	updateBm(spr, pth, "BB", responseString)

	ruta := ""
	for i := 0; i < len(copyPath); i++ {
		ruta += "/" + copyPath[i]
	}
	Message("MKDIR", "Se ha creado el directorio "+ruta, responseString)
	file.Close()
}

func updateBm(spr Structs.SuperBlock, pth string, t string, responseString string) {
	ch := 'x'
	var num int

	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))
	if err != nil {
		Error("MKDIR", "No se ha encontrado el disco", responseString)
		return
	}
	if t == "BI" {
		file.Seek(spr.S_bm_inode_start, 0)
		for i := 0; i < int(spr.S_inodes_count); i++ {
			data := readBytes(file, int(unsafe.Sizeof(ch)))
			buffer := bytes.NewBuffer(data)
			err_ := binary.Read(buffer, binary.BigEndian, &ch)
			if err_ != nil {
				Error("MKDIR", "Error al leer el archivo", responseString)
				return
			}
			if ch == '0' {
				num = i
				break
			}
		}
		file.Close()

		file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
		if err != nil {
			Error("MKDIR", "No se ha encontrado el disco", responseString)
			return
		}
		zero := '1'
		file.Seek(spr.S_bm_inode_start, 0)
		for i := 0; i < num+1; i++ {
			var binaryZero bytes.Buffer
			binary.Write(&binaryZero, binary.BigEndian, zero)
			WritingBytes(file, binaryZero.Bytes())
		}
		file.Close()
	} else {
		file.Seek(spr.S_bm_block_start, 0)
		for i := 0; i < int(spr.S_blocks_count); i++ {
			data := readBytes(file, int(unsafe.Sizeof(ch)))
			buffer := bytes.NewBuffer(data)
			err_ := binary.Read(buffer, binary.BigEndian, &ch)
			if err_ != nil {
				Error("MKDIR", "Error al leer el archivo", responseString)
				return
			}
			if ch == '0' {
				num = i
				break
			}
		}
		file.Close()

		file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
		if err != nil {
			Error("MKDIR", "No se ha encontrado el disco", responseString)
			return
		}
		zero := '1'
		file.Seek(spr.S_bm_block_start, 0)
		for i := 0; i < num+1; i++ {
			var binaryZero bytes.Buffer
			binary.Write(&binaryZero, binary.BigEndian, zero)
			WritingBytes(file, binaryZero.Bytes())
		}
		file.Close()
	}
	file.Close()
}

func GetPath(path string) []string {
	var result []string
	if path == "" {
		return result
	}
	aux := strings.Split(path, "/")
	for i := 1; i < len(aux); i++ {
		result = append(result, aux[i])
	}
	return result
}

func GetFree(spr Structs.SuperBlock, pth string, t string, responseString string) int64 {
	ch := '2'
	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))
	if err != nil {
		Error("MKDIR", "No se ha encontrado el disco", responseString)
		return -1
	}
	if t == "BI" {
		file.Seek(spr.S_bm_inode_start, 0)
		for i := 0; i < int(spr.S_inodes_count); i++ {
			data := readBytes(file, int(unsafe.Sizeof(ch)))
			buffer := bytes.NewBuffer(data)
			err_ := binary.Read(buffer, binary.BigEndian, &ch)
			if err_ != nil {
				Error("MKDIR", "Error al leer el archivo", responseString)
				return -1
			}
			if ch == '0' {
				file.Close()
				return int64(i)
			}
		}
	} else {
		file.Seek(spr.S_bm_block_start, 0)
		for i := 0; i < int(spr.S_blocks_count); i++ {
			data := readBytes(file, int(unsafe.Sizeof(ch)))
			buffer := bytes.NewBuffer(data)
			err_ := binary.Read(buffer, binary.BigEndian, &ch)
			if err_ != nil {
				Error("MKDIR", "Error al leer el archivo", responseString)
				return -1
			}
			if ch == '0' {
				file.Close()
				return int64(i)
			}
		}
	}
	return -1
}

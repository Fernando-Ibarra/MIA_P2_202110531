package Comands

import (
	"Backend/Structs"
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func DataFile(context []string, partition Structs.Partition, pth string, responseString string) {
	rBoolean := false
	size := ""
	path := ""
	cont := ""

	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "path") {
			path = tk[1]
		} else if Compare(tk[0], "r") {
			rBoolean = true
		} else if Compare(tk[0], "size") {
			size = tk[1]
		} else if Compare(tk[0], "cont") {
			cont = tk[1]
		}
	}
	if path == "" {
		Error("MKFILE", "Se necesitan parametros obligatorios para crear un directorio", responseString)
		return
	}
	if cont == "" {
		cont = ""
	}
	if size == "" {
		size = "0"
	}
	tmp := GetPath(path)
	mkfile(tmp, rBoolean, partition, pth, size, cont, responseString)
}

func mkfile(path []string, r bool, partition Structs.Partition, pth string, s string, cont string, responseString string) {
	copyPath := path
	size, err := strconv.Atoi(s)
	spr := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	if size < 0 {
		Error("MKFILE", "NO SE PUEDEN CREAR ARCHIVOS DE TAMAÑO NEGATIVO", responseString)
		return
	}

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

	counterFilesBlockse := 0
	var past int64
	var bi int64
	var bb int64
	fnd := false
	inodetmp := Structs.NewInodos()
	foldertmp := Structs.NewDirectoriesBlocks()
	fileBlock := Structs.FilesBlocks{}

	newf = path[len(path)-1]
	var father int64
	var fatherFileFolder int64
	var aux []string
	for i := 0; i < len(path); i++ {
		aux = append(aux, path[i])
	}
	path = aux
	var stack string
	fileWritten := false
	fatherSpace := false

	for v := 0; v < len(path)-1; v++ {
		fnd = false
		for i := 0; i < 16; i++ {
			if i < 16 {
				if inode.I_block[i] != -1 {
					folder = Structs.NewDirectoriesBlocks()
					file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
					fatherFileFolder = inode.I_block[i]
					data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
					buffer = bytes.NewBuffer(data)
					err_ = binary.Read(buffer, binary.BigEndian, &folder)
					if err_ != nil {
						Error("MKFILE", "Error al leer el archivo", responseString)
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
							// inodo padre de la carpeta donde esta el archivo
							file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*folder.B_content[j].B_inodo, 0)
							data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
							buffer = bytes.NewBuffer(data)
							err_ = binary.Read(buffer, binary.BigEndian, &inode)
							if err_ != nil {
								Error("MKFILE", "Error al leer el archivo", responseString)
								return
							}

							if inode.I_uid != int64(Logged.Uid) {
								Error("MKFILE", "No tiene permisos para crear carpetas en este directorio", responseString)
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
					Error("MKFILE", "Error al leer el archivo", responseString)
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
						Error("MKFILE", "Error al leer el archivo", responseString)
						return
					}
					return
				}
			} else {
				address := ""
				for i := 0; i < len(path); i++ {
					address += "/" + path[i]
				}
				Error("MKFILE", "No se pudo crear el archivo "+address+", no existen directorios", responseString)
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
				// carpeta donde va tendria que estar el archivo
				fatherFileFolder = inode.I_block[i] // folder aux
				file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
				data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
				buffer = bytes.NewBuffer(data)
				err_ = binary.Read(buffer, binary.BigEndian, &folder)
				if err_ != nil {
					Error("MKFILE", "Error al leer el archivo", responseString)
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
							Error("MKFILE", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad", responseString)
							return
						}
						bb = GetFree(spr, pth, "BB", responseString)
						if bb == -1 {
							Error("MKFILE", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad", responseString)
							return
						}

						if strings.Contains(newf, ".") {
							inodetmp.I_uid = int64(Logged.Uid)
							inodetmp.I_gid = int64(Logged.Gid)

							dateNow := time.Now().String()
							copy(inodetmp.I_atime[:], spr.S_mtime[:])
							copy(inodetmp.I_ctime[:], dateNow)
							copy(inodetmp.I_mtime[:], dateNow)
							inodetmp.I_type = 1
							inodetmp.I_perm = 664

							folder.B_content[j].B_inodo = bi
							copy(folder.B_content[j].B_name[:], newf)

							content := ""
							variableContent := ""
							if cont != "" {
								dataFil, errC := os.Open(cont)
								if errC != nil {
									Error("MKFILE", "LA RUTA DEL ARCHIVO A LEER NO EXISTE", responseString)
									return
								}

								bufferR := bufio.NewReader(dataFil)
								for {
									line, errL := bufferR.ReadString('\n')
									if errL != nil {
										if errL == io.EOF {
											break
										}
										return
									}
									variableContent += line
								}
								dataFil.Close()

								for p := 0; p < size; p++ {
									content += string(variableContent[p])
								}
							} else {
								internalCounter := 0
								for d := 0; d < size; d++ {
									if d%10 == 0 {
										internalCounter = 0
										content += strconv.Itoa(internalCounter)
										internalCounter++
									} else {
										content += strconv.Itoa(internalCounter)
										internalCounter++
									}
								}
							}

							inodetmp.I_s = int64(unsafe.Sizeof(Structs.DirectoriesBlocks{})) + int64(len(content))

							tam := len(content)
							var filesBlocks []string
							if tam > 10 {
								for tam >= 10 {
									auxFiles := ""
									for m := 0; m <= 10; m++ {
										auxFiles += string(content[m])
									}
									filesBlocks = append(filesBlocks, auxFiles)
									content = strings.ReplaceAll(content, auxFiles, "")
									tam = len(content)
								}
								if tam <= 10 && tam != 0 {
									filesBlocks = append(filesBlocks, content)
								}
							} else {
								filesBlocks = append(filesBlocks, content)
							}

							if len(filesBlocks) > 16 {
								Error("MKFILE", "SE HA LLEGADO A LA MÁXIMA CAPACIDAD DE APUNDODORES DEL INODO PADRE", responseString)
								return
							}
							file.Close()

							file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
							if err != nil {
								Error("MKFILE", "No se ha encontrado el disco", responseString)
								return
							}

							for w := 0; w < len(filesBlocks); w++ {
								var control int
								control = int(bb) + w
								var fbAux Structs.FilesBlocks
								if inode.I_block[w] == -1 {
									file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(control), 0)
									var binAux bytes.Buffer
									binary.Write(&binAux, binary.BigEndian, fbAux)
									WritingBytes(file, binAux.Bytes())
								} else {
									fbAux = fileBlock
								}

								copy(fbAux.B_content[:], filesBlocks[w])
								file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(control), 0)
								var bin1 bytes.Buffer
								binary.Write(&bin1, binary.BigEndian, fbAux)
								WritingBytes(file, bin1.Bytes())
								counterFilesBlockse++
							}

							for w := 0; w < len(filesBlocks); w++ {
								inodetmp.I_block[w] = bb + int64(w)
							}

							fileWritten = true
							fatherSpace = true
							fnd = true
							i = 20
							break
						} else {
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
						Error("MKFILE", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad", responseString)
						return
					}
					past = GetFree(spr, pth, "BB", responseString)
					if past == -1 {
						Error("MKFILE", "No se ha podido crear el directorio, el sistema de archivos ha alcanzado su maxima capacidad", responseString)
						return
					}

					bb = GetFree(spr, pth, "BB", responseString)

					if strings.Contains(newf, ".") {
						// folder previo donde se ingresara el archivo
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

						dateNow := time.Now().String()
						copy(inodetmp.I_atime[:], spr.S_mtime[:])
						copy(inodetmp.I_ctime[:], dateNow)
						copy(inodetmp.I_mtime[:], dateNow)
						inodetmp.I_type = 1
						inodetmp.I_perm = 664

						// ARCHIVO MAYOR A 10 CARACTERES
						content := ""
						variableContent := ""
						if cont != "" {
							dataFil, errC := os.Open(cont)
							if errC != nil {
								Error("MKFILE", "LA RUTA DEL ARCHIVO A LEER NO EXISTE", responseString)
								return
							}

							bufferR := bufio.NewReader(dataFil)
							for {
								line, errL := bufferR.ReadString('\n')
								if errL != nil {
									if errL == io.EOF {
										break
									}
									return
								}
								variableContent += line
							}
							dataFil.Close()

							for p := 0; p < size; p++ {
								content += string(variableContent[p])
							}
						} else {
							internalCounter := 0
							for d := 0; d < size; d++ {
								if d%10 == 0 {
									internalCounter = 0
									content += strconv.Itoa(internalCounter)
									internalCounter++
								} else {
									content += strconv.Itoa(internalCounter)
									internalCounter++
								}
							}
						}

						inodetmp.I_s = int64(unsafe.Sizeof(Structs.DirectoriesBlocks{})) + int64(len(content))

						tam := len(content)
						var filesBlocks []string
						if tam > 10 {
							for tam > 10 {
								auxFiles := ""
								for m := 0; m < 10; m++ {
									auxFiles += string(content[m])
								}
								filesBlocks = append(filesBlocks, auxFiles)
								content = strings.ReplaceAll(content, auxFiles, "")
								tam = len(content)
							}
							if tam < 10 && tam != 0 {
								filesBlocks = append(filesBlocks, content)
							}
						} else {
							filesBlocks = append(filesBlocks, content)
						}

						if len(filesBlocks) > 16 {
							Error("MKFILE", "SE HA LLEGADO A LA MÁXIMA CAPACIDAD DE APUNDODORES DEL INODO PADRE", responseString)
							return
						}

						file.Close()

						file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
						if err != nil {
							Error("MKFILE", "No se ha encontrado el disco", responseString)
							return
						}

						for w := 0; w < len(filesBlocks); w++ {
							var control int
							control = int(bb) + w
							var fbAux Structs.FilesBlocks
							if inode.I_block[w] == -1 {
								file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(control), 0)
								var binAux bytes.Buffer
								binary.Write(&binAux, binary.BigEndian, fbAux)
								WritingBytes(file, binAux.Bytes())
							} else {
								fbAux = fileBlock
							}

							copy(fbAux.B_content[:], filesBlocks[w])
							file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(control), 0)
							var bin1 bytes.Buffer
							binary.Write(&bin1, binary.BigEndian, fbAux)
							WritingBytes(file, bin1.Bytes())
							counterFilesBlockse++
						}

						for w := 0; w < len(filesBlocks); w++ {
							inodetmp.I_block[w] = bb + int64(w)
						}
						inode.I_block[i] = past

						file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*father, 0)
						var binInode bytes.Buffer
						binary.Write(&binInode, binary.BigEndian, inode)
						WritingBytes(file, binInode.Bytes())

						fileWritten = true
						fatherSpace = false
						break
					} else {
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
						inodetmp.I_type = 1
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
							Error("MKFILE", "No se ha encontrado el disco", responseString)
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
	}

	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("MKFILE", "No se ha encontrado el disco", responseString)
		return
	}

	if fileWritten {
		if fatherSpace {
			file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*fatherFileFolder, 0)
			var binInode bytes.Buffer
			binary.Write(&binInode, binary.BigEndian, inode)
			WritingBytes(file, binInode.Bytes())

			file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*past+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*past, 0)
			var binFile bytes.Buffer
			binary.Write(&binFile, binary.BigEndian, folder)
			WritingBytes(file, binFile.Bytes())

			file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
			var binInodeFile bytes.Buffer
			binary.Write(&binInodeFile, binary.BigEndian, inodetmp)
			WritingBytes(file, binInodeFile.Bytes())

		} else {
			file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*past+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*past, 0)
			var binFile bytes.Buffer
			binary.Write(&binFile, binary.BigEndian, folder)
			WritingBytes(file, binFile.Bytes())

			file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
			var binInodeFile bytes.Buffer
			binary.Write(&binInodeFile, binary.BigEndian, inodetmp)
			WritingBytes(file, binInodeFile.Bytes())
		}

		if counterFilesBlockse > 0 {
			updateBm(spr, pth, "BI", responseString)
			for i := 0; i < counterFilesBlockse; i++ {
				updateBm(spr, pth, "BB", responseString)
			}
		} else {
			updateBm(spr, pth, "BI", responseString)
			updateBm(spr, pth, "BB", responseString)
		}

		ruta := ""
		for i := 0; i < len(copyPath); i++ {
			ruta += "/" + copyPath[i]
		}
		Message("MKFILE", "Se ha creado el archivo en la ruta "+ruta, responseString)
	} else {
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
		Message("MKFILE", "Se ha creado el directorio "+ruta, responseString)
		file.Close()
	}
}

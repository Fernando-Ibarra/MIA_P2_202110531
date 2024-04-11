package Comands

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"os"
	"strings"
	"unsafe"
)

func DataCat(context []string, partition Structs.Partition, pth string, responseString string) {
	var filesPath []string

	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		ntk := tk[0]
		newTk := ntk[:len(ntk)-1]
		if Compare(newTk, "file") {
			filesPath = append(filesPath, tk[1])
		}
	}

	if len(filesPath) == 0 {
		Error("CAT", "Se necesitan parametros obligatorios ", responseString)
		return
	}

	for i := 0; i < len(filesPath); i++ {
		tmp := GetPath(filesPath[i])
		content := cat(tmp, partition, pth, responseString)
		Message("CAT", "Contenido: "+content, responseString)
	}
}

func cat(path []string, partition Structs.Partition, pth string, responseString string) string {
	spr := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))
	if err != nil {
		Error("CAT", "No se ha encontrado el disco", responseString)
		return ""
	}

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &spr)
	if err_ != nil {
		Error("CAT", "Error al leer el archivo", responseString)
		return ""
	}

	file.Seek(spr.S_inode_start, 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("CAT", "Error al leer el archivo", responseString)
		return ""
	}

	var fileToFound string
	if len(path) == 0 {
		Error("CAT", "No se ha brindado un path valido", responseString)
		return ""
	}

	var bi int64
	fnd := false
	inodetmp := Structs.NewInodos()
	content := ""

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
						Error("CAT", "Error al leer el archivo", responseString)
						return ""
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
								Error("CAT", "Error al leer el archivo", responseString)
								return ""
							}

							if inode.I_uid != int64(Logged.Uid) {
								Error("CAT", "No tiene permisos para crear carpetas en este directorio", responseString)
								return ""
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
			Error("CAT", "No se pudo crear el directorio "+address+", no existen directorios", responseString)
			return ""
		}
	}

	fnd = false
	content = ""
	for i := 0; i < 16; i++ {
		if inode.I_block[i] != -1 {
			if i < 16 {
				// carpeta donde va tendria que estar el archivo
				file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
				data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
				buffer = bytes.NewBuffer(data)
				err_ = binary.Read(buffer, binary.BigEndian, &folder)
				if err_ != nil {
					Error("CAT", "Error al leer el archivo", responseString)
					return ""
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
							bi = folder.B_content[2].B_inodo
							file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
							data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
							buffer = bytes.NewBuffer(data)
							err_ = binary.Read(buffer, binary.BigEndian, &inodetmp)
							if err_ != nil {
								Error("CAT", "Error al leer el archivo", responseString)
								return ""
							}

							fileBlock := Structs.FilesBlocks{}
							for w := 0; w < 16; w++ {
								if inodetmp.I_block[w] != -1 {
									var control int64
									control = inodetmp.I_block[w]
									file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(control), 0)
									data = readBytes(file, int(unsafe.Sizeof(Structs.FilesBlocks{})))
									buffer = bytes.NewBuffer(data)
									err_ = binary.Read(buffer, binary.BigEndian, &fileBlock)
									if err_ != nil {
										Error("CAT", "Error al leer el archivo", responseString)
										return ""
									}
									for nam := 0; nam < len(fileBlock.B_content); nam++ {
										if fileBlock.B_content[nam] == 0 {
											continue
										}
										content += string(fileBlock.B_content[nam])
									}
								} else {
									break
								}
							}
							break
						} else if Compare(nameAux2, fileToFound) {
							bi = folder.B_content[3].B_inodo
							file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
							data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
							buffer = bytes.NewBuffer(data)
							err_ = binary.Read(buffer, binary.BigEndian, &inodetmp)
							if err_ != nil {
								Error("CAT", "Error al leer el archivo", responseString)
								return ""
							}

							fileBlock := Structs.FilesBlocks{}
							for w := 0; w < 16; w++ {
								if inodetmp.I_block[w] != -1 {
									var control int64
									control = inodetmp.I_block[w]
									file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*int64(control), 0)
									data = readBytes(file, int(unsafe.Sizeof(Structs.FilesBlocks{})))
									buffer = bytes.NewBuffer(data)
									err_ = binary.Read(buffer, binary.BigEndian, &fileBlock)
									if err_ != nil {
										Error("CAT", "Error al leer el archivo", responseString)
										return ""
									}
									for nam := 0; nam < len(fileBlock.B_content); nam++ {
										if fileBlock.B_content[nam] == 0 {
											continue
										}
										content += string(fileBlock.B_content[nam])
									}
								} else {
									break
								}
							}
							break
						}
					} else {
						break
					}
				}
			}
		}
	}
	return content
}

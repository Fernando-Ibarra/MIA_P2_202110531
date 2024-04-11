package Comands

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"os"
	"strings"
	"unsafe"
)

func DataRename(context []string, partition Structs.Partition, pth string, responseString string) {
	name := ""
	path := ""
	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "path") {
			path = tk[1]
		} else if Compare(tk[0], "name") {
			name = tk[1]
		}
	}
	if path == "" || name == "" {
		Error("RENAME", "Se necesitan parametros obligatorios para crear un directorio", responseString)
		return
	}
	tmp := GetPath(path)
	rename(tmp, partition, pth, name, responseString)
}

func rename(path []string, partition Structs.Partition, pth string, name string, responseString string) {
	copyPath := path
	spr := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))
	if err != nil {
		Error("RENAME", "No se ha encontrado el disco", responseString)
		return
	}

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &spr)
	if err_ != nil {
		Error("RENAME", "Error al leer el archivo", responseString)
		return
	}

	file.Seek(spr.S_inode_start, 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("RENAME", "Error al leer el archivo", responseString)
		return
	}

	var fileToFound string
	if len(path) == 0 {
		Error("RENAME", "No se ha brindado un path valido", responseString)
		return
	}

	fnd := false

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
						Error("RENAME", "Error al leer el archivo", responseString)
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
								Error("RENAME", "Error al leer el archivo", responseString)
								return
							}

							if inode.I_uid != int64(Logged.Uid) {
								Error("RENAME", "No tiene permisos para crear carpetas en este directorio", responseString)
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
			Error("RENAME", "No se pudo crear el directorio "+address+", no existen directorios", responseString)
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
					Error("RENAME", "Error al leer el archivo", responseString)
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

						copy(folder.B_content[2].B_name[:], name)

						file.Close()
						file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
						if err != nil {
							Error("RENAME", "No se ha encontrado el disco", responseString)
							return
						}

						file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
						var folderChange bytes.Buffer
						binary.Write(&folderChange, binary.BigEndian, folder)
						WritingBytes(file, folderChange.Bytes())

						file.Close()
						file, err = os.Open(strings.ReplaceAll(pth, "\"", ""))
						if err != nil {
							Error("RENAME", "No se ha encontrado el disco", responseString)
							return
						}

						ruta := ""
						for p := 0; p < len(copyPath)-1; p++ {
							ruta += "/" + copyPath[p]
						}
						ruta += "/" + name
						Message("RENAME", "Se ha actualiza la ruta: "+ruta, responseString)
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
						copy(folder.B_content[3].B_name[:], name)

						file.Close()
						file, err = os.OpenFile(strings.ReplaceAll(pth, "\"", ""), os.O_WRONLY, os.ModeAppend)
						if err != nil {
							Error("RENAME", "No se ha encontrado el disco", responseString)
							return
						}

						file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
						var folderChange bytes.Buffer
						binary.Write(&folderChange, binary.BigEndian, folder)
						WritingBytes(file, folderChange.Bytes())

						file.Close()
						file, err = os.Open(strings.ReplaceAll(pth, "\"", ""))
						if err != nil {
							Error("RENAME", "No se ha encontrado el disco", responseString)
							return
						}

						ruta := ""
						for p := 0; p < len(copyPath)-1; p++ {
							ruta += "/" + copyPath[p]
						}
						ruta += "/" + name
						Message("RENAME", "Se ha actualiza la ruta: "+ruta, responseString)
						break
					}
				}
			}
		}
	}
}

package Comands

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unsafe"
)

func MakeJson() string {

	var disksPath []string

	// current Path
	currentPath, _ := os.Getwd()
	diskPath := currentPath + "/MIA/P2/"

	// read files in the folder
	files, err := ioutil.ReadDir(diskPath)
	if err != nil {
		fmt.Println("Error al leer la carpeta - JSON TREE ")
		return ""
	}

	//  jsonResponse := "{\"disks\":
	jsonResponse := "["
	// clean disk arrange
	for _, disk := range files {
		if strings.Contains(disk.Name(), ".dsk") {
			currentPath, _ = os.Getwd()
			pathTmp := currentPath + "/MIA/P2/" + disk.Name()
			disksPath = append(disksPath, pathTmp)
			fmt.Println("Ruta de disco: ", disksPath)

			_, err2 := os.Open(strings.ReplaceAll(pathTmp, "\"", ""))
			if err2 != nil {
				fmt.Println("ERROR MJ 1")
				return ""
			}

			var partitions [4]Structs.Partition
			var logicPartitions []Structs.EBR

			mbr := readDisk(pathTmp, "")
			partitions[0] = mbr.Mbr_partitions_1
			partitions[1] = mbr.Mbr_partitions_2
			partitions[2] = mbr.Mbr_partitions_3
			partitions[3] = mbr.Mbr_partitions_4

			//
			jsonResponse += "{"
			jsonResponse += "\"name\": " + "\"" + disk.Name() + "\","
			jsonResponse += "\"partitions\":["
			for i := 0; i < len(partitions); i++ {
				if partitions[i].Part_type == 'E' {
					logicPartitions = GetLogics(partitions[i], pathTmp, "")
					for k := 0; k < len(logicPartitions); k++ {
						logicPartitionName := ""
						for m := 0; m < len(logicPartitions[k].Part_name); m++ {
							if logicPartitions[k].Part_name[m] != 0 {
								logicPartitionName += string(logicPartitions[k].Part_name[m])
							}
						}
						jsonResponse += "{"
						jsonResponse += "\"name\":" + "\" " + logicPartitionName + "\","
						jsonResponse += "\"users\":" + "["
						users := getUsers(pathTmp, partitions[i])
						jsonResponse += users
						jsonResponse += "],"
						jsonResponse += "\"fileSystem\":" + "["
						response := getFileSystem(pathTmp, partitions[i])
						jsonResponse += response
						jsonResponse += "]"
						jsonResponse += "},"
					}
				} else {
					partitionName := ""
					for j := 0; j < len(partitions[i].Part_name); j++ {
						if partitions[i].Part_name[j] != 0 {
							partitionName += string(partitions[i].Part_name[j])
						}
					}
					if len(partitionName) > 0 {
						jsonResponse += "{"
						jsonResponse += "\"name\":" + "\" " + partitionName + "\","
						jsonResponse += "\"users\":" + "["
						users := getUsers(pathTmp, partitions[i])
						jsonResponse += users
						jsonResponse += "],"
						jsonResponse += "\"fileSystem\":" + "["
						response := getFileSystem(pathTmp, partitions[i])
						jsonResponse += response
						jsonResponse += "]"
						jsonResponse += "},"
					}
				}
			}
			jsonResponse += "]"
			jsonResponse += "},"

		}
	}
	jsonResponse += "]"
	// jsonResponse += "]

	// number of disks
	// disks := len(disksPath)
	fmt.Println(jsonResponse)
	return jsonResponse
}

func getFileSystem(path string, partition Structs.Partition) string {
	fileSystem := ""
	super := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		fmt.Println("ERROR - 1")
		return ""
	}

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		fmt.Println("ERROR - 2")
		return ""
	}

	if int(super.S_block_start) == 0 {
		return ""
	}

	file.Seek(super.S_inode_start, 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		fmt.Println("ERROR - 3")
		return ""
	}
	for i := 0; i < 16; i++ {
		if inode.I_block[i] != -1 {
			folder = Structs.NewDirectoriesBlocks()
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
			data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &folder)

			if folder.B_content[2].B_inodo != -1 {
				name := ""
				for nam := 0; nam < len(folder.B_content[2].B_name); nam++ {
					if folder.B_content[2].B_name[nam] == 0 {
						continue
					}
					name += string(folder.B_content[2].B_name[nam])
				}
				fileSystem += "{\"file\":"
				fileSystem += "{\"name\":" + "\"" + name + "\","
				fileSystem += "\"content\":" + recursiveFileSystem(folder.B_content[2].B_inodo, partition, path)
				fileSystem += "},"
				fileSystem += "},"
			}

			if folder.B_content[3].B_inodo != -1 {
				name := ""
				for nam := 0; nam < len(folder.B_content[3].B_name); nam++ {
					if folder.B_content[3].B_name[nam] == 0 {
						continue
					}
					name += string(folder.B_content[3].B_name[nam])
				}
				fileSystem += "{\"folder\":"
				fileSystem += "{\"name\":" + "\"" + name + "\","

				fileSystem += "\"content\":" + recursiveFileSystem(folder.B_content[3].B_inodo, partition, path)
				// fileSystem += recursiveFileSystem(folder.B_content[3].B_inodo, partition, path)
				fileSystem += "}"
				fileSystem += "},"
			}
		} else {
			break
		}
	}
	fmt.Println(fileSystem)
	return fileSystem
}

func recursiveFileSystem(control int64, partition Structs.Partition, path string) string {
	fileSystem := ""
	super := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		fmt.Println("ERROR PATH" + path)
		fmt.Println("ERROR - 4")
		fmt.Println(err)
		return ""
	}

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		fmt.Println("ERROR - 5")
		return ""
	}

	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*control, 0)
	data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		fmt.Println("ERROR - 5")
		return ""
	}

	content := ""
	if inode.I_type == 0 {
		fileSystem += "["
		for i := 0; i < 16; i++ {
			if i < 16 {
				if inode.I_block[i] != -1 {
					folder = Structs.NewDirectoriesBlocks()
					file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)
					data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
					buffer = bytes.NewBuffer(data)
					err_ = binary.Read(buffer, binary.BigEndian, &folder)
					// fileSystem += "["

					if folder.B_content[2].B_inodo != -1 {
						name := ""
						for nam := 0; nam < len(folder.B_content[2].B_name); nam++ {
							if folder.B_content[2].B_name[nam] == 0 {
								continue
							}
							name += string(folder.B_content[2].B_name[nam])
						}

						// fileSystem += "{\"" + name + "\":" + "["
						// response := recursiveFileSystem(folder.B_content[2].B_inodo, partition, path)
						/* if response != "" {
							fileSystem += response
						} else {
							continue
						} */
						// fileSystem += "]},"

						fileSystem += "{\"folder\":"
						fileSystem += "{\"name\":" + "\"" + name + "\","
						// fileSystem += "\"content\":" + recursiveFileSystem(folder.B_content[3].B_inodo, partition, path)

						response := recursiveFileSystem(folder.B_content[2].B_inodo, partition, path)
						if response != "" {
							fileSystem += "\"content\":" + response
						} else {
							continue
						}

						// fileSystem += recursiveFileSystem(folder.B_content[3].B_inodo, partition, path)
						fileSystem += "}"
						fileSystem += "},"
					}

					if folder.B_content[3].B_inodo != -1 {
						name := ""
						for nam := 0; nam < len(folder.B_content[3].B_name); nam++ {
							if folder.B_content[3].B_name[nam] == 0 {
								continue
							}
							name += string(folder.B_content[3].B_name[nam])
						}
						// fileSystem += "{\"" + name + "\":" + "["
						// response := recursiveFileSystem(folder.B_content[3].B_inodo, partition, path)
						/* if response != "" {
							fileSystem += response
						} else {
							continue
						}*/
						// fileSystem += "]},"

						fileSystem += "{\"folder\":"
						fileSystem += "{\"name\":" + "\"" + name + "\","
						// fileSystem += "\"content\":" + recursiveFileSystem(folder.B_content[3].B_inodo, partition, path)

						response := recursiveFileSystem(folder.B_content[3].B_inodo, partition, path)
						if response != "" {
							fileSystem += "\"content\":" + response
						} else {
							continue
						}

						// fileSystem += recursiveFileSystem(folder.B_content[3].B_inodo, partition, path)
						fileSystem += "}"
						fileSystem += "},"
					}

					// fileSystem += "],"

				} else {
					break
				}
			} else {
				break
			}
		}
		fileSystem += "],"
	} else if inode.I_type == 1 {
		for i := 0; i < 16; i++ {
			if inode.I_block[i] != -1 {
				fileB := Structs.FilesBlocks{}
				file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*inode.I_block[i], 0)
				data = readBytes(file, int(unsafe.Sizeof(Structs.FilesBlocks{})))
				buffer = bytes.NewBuffer(data)
				err_ = binary.Read(buffer, binary.BigEndian, &fileB)

				contentFile := ""
				for nam := 0; nam < len(fileB.B_content); nam++ {
					if fileB.B_content[nam] == 0 {
						continue
					}
					contentFile += string(fileB.B_content[nam])
				}
				contentFile = strings.ReplaceAll(contentFile, ",", ".")
				contentFile = strings.ReplaceAll(contentFile, " ", "")
				content += strings.ReplaceAll(contentFile, "\n", "")
			} else {
				break
			}
		}
		fileSystem += "\"" + content + "\""
	}

	return fileSystem
}

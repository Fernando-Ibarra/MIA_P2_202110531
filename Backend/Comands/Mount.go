package Comands

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

var DiskMount [99]DiskMounted

type DiskMounted struct {
	Path       [150]byte
	State      byte
	Partitions [26]PartitionMounted
}

type PartitionMounted struct {
	Letter byte
	State  byte
	Name   [20]byte
}

func DataMount(tokens []string, responseString string) {
	driveLetter := ""
	name := ""
	letter := ""
	for i := 0; i < len(tokens); i++ {
		current := tokens[i]
		command := strings.Split(current, "=")
		if Compare(command[0], "name") {
			name = command[1]
		} else if Compare(command[0], "driveletter") {
			currentPath, _ := os.Getwd()
			letter = command[1]
			driveLetter = currentPath + "/MIA/P2/" + command[1] + ".dsk"
		}
	}
	if driveLetter == "" || name == "" {
		Error("MOUNT", "El comando MOUNT requiere parámetros obligatorios", responseString)
		return
	}
	mount(driveLetter, name, letter, responseString)
	listMount()
}

func mount(d string, n string, l string, responseString string) {
	file, error_ := os.Open(d)
	if error_ != nil {
		Error("MOUNT", "No se ha podido abrir el archivo", responseString)
		return
	}
	disk := Structs.NewMBR()
	file.Seek(0, 0)

	data := readBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &disk)
	if err_ != nil {
		Error("MOUNT", "Error al leer el archivo", responseString)
		return
	}
	err := file.Close()
	if err != nil {
		return
	}

	partition := SearchPartitions(disk, n, d, responseString)
	if partition == nil {
		Error("MOUNT", "No se encontró la partición "+n, responseString)
		return
	}
	if partition.Part_type == 'E' || partition.Part_type == 'L' {
		var name [16]byte
		copy(name[:], n)
		if partition.Part_name == name && partition.Part_type == 'E' {
			Error("MOUNT", "No se puede montar una partición extendida", responseString)
			return
		} else {
			ebrs := GetLogics(*partition, d, responseString)
			founded := false
			if len(ebrs) != 0 {
				for i := 0; i < len(ebrs); i++ {
					ebr := ebrs[i]
					nameEbr := ""
					for j := 0; j < len(ebr.Part_name); j++ {
						if ebr.Part_name[j] != 0 {
							nameEbr += string(ebr.Part_name[j])
						}
					}
					if Compare(nameEbr, n) && ebr.Part_mount == '1' {
						founded = true
						n = nameEbr
						break
					} else if nameEbr == n && ebr.Part_mount == '0' {
						Error("MOUNT", "No se puede montar una partición lógica eliminada", responseString)
						return
					}
				}
				if !founded {
					Error("MOUNT", "No se encontró la partición lógica", responseString)
					return
				}
			}
		}
	}
	for i := 0; i < 99; i++ {
		var path [150]byte
		copy(path[:], d)
		if DiskMount[i].Path == path {
			for j := 0; j < 26; j++ {
				var name [20]byte
				copy(name[:], n)
				if DiskMount[i].Partitions[j].Name == name {
					Error("MOUNT", "Ya se ha montado la partición "+n, responseString)
					return
				}
				if DiskMount[i].Partitions[j].State == 0 {
					DiskMount[i].Partitions[j].State = 1
					DiskMount[i].Partitions[j].Letter = l[0]
					copy(DiskMount[i].Partitions[j].Name[:], n)
					res := l + strconv.Itoa(i+1) + strconv.Itoa(31)
					Message("MOUNT", "Se ha realizado correctamente el mount -id="+res, responseString)
					return
				}
			}
		}
	}
	for i := 0; i < 99; i++ {
		if DiskMount[i].State == 0 {
			DiskMount[i].State = 1
			copy(DiskMount[i].Path[:], d)
			for j := 0; j < 26; j++ {
				if DiskMount[i].Partitions[j].State == 0 {
					DiskMount[i].Partitions[j].State = 1
					DiskMount[i].Partitions[j].Letter = l[0]
					copy(DiskMount[i].Partitions[j].Name[:], n)
					res := l + strconv.Itoa(i+1) + strconv.Itoa(31)
					Message("MOUNT", "Se ha realizado correctamente el mount -id="+res, responseString)
					return
				}
			}
		}
	}
}

func GetMount(comand string, id string, p *string, responseString string) Structs.Partition {
	if !(id[2] == '3' && id[3] == '1') {
		Error(comand, "El primer identificador no es válido", responseString)
		return Structs.Partition{}
	}
	letter := id[0]
	j, _ := strconv.Atoi(string(id[1] - 1))
	if j < 0 {
		Error(comand, "El primer identificador no es válido", responseString)
		return Structs.Partition{}
	}
	for i := 0; i < 99; i++ {
		if DiskMount[i].Partitions[j].State == 1 {
			if DiskMount[i].Partitions[j].Letter == letter {
				path := ""
				for k := 0; k < len(DiskMount[i].Path); k++ {
					if DiskMount[i].Path[k] != 0 {
						path += string(DiskMount[i].Path[k])
					}
				}
				file, erro := os.Open(strings.ReplaceAll(path, "\"", ""))
				if erro != nil {
					Error(comand, "No se encontro el disco", responseString)
					return Structs.Partition{}
				}
				disk := Structs.NewMBR()
				file.Seek(0, 0)
				data := readBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
				buffer := bytes.NewBuffer(data)
				err_ := binary.Read(buffer, binary.BigEndian, &disk)
				if err_ != nil {
					Error("FDISK", "Error al leer el archivo", responseString)
					return Structs.Partition{}
				}
				file.Close()

				partitionName := ""
				for k := 0; k < len(DiskMount[i].Partitions[j].Name); k++ {
					if DiskMount[i].Partitions[j].Name[k] != 0 {
						partitionName += string(DiskMount[i].Partitions[j].Name[k])
					}
				}
				*p = path
				return *SearchPartitions(disk, partitionName, path, responseString)
			}
		}
	}
	return Structs.Partition{}
}

func listMount() {
	fmt.Println("\n------------------------------ LISTADO DE MOUNTS ------------------------------")
	for i := 0; i < 99; i++ {
		for j := 0; j < 26; j++ {
			if DiskMount[i].Partitions[j].State == 1 {
				name := ""
				for k := 0; k < len(DiskMount[i].Partitions[j].Name); k++ {
					if DiskMount[i].Partitions[j].Name[k] != 0 {
						name += string(DiskMount[i].Partitions[j].Name[k])
					}
				}
				fmt.Println("\t id=" + string(DiskMount[i].Partitions[j].Letter) + strconv.Itoa(j+1) + "31, Nombre: " + name)
			}
		}
	}
}

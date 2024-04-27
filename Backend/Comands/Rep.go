package Comands

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

func DataRep(context []string, responseString string) {
	name := ""
	pathOut := ""
	id := ""
	ruta := ""
	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "id") {
			id = tk[1]
		} else if Compare(tk[0], "name") {
			name = tk[1]
		} else if Compare(tk[0], "path") {
			pathOut = tk[1]
		} else if Compare(tk[0], "ruta") {
			ruta = tk[1]
		}
	}

	if id == "" || pathOut == "" || name == "" {
		Error("REP", "Se necesitan parámetros obligatorios para el comando rep", responseString)
		return
	}

	if Compare(name, "file") && ruta == "" {
		Error("REP", "Se necesitan parámetros obligatorios para el comando rep de tipo file", responseString)
		return
	}

	aux := strings.Split(pathOut, "/")

	last := len(aux)

	nameDot := aux[last-1]

	nameG := nameDot[:len(nameDot)-3]

	currentPath, _ := os.Getwd()
	diskPath := currentPath + "/MIA/P2/Rep/"
	pathOut = diskPath + nameG

	fmt.Println(pathOut)

	if Compare(name, "mbr") {
		repMBR(id, pathOut, responseString)
	} else if Compare(name, "sb") {
		repSuperBlock(id, pathOut, responseString)
	} else if Compare(name, "disk") {
		repDisk(id, pathOut, responseString)
	} else if Compare(name, "bm_inode") {
		repBM(id, pathOut, "BI", responseString)
	} else if Compare(name, "bm_bloc") {
		repBM(id, pathOut, "BB", responseString)
	} else if Compare(name, "inode") {
		repInode(id, pathOut, responseString)
	} else if Compare(name, "block") {
		repBlock(id, pathOut, responseString)
	} else if Compare(name, "tree") {
		repTree(id, pathOut, responseString)
	} else if Compare(name, "journaling") {
		repJournaling(id, pathOut, responseString)
	} else if Compare(name, "file") {
		repCat(id, pathOut, ruta, responseString)
	} else if Compare(name, "ls") {
		repLs(id, pathOut, responseString)
	}
}

func repMBR(id string, pathOut string, responseString string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido", responseString)
		return
	}
	letter := id[0]
	currentPath, _ := os.Getwd()
	driveLetter := currentPath + "/MIA/P2/" + string(letter) + ".dsk"

	aux := strings.Split(pathOut, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos", responseString)
		return
	}
	pd := aux[0] + ".dot"

	var partitions [4]Structs.Partition
	var logicPartitions []Structs.EBR

	_, err := os.OpenFile(strings.ReplaceAll(driveLetter, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("FDISK", "Error al abrir el archivo", responseString)
		return
	}

	mbr := readDisk(driveLetter, responseString)
	partitions[0] = mbr.Mbr_partitions_1
	partitions[1] = mbr.Mbr_partitions_2
	partitions[2] = mbr.Mbr_partitions_3
	partitions[3] = mbr.Mbr_partitions_4

	text := "digraph MBR{\n"
	text += "node [ shape=none fontname=Arial ]\n"
	text += "n1 [ label = <\n"
	text += "<table>\n"
	text += "<tr><td colspan=\"2\" bgcolor=\"blueviolet\"><font color=\"white\">REPORTE DE MBR</font></td></tr>\n"
	text += "<tr><td bgcolor=\"white\">mbr_tamano</td><td bgcolor=\"white\">" + strconv.Itoa(int(mbr.Mbr_tamano)) + "</td></tr>\n"
	fechaC := ""
	for i := 0; i < len(mbr.Mbr_fecha_creacion); i++ {
		if mbr.Mbr_fecha_creacion[i] != 0 {
			fechaC += string(mbr.Mbr_fecha_creacion[i])
		}
	}
	text += "<tr><td bgcolor=\"thistle\">mbr_fecha_creacion</td><td bgcolor=\"thistle\">" + fechaC + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">mbr_dsk_signature</td><td bgcolor=\"white\">" + strconv.Itoa(int(mbr.Mbr_dsk_signature)) + "</td></tr>\n"
	for i := 0; i < len(partitions); i++ {
		if partitions[i].Part_type == 'E' {
			text += "<tr><td colspan=\"2\" bgcolor=\"blueviolet\"><font color=\"white\">Particion Extendida</font></td></tr>\n"
			text += "<tr><td bgcolor=\"white\">part_status</td><td bgcolor=\"white\">" + string(partitions[i].Part_status) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"thistle\">part_type</td><td bgcolor=\"thistle\">" + string(partitions[i].Part_type) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"white\">part_fit</td><td bgcolor=\"white\">" + string(partitions[i].Part_fit) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"thistle\">part_start</td><td bgcolor=\"thistle\">" + strconv.Itoa(int(partitions[i].Part_start)) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"white\">part_s</td><td bgcolor=\"white\">" + strconv.Itoa(int(partitions[i].Part_s)) + "</td></tr>\n"
			partitionName := ""
			for j := 0; j < len(partitions[i].Part_name); j++ {
				if partitions[i].Part_name[j] != 0 {
					partitionName += string(partitions[i].Part_name[j])
				}
			}
			text += "<tr><td bgcolor=\"thistle\">part_name</td><td bgcolor=\"thistle\">" + partitionName + "</td></tr>\n"

			logicPartitions = GetLogics(partitions[i], driveLetter, responseString)
			for k := 0; k < len(logicPartitions); k++ {
				text += "<tr><td colspan=\"2\" bgcolor=\"salmon\"><font color=\"white\">Particion Lógica - EBR</font></td></tr>\n"
				text += "<tr><td bgcolor=\"white\">part_mount</td><td bgcolor=\"white\">" + string(logicPartitions[k].Part_mount) + "</td></tr>\n"
				text += "<tr><td bgcolor=\"lightsalmon\">part_fit</td><td bgcolor=\"lightsalmon\">" + string(logicPartitions[k].Part_fit) + "</td></tr>\n"
				text += "<tr><td bgcolor=\"white\">part_start</td><td bgcolor=\"white\">" + strconv.Itoa(int(logicPartitions[k].Part_start)) + "</td></tr>\n"
				text += "<tr><td bgcolor=\"lightsalmon\">part_s</td><td bgcolor=\"lightsalmon\">" + strconv.Itoa(int(logicPartitions[k].Part_s)) + "</td></tr>\n"
				text += "<tr><td bgcolor=\"white\">part_next</td><td bgcolor=\"white\">" + strconv.Itoa(int(logicPartitions[k].Part_next)) + "</td></tr>\n"
				logicPartitionName := ""
				for m := 0; m < len(logicPartitions[k].Part_name); m++ {
					if logicPartitions[k].Part_name[m] != 0 {
						logicPartitionName += string(logicPartitions[k].Part_name[m])
					}
				}
				text += "<tr><td bgcolor=\"lightsalmon\">part_name</td><td bgcolor=\"lightsalmon\">" + logicPartitionName + "</td></tr>\n"
			}

		} else {
			text += "<tr><td colspan=\"2\" bgcolor=\"blueviolet\"><font color=\"white\">Particion Primaria</font></td></tr>\n"
			text += "<tr><td bgcolor=\"white\">part_status</td><td bgcolor=\"white\">" + string(partitions[i].Part_status) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"thistle\">part_type</td><td bgcolor=\"thistle\">" + string(partitions[i].Part_type) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"white\">part_fit</td><td bgcolor=\"white\">" + string(partitions[i].Part_fit) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"thistle\">part_start</td><td bgcolor=\"thistle\">" + strconv.Itoa(int(partitions[i].Part_start)) + "</td></tr>\n"
			text += "<tr><td bgcolor=\"white\">part_s</td><td bgcolor=\"white\">" + strconv.Itoa(int(partitions[i].Part_s)) + "</td></tr>\n"
			partitionName := ""
			for j := 0; j < len(partitions[i].Part_name); j++ {
				if partitions[i].Part_name[j] != 0 {
					partitionName += string(partitions[i].Part_name[j])
				}
			}
			if len(partitionName) > 0 {
				text += "<tr><td bgcolor=\"thistle\">part_name</td><td bgcolor=\"thistle\">" + partitionName + "</td></tr>\n"
			} else {
				text += "<tr><td bgcolor=\"thistle\">part_name</td><td bgcolor=\"thistle\"></td></tr>\n"
			}

		}
	}
	text += "</table>\n"
	text += "> ]\n"
	text += "}\n"

	CreateFile(pd)
	WriteFile(text, pd)
	// termination := strings.Split(pathOut, ".")
	// Execute(pathOut, pd, termination[1], responseString)
	Message("REP", "Reporte de MBR se ha generado correctamente en"+pathOut, responseString)
}

func repSuperBlock(id string, pathOut string, responseString string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido", responseString)
		return
	}

	var path string
	partition := GetMount("MKGRP", id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("MKGRP", "No se encontró la partición montada con el id: "+id, responseString)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKGRP", "No se ha encontrado el disco", responseString)
		return
	}

	aux := strings.Split(path, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos", responseString)
		return
	}
	pd := aux[0] + ".dot"

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("MKGRP", "Error al leer el archivo", responseString)
		return
	}

	text := "digraph SuperBloque{\n"
	text += "node [ shape=none fontname=Arial ]\n"
	text += "n1 [ label = <\n"
	text += "<table>\n"
	text += "<tr><td colspan=\"2\" bgcolor=\"palegreen4\"><font color=\"white\">REPORTE DE SuperBloque</font></td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_filesystem_type</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_filesystem_type)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_inodes_count</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_inodes_count)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_blocks_count</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_blocks_count)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_free_inodes_count</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_free_inodes_count)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_free_blocks_count</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_free_blocks_count)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_mtime</td><td bgcolor=\"palegreen2\">" + string(super.S_mtime[:]) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_umtime</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_umtime)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_mnt_count</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_mnt_count)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_magic</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_magic)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_inode_s</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_inode_s)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_block_s</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_block_s)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_firts_ino</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_firts_ino)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_firts_blo</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_firts_blo)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_bm_inode_start</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_bm_inode_start)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_bm_block_start</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_bm_block_start)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"palegreen2\">s_inode_start</td><td bgcolor=\"palegreen2\">" + strconv.Itoa(int(super.S_inode_start)) + "</td></tr>\n"
	text += "<tr><td bgcolor=\"white\">s_block_start</td><td bgcolor=\"white\">" + strconv.Itoa(int(super.S_block_start)) + "</td></tr>\n"
	text += "</table>\n"
	text += "> ]\n"
	text += "}\n"

	file.Close()

	CreateFile(pd)
	WriteFile(text, pd)
	// termination := strings.Split(pathOut, ".")
	// Execute(pathOut, pd, termination[1], responseString)
	Message("REP", "Reporte de Superbloque se ha generado correctamente en"+pathOut, responseString)
}

func repDisk(id string, pathOut string, responseString string) {
	var path string
	GetMount("REP", id, &path, responseString)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco", responseString)
		return
	}

	var disk Structs.MBR
	file.Seek(0, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &disk)
	if err_ != nil {
		Error("REP", "Error al leer el archivo", responseString)
		return
	}
	file.Close()

	aux := strings.Split(pathOut, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos", responseString)
		return
	}
	pd := aux[0] + ".dot"

	folder := ""
	address := strings.Split(pd, "/")

	fileaux, _ := os.Open(strings.ReplaceAll(pd, "\"", ""))
	if fileaux == nil {
		for i := 0; i < len(address); i++ {
			folder += "/" + address[i]
			if _, err_2 := os.Stat(folder); os.IsNotExist(err_2) {
				os.Mkdir(folder, 0777)
			}
		}
		os.Remove(pd)
	} else {
		fileaux.Close()
	}

	partitions := GetPartitions(disk)
	var extended Structs.Partition
	ext := false
	for i := 0; i < 4; i++ {
		if partitions[i].Part_status == '1' {
			if partitions[i].Part_type == "E"[0] || partitions[i].Part_type == "e"[0] {
				ext = true
				extended = partitions[i]
			}
		}
	}

	content := "digraph Disk{\n"
	content += "rankdir=TB;\n"
	content += "forcelabels=true;\n"
	content += "graph [dpi = \"600\"];\n"
	content += "node [ shape=plaintext fontname=Arial ]\n"
	content += "n1 [ label = <\n"
	content += "<table>\n"
	content += "<tr>\n"

	var positions [5]int64
	var positionsii [5]int64

	positions[0] = disk.Mbr_partitions_1.Part_start - (1 + int64(unsafe.Sizeof(Structs.MBR{})))
	positions[1] = disk.Mbr_partitions_2.Part_start - disk.Mbr_partitions_1.Part_start + disk.Mbr_partitions_1.Part_s
	positions[2] = disk.Mbr_partitions_3.Part_start - disk.Mbr_partitions_2.Part_start + disk.Mbr_partitions_2.Part_s
	positions[3] = disk.Mbr_partitions_4.Part_start - disk.Mbr_partitions_3.Part_start + disk.Mbr_partitions_3.Part_s
	positions[4] = disk.Mbr_tamano + 1 - disk.Mbr_partitions_4.Part_start + disk.Mbr_partitions_4.Part_s

	copy(positionsii[:], positions[:])

	logic := 0
	tmpLogic := ""

	if ext {
		tmpLogic += "<tr>\n"
		auxEBR := Structs.NewEBR()
		file, err = os.Open(strings.ReplaceAll(path, "\"", ""))

		if err != nil {
			Error("REP", "No se ha encontrado el disco", responseString)
			return
		}

		file.Seek(extended.Part_start, 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &auxEBR)
		if err_ != nil {
			Error("REP", "Error al leer el archivo", responseString)
			return
		}
		file.Close()

		var tamGen int64 = 0
		for auxEBR.Part_next != -1 {
			tamGen += auxEBR.Part_s
			res := float64(auxEBR.Part_s) / float64(disk.Mbr_tamano)
			res = res * 100
			tmpLogic += "<td>\"EBR\"</td>"
			s := fmt.Sprintf("%.2f", res)
			tmpLogic += "<td>\"Lógica \n " + s + "% de la partición extendida</td>\n"

			resta := float64(auxEBR.Part_next) - (float64(auxEBR.Part_start) + float64(auxEBR.Part_s))
			resta = resta / float64(disk.Mbr_tamano)
			resta = resta * 10000.00
			resta = math.Round(resta) / 100.00
			if resta != 0 {
				s = fmt.Sprintf("%f", resta)
				tmpLogic += "<td>\"Lógica\n " + s + "% libre de la partición extendida</td>\n"
				logic++
			}
			logic += 2
			file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
			if err != nil {
				Error("REP", "No se ha encontrado el disco", responseString)
				return
			}

			file.Seek(auxEBR.Part_next, 0)
			data = readBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &auxEBR)
			if err_ != nil {
				Error("REP", "Error al leer el archivo", responseString)
				return
			}
			file.Close()
		}
		resta := float64(extended.Part_s) - float64(tamGen)
		resta = resta / float64(disk.Mbr_tamano)
		resta = math.Round(resta * 100)
		if resta != 0 {
			s := fmt.Sprintf("%.2f", resta)
			tmpLogic += "<td>\"Libre \n " + s + "% de la partición extendida \"</td>\n"
			logic++
		}
		tmpLogic += "</tr>\n"
		logic += 2
	}
	var tamPrim int64
	for i := 0; i < 4; i++ {
		if partitions[i].Part_type == 'E' {
			tamPrim += partitions[i].Part_s
			res := float64(partitions[i].Part_s) / float64(disk.Mbr_tamano)
			res = math.Round(res*10000.00) / 100.00
			s := fmt.Sprintf("%.3f", res)
			content += "<td COLSPAN='" + strconv.Itoa(logic) + "'>Extendida \n" + s + "% del disco</td>\n"
		} else if partitions[i].Part_start != -1 {
			tamPrim += partitions[i].Part_s
			res := float64(partitions[i].Part_s) / float64(disk.Mbr_tamano)
			res = math.Round(res*10000.00) / 100.00
			s := fmt.Sprintf("%.3f", res)
			content += "<td ROWSPAN='2'>Primaria \n" + s + "% del disco</td>\n"
		}
	}

	if tamPrim != 0 {
		libre := disk.Mbr_tamano - tamPrim
		res := float64(libre) / float64(disk.Mbr_tamano)
		res = math.Round(res * 100)
		s := fmt.Sprintf("%.3f", res)
		content += "<td ROWSPAN='2'>Libre\n" + s + "% del disco</td>"
	}

	content += "</tr>\n"
	content += tmpLogic
	content += "</table>\n"
	content += "> ]\n"
	content += "}\n"

	CreateFile(pd)
	WriteFile(content, pd)
	// termination := strings.Split(pathOut, ".")
	// Execute(pathOut, pd, termination[1], responseString)
	Message("REP", "Reporte del Disco se ha generado correctamente en"+pathOut, responseString)
}

func repBM(id string, pathOut, t string, responseString string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido", responseString)
		return
	}

	var path string
	partition := GetMount("MKGRP", id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("REP", "No se encontró la partición montada con el id: "+id, responseString)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco", responseString)
		return
	}

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("REP", "Error al leer el archivo", responseString)
		return
	}

	CreateFile(pathOut + "txt")
	content := ""
	counter := 1
	ch := '2'
	if t == "BI" {
		file.Seek(super.S_bm_inode_start, 0)
		for i := 0; i < int(super.S_inodes_count); i++ {
			data = readBytes(file, int(unsafe.Sizeof(ch)))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &ch)
			if err_ != nil {
				Error("REP", "Error al leer el archivo", responseString)
				return
			}

			element := fromBMtoFile(strconv.Itoa(int(ch)))
			if element == "-1" {
				break
			}
			if counter == 20 {
				content += element + "\n"
				counter = 1
			} else {
				content += element
				counter++
			}
		}
	} else {
		file.Seek(super.S_bm_block_start, 0)
		for i := 0; i < int(super.S_inodes_count); i++ {
			data = readBytes(file, int(unsafe.Sizeof(ch)))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &ch)
			if err_ != nil {
				Error("REP", "Error al leer el archivo", responseString)
				return
			}
			element := fromBMtoFile(strconv.Itoa(int(ch)))
			if element == "-1" {
				break
			}
			if counter == 20 {
				content += element + "\n"
				counter = 1
			} else {
				content += element
				counter++
			}
		}
	}

	WriteFile(content, pathOut+"txt")
	if t == "BI" {
		Message("REP", "Reporte de los bitmaps de Inodos "+pathOut+", creado correctamente", responseString)
	} else {
		Message("REP", "Reporte de los bitmaps de bloques "+pathOut+", creado correctamente", responseString)
	}
}

func fromBMtoFile(ch string) string {
	if ch == "48" {
		return "0"
	} else if ch == "49" {
		return "1"
	} else {
		return "-1"
	}
	return "-1"
}

func repInode(id string, pathOut string, responseString string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido", responseString)
		return
	}

	var path string
	partition := GetMount("MKGRP", id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("REP", "No se encontró la partición montada con el id: "+Logged.Id, responseString)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco", responseString)
		return
	}

	aux := strings.Split(pathOut, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos", responseString)
		return
	}
	pd := aux[0] + ".dot"
	fmt.Println(pd)

	super := Structs.NewSuperBlock()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("REP", "Error al leer el archivo", responseString)
		return
	}

	content := "digraph Inodos{\n"
	content += "node [ shape=plaintext fontname=Arial ]\n"

	var inodes []Structs.Inodos
	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start, 0)
	for i := 0; i < int(super.S_inodes_count); i++ {
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("REP", "Error al leer el archivo", responseString)
			return
		}

		if inode.I_uid == -1 {
			break
		}
		inodes = append(inodes, inode)
	}

	for i := 0; i < len(inodes); i++ {
		content += "A" + strconv.Itoa(i)
		content += "[label= <"
		content += "<table border=\"1\" cellborder=\"0\">\n"
		content += "<tr><td bgcolor=\"dodgerblue4\" colspan=\"2\" ><font color=\"white\">Inodo " + strconv.Itoa(i) + "</font></td></tr>\n"
		content += "<tr><td bgcolor=\"deepskyblue\">I_uid</td><td bgcolor=\"deepskyblue\">" + strconv.Itoa(int(inodes[i].I_uid)) + "</td></tr>\n"
		content += "<tr><td bgcolor=\"deepskyblue\">I_gid</td><td bgcolor=\"deepskyblue\">" + strconv.Itoa(int(inodes[i].I_gid)) + "</td></tr>\n"
		content += "<tr><td bgcolor=\"deepskyblue\">I_s</td><td bgcolor=\"deepskyblue\">" + strconv.Itoa(int(inodes[i].I_s)) + "</td></tr>\n"
		atime := ""
		for k := 0; k < len(inodes[i].I_atime); k++ {
			if inodes[i].I_atime[k] != 0 {
				atime += string(inodes[i].I_atime[k])
			}
		}
		content += "<tr><td bgcolor=\"deepskyblue\">I_atime</td><td bgcolor=\"deepskyblue\">" + atime + "</td></tr>\n"
		ctime := ""
		for k := 0; k < len(inodes[i].I_ctime); k++ {
			if inodes[i].I_ctime[k] != 0 {
				ctime += string(inodes[i].I_ctime[k])
			}
		}
		content += "<tr><td bgcolor=\"deepskyblue\">I_ctime</td><td bgcolor=\"deepskyblue\">" + ctime + "</td></tr>\n"
		mtime := ""
		for k := 0; k < len(inodes[i].I_mtime); k++ {
			if inodes[i].I_mtime[k] != 0 {
				mtime += string(inodes[i].I_mtime[k])
			}
		}
		content += "<tr><td bgcolor=\"deepskyblue\">I_mtime</td><td bgcolor=\"deepskyblue\">" + mtime + "</td></tr>\n"
		for j := 0; j < len(inodes[i].I_block); j++ {
			if j > 12 {
				content += "<tr><td bgcolor=\"azure3\">I_block " + strconv.Itoa(j+1) + " </td><td bgcolor=\"azure3\">" + strconv.Itoa(int(inodes[i].I_block[j])) + "</td></tr>\n"
			} else {
				content += "<tr><td bgcolor=\"aliceblue\">I_block " + strconv.Itoa(j+1) + " </td><td bgcolor=\"aliceblue\">" + strconv.Itoa(int(inodes[i].I_block[j])) + "</td></tr>\n"

			}
		}
		content += "<tr><td bgcolor=\"deepskyblue\">I_type</td><td bgcolor=\"deepskyblue\">" + strconv.Itoa(int(inodes[i].I_type)) + "</td></tr>\n"
		content += "<tr><td bgcolor=\"deepskyblue\">I_perm</td><td bgcolor=\"deepskyblue\">" + strconv.Itoa(int(inodes[i].I_perm)) + "</td></tr>\n"
		content += "</table>\n"
		content += ">]"
		content += "\n"
	}

	content += "\n"

	for i := 0; i < len(inodes); i++ {
		if i == 0 {
			content += "A" + strconv.Itoa(i)
		} else {
			content += " -> " + "A" + strconv.Itoa(i)
		}
	}

	content += "\n"
	content += "{ rank=same "
	for i := 0; i < len(inodes); i++ {
		content += "A" + strconv.Itoa(i) + " "
	}
	content += "}"

	content += "\n"
	content += "}\n"

	CreateFile(pd)
	WriteFile(content, pd)
	// termination := strings.Split(pathOut, ".")
	// Execute(pathOut, pd, termination[1], responseString)
	Message("REP", "Reporte de Inodos se ha generado correctamente en"+pathOut, responseString)
}

func repBlock(id string, pathOut string, responseString string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido", responseString)
		return
	}

	super := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	var path string
	partition := GetMount("MKGRP", id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("REP", "No se encontró la partición montada con el id: "+Logged.Id, responseString)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco", responseString)
		return
	}

	aux := strings.Split(pathOut, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos", responseString)
		return
	}
	pd := aux[0] + ".dot"

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("REP", "Error al leer el archivo", responseString)
		return
	}

	var inodes []Structs.Inodos
	inode = Structs.NewInodos()
	file.Seek(super.S_inode_start, 0)
	for i := 0; i < int(super.S_inodes_count); i++ {
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("REP", "Error al leer el archivo", responseString)
			return
		}

		if inode.I_uid == -1 {
			break
		}
		inodes = append(inodes, inode)
	}

	counter := 0
	content := "digraph Bloques{\n"
	content += "node [ shape=plaintext fontname=Arial ]\n"

	file.Seek(super.S_inode_start, 0)
	for v := 0; v < len(inodes); v++ {
		file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*int64(v), 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("REP", "Error al leer el archivo", responseString)
			return
		}
		if inode.I_type == 0 {
			for i := 0; i < 16; i++ {
				if i < 16 {
					if inode.I_block[i] != -1 {
						folder = Structs.NewDirectoriesBlocks()
						file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)

						data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &folder)
						if err_ != nil {
							Error("REP", "Error al leer el archivo", responseString)
							return
						}

						content += "\nA" + strconv.Itoa(counter)
						content += "[label= <"
						content += "<table border=\"1\" cellborder=\"0\">\n"
						content += "<tr><td bgcolor=\"coral3\" colspan=\"2\"><font color=\"white\">Bloque Carpeta " + strconv.Itoa(counter) + "</font></td></tr>\n"
						content += "<tr><td bgcolor=\"coral\">B_name</td><td bgcolor=\"coral1\">B_Inodo</td></tr>\n"
						for j := 0; j < 4; j++ {
							name := ""
							for nam := 0; nam < len(folder.B_content[j].B_name); nam++ {
								if folder.B_content[j].B_name[nam] == 0 {
									continue
								}
								name += string(folder.B_content[j].B_name[nam])
							}
							content += "<tr><td bgcolor=\"azure\">" + name + "</td><td bgcolor=\"azure\">" + strconv.Itoa(int(folder.B_content[j].B_inodo)) + "</td></tr>\n"
						}
						content += "</table>\n"
						content += ">]"
						counter++
					}
				}
			}
		} else if inode.I_type == 1 {
			for i := 0; i < 16; i++ {
				if i < 16 {
					if inode.I_block[i] != -1 {
						var folderAux Structs.FilesBlocks
						file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*inode.I_block[i], 0)

						data = readBytes(file, int(unsafe.Sizeof(Structs.FilesBlocks{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &folderAux)
						if err_ != nil {
							Error("REP", "Error al leer el archivo", responseString)
							return
						}

						content += "\nA" + strconv.Itoa(counter)
						content += "[label= <"
						content += "<table border=\"1\" cellborder=\"0\">\n"
						content += "<tr><td bgcolor=\"palegreen\"> Bloque Archivo " + strconv.Itoa(counter) + "</td></tr>\n"
						folderContent := ""
						for k := 0; k < len(folderAux.B_content); k++ {
							if folderAux.B_content[k] == 0 {
								continue
							}
							regex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúüñ.,]+$`)
							if !regex.MatchString(string(folderAux.B_content[k])) {
								continue
							} else {
								folderContent += string(folderAux.B_content[k])
							}
						}
						content += "<tr><td>" + folderContent + "</td></tr>\n"
						content += "</table>\n"
						content += ">]"
						counter++
					}
				}
			}
		} else {
			continue
		}
	}

	content += "\n"

	for i := 0; i < counter; i++ {
		if i == 0 {
			content += "A" + strconv.Itoa(i)
		} else {
			content += " -> " + "A" + strconv.Itoa(i)
		}
	}

	content += "\n"
	content += "{ rank=same "
	for i := 0; i < counter; i++ {
		content += "A" + strconv.Itoa(i) + " "
	}
	content += "}"
	content += "\n"
	content += "}\n"

	CreateFile(pd)
	WriteFile(content, pd)
	// termination := strings.Split(pathOut, ".")
	// Execute(pathOut, pd, termination[1], responseString)
	Message("REP", "Reporte de Bloques se ha generado correctamente en"+pathOut, responseString)
}

func repTree(id string, pathOut string, responseString string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido", responseString)
		return
	}

	super := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	var path string
	partition := GetMount("REP", id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("REP", "No se encontró la partición montada con el id: "+Logged.Id, responseString)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco", responseString)
		return
	}

	aux := strings.Split(pathOut, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos", responseString)
		return
	}
	pd := aux[0] + ".dot"

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("REP", "Error al leer el archivo", responseString)
		return
	}

	var inodes []Structs.Inodos
	inode = Structs.NewInodos()
	file.Seek(super.S_inode_start, 0)
	for i := 0; i < int(super.S_inodes_count); i++ {
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("REP", "Error al leer el archivo", responseString)
			return
		}

		if inode.I_uid == -1 {
			break
		}
		inodes = append(inodes, inode)
	}

	content := "digraph Tree{\n"
	content += "rankdir=LR;\n"
	content += "node [ shape=plaintext fontname=Arial ]\n"

	// GRAFICANDO INODOS
	file.Seek(super.S_inode_start, 0)
	for i := 0; i < len(inodes); i++ {
		pos := int(super.S_inode_start) + int(unsafe.Sizeof(Structs.Inodos{}))*i
		content += "I" + strconv.Itoa(pos)
		content += "[label= <"
		content += "<table border=\"1\" cellborder=\"0\">\n"
		content += "<tr><td bgcolor=\"dodgerblue4\" colspan=\"2\" ><font color=\"white\">Inodo " + strconv.Itoa(pos) + "</font></td></tr>\n"
		content += "<tr><td bgcolor=\"deepskyblue\">I_uid</td><td bgcolor=\"deepskyblue\">" + strconv.Itoa(int(inodes[i].I_uid)) + "</td></tr>\n"
		content += "<tr><td bgcolor=\"deepskyblue\">I_gid</td><td bgcolor=\"deepskyblue\">" + strconv.Itoa(int(inodes[i].I_gid)) + "</td></tr>\n"
		content += "<tr><td bgcolor=\"deepskyblue\">I_s</td><td bgcolor=\"deepskyblue\">" + strconv.Itoa(int(inodes[i].I_s)) + "</td></tr>\n"
		atime := ""
		for k := 0; k < len(inodes[i].I_atime); k++ {
			if inodes[i].I_atime[k] != 0 {
				atime += string(inodes[i].I_atime[k])
			}
		}
		content += "<tr><td bgcolor=\"deepskyblue\">I_atime</td><td bgcolor=\"deepskyblue\">" + atime + "</td></tr>\n"
		ctime := ""
		for k := 0; k < len(inodes[i].I_ctime); k++ {
			if inodes[i].I_ctime[k] != 0 {
				ctime += string(inodes[i].I_ctime[k])
			}
		}
		content += "<tr><td bgcolor=\"deepskyblue\">I_ctime</td><td bgcolor=\"deepskyblue\">" + ctime + "</td></tr>\n"
		mtime := ""
		for k := 0; k < len(inodes[i].I_mtime); k++ {
			if inodes[i].I_mtime[k] != 0 {
				mtime += string(inodes[i].I_mtime[k])
			}
		}
		content += "<tr><td bgcolor=\"deepskyblue\">I_mtime</td><td bgcolor=\"deepskyblue\">" + mtime + "</td></tr>\n"
		for j := 0; j < len(inodes[i].I_block); j++ {
			if j > 12 {
				content += "<tr><td bgcolor=\"azure3\">I_block " + strconv.Itoa(j+1) + " </td><td port=\"b" + strconv.Itoa(j) + "\" bgcolor=\"azure3\">" + strconv.Itoa(int(inodes[i].I_block[j])) + "</td></tr>\n"
			} else {
				content += "<tr><td bgcolor=\"aliceblue\">I_block " + strconv.Itoa(j+1) + " </td><td port=\"b" + strconv.Itoa(j) + "\" bgcolor=\"aliceblue\">" + strconv.Itoa(int(inodes[i].I_block[j])) + "</td></tr>\n"
			}
		}
		content += "<tr><td bgcolor=\"deepskyblue\">I_type</td><td bgcolor=\"deepskyblue\">" + strconv.Itoa(int(inodes[i].I_type)) + "</td></tr>\n"
		content += "<tr><td bgcolor=\"deepskyblue\">I_perm</td><td bgcolor=\"deepskyblue\">" + strconv.Itoa(int(inodes[i].I_perm)) + "</td></tr>\n"
		content += "</table>\n"
		content += ">]"
		content += "\n"
	}

	content += "\n"

	// GRAFICANDO BLOQUES
	file.Seek(super.S_inode_start, 0)
	for i := 0; i < len(inodes); i++ {
		file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*int64(i), 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("REP", "Error al leer el archivo", responseString)
			return
		}
		if inode.I_type == 0 {
			for j := 0; j < 16; j++ {
				if j < 16 {
					if inode.I_block[j] != -1 {
						folder = Structs.NewDirectoriesBlocks()
						file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[j]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[j], 0)
						posBloque := int(super.S_block_start) + int(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*int(inode.I_block[j]) + int(unsafe.Sizeof(Structs.FilesBlocks{}))*32*int(inode.I_block[j])
						data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &folder)
						if err_ != nil {
							Error("REP", "Error al leer el archivo", responseString)
							return
						}

						content += "B" + strconv.Itoa(posBloque)
						content += " [label= <"
						content += "<table border=\"1\" cellborder=\"0\">\n"
						content += "<tr><td bgcolor=\"coral3\" colspan=\"2\" ><font color=\"white\">Bloque Carpeta " + strconv.Itoa(posBloque) + "</font></td></tr>\n"
						content += "<tr><td bgcolor=\"coral\">B_name</td><td bgcolor=\"coral1\">B_Inodo</td></tr>\n"
						for k := 0; k < 4; k++ {
							name := ""
							for nam := 0; nam < len(folder.B_content[k].B_name); nam++ {
								if folder.B_content[k].B_name[nam] == 0 {
									continue
								}
								name += string(folder.B_content[k].B_name[nam])
							}
							if k > 1 {
								content += "<tr><td bgcolor=\"azure\">" + name + "</td><td port=\"i" + strconv.Itoa(k) + "\" bgcolor=\"azure\">" + strconv.Itoa(int(folder.B_content[k].B_inodo)) + "</td></tr>\n"
							} else {
								content += "<tr><td port=\"i" + strconv.Itoa(k) + "\" bgcolor=\"azure\">" + name + "</td><td bgcolor=\"azure\">" + strconv.Itoa(int(folder.B_content[k].B_inodo)) + "</td></tr>\n"
							}
						}
						content += "</table>\n"
						content += ">]\n"
						content += "\n"
					}
				}
			}
		} else if inode.I_type == 1 {
			for j := 0; j < 16; j++ {
				if j < 16 {
					if inode.I_block[j] != -1 {
						var folderAux Structs.FilesBlocks
						file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*inode.I_block[j], 0)
						posBloque := int(super.S_block_start) + int(unsafe.Sizeof(Structs.DirectoriesBlocks{})) + int(unsafe.Sizeof(Structs.FilesBlocks{}))*int(inode.I_block[j])
						data = readBytes(file, int(unsafe.Sizeof(Structs.FilesBlocks{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &folderAux)
						if err_ != nil {
							Error("REP", "Error al leer el archivo", responseString)
							return
						}

						content += "B" + strconv.Itoa(posBloque)
						content += " [label= <"
						content += "<table border=\"1\" cellborder=\"0\">\n"
						content += "<tr><td bgcolor=\"palegreen\"> Bloque Archivo " + strconv.Itoa(posBloque) + "</td></tr>\n"
						folderContent := ""
						for k := 0; k < len(folderAux.B_content); k++ {
							if folderAux.B_content[k] == 0 {
								continue
							}
							regex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúüñ.,]+$`)
							if !regex.MatchString(string(folderAux.B_content[k])) {
								continue
							} else {
								folderContent += string(folderAux.B_content[k])
							}
						}
						content += "<tr><td>" + folderContent + "</td></tr>\n"
						content += "</table>\n"
						content += ">]\n"
						content += "\n"
					}
				}
			}
		} else {
			continue
		}
	}

	content += "\n"

	// UNIENDO INODOS CON CARPETAS
	file.Seek(super.S_inode_start, 0)
	for i := 0; i < len(inodes); i++ {
		posInodo := int(super.S_inode_start) + int(unsafe.Sizeof(Structs.Inodos{}))*i
		file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*int64(i), 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("REP", "Error al leer el archivo", responseString)
			return
		}
		if inode.I_type == 0 {
			for j := 0; j < 16; j++ {
				if j < 16 {
					if inode.I_block[j] != -1 {
						folder = Structs.NewDirectoriesBlocks()
						file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[j]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[j], 0)
						posBloque := int(super.S_block_start) + int(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*int(inode.I_block[j]) + int(unsafe.Sizeof(Structs.FilesBlocks{}))*32*int(inode.I_block[j])
						data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &folder)
						if err_ != nil {
							Error("REP", "Error al leer el archivo", responseString)
							return
						}

						// UNIENDO CARPETAS CON ELLA MISMA
						content += "B" + strconv.Itoa(posBloque) + ":i" + strconv.Itoa(0)
						content += " -> "
						content += "B" + strconv.Itoa(posBloque)
						content += "\n"

						// UNIENDO INIDO CON CARPETA HIJO
						content += "I" + strconv.Itoa(posInodo) + ":b" + strconv.Itoa(j)
						content += " -> "
						content += "B" + strconv.Itoa(posBloque)
						content += "\n"

						// UNIENDO CARPETAS CON INODOS PADRES
						content += "B" + strconv.Itoa(posBloque) + ":i" + strconv.Itoa(1)
						content += " -> "
						content += "I" + strconv.Itoa(posInodo) + ":b" + strconv.Itoa(j)
						content += "\n"

						// UNIENDO CARPETAS CON INODOS HIJOS
						folder = Structs.NewDirectoriesBlocks()
						file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inodes[i].I_block[j]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inodes[i].I_block[j], 0)
						data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &folder)
						if err_ != nil {
							Error("REP", "Error al leer el archivo", responseString)
							return
						}

						for k := 0; k < 4; k++ {
							name := ""
							for nam := 0; nam < len(folder.B_content[k].B_name); nam++ {
								if folder.B_content[k].B_name[nam] == 0 {
									continue
								}
								name += string(folder.B_content[k].B_name[nam])
							}
							if Compare(name, ".") || Compare(name, "..") || Compare(name, "-") {
								continue
							} else {
								if k == 2 && folder.B_content[2].B_inodo != -1 {
									content += "B" + strconv.Itoa(posBloque) + ":i" + strconv.Itoa(2)
									content += " -> "
									content += "I" + strconv.Itoa(int(super.S_inode_start)+int(unsafe.Sizeof(Structs.Inodos{}))*int(folder.B_content[2].B_inodo))
									content += "\n"
								} else if k == 3 && folder.B_content[3].B_inodo != -1 {
									content += "B" + strconv.Itoa(posBloque) + ":i" + strconv.Itoa(3)
									content += " -> "
									content += "I" + strconv.Itoa(int(super.S_inode_start)+int(unsafe.Sizeof(Structs.Inodos{}))*int(folder.B_content[3].B_inodo))
									content += "\n"
								}
							}
						}
					}
				}
			}
		} else if inodes[i].I_type == 1 {
			for j := 0; j < 16; j++ {
				if j < 16 {
					if inodes[i].I_block[j] != -1 {
						var folderAux Structs.FilesBlocks
						file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*inode.I_block[j], 0)
						posBloque := int(super.S_block_start) + int(unsafe.Sizeof(Structs.DirectoriesBlocks{})) + int(unsafe.Sizeof(Structs.FilesBlocks{}))*int(inode.I_block[j])
						data = readBytes(file, int(unsafe.Sizeof(Structs.FilesBlocks{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &folderAux)
						if err_ != nil {
							Error("REP", "Error al leer el archivo", responseString)
							return
						}

						content += "I" + strconv.Itoa(posInodo) + ":b" + strconv.Itoa(j)
						content += " -> "
						content += "B" + strconv.Itoa(posBloque)
						content += "\n"
					}
				}
			}
		} else {
			continue
		}
		content += "\n"
	}

	content += "\n"
	content += "}\n"

	CreateFile(pd)
	WriteFile(content, pd)
	// termination := strings.Split(pathOut, ".")
	// Execute(pathOut, pd, termination[1], responseString)
	Message("REP", "Reporte de Tree se ha generado correctamente en"+pathOut, responseString)
}

func repJournaling(id string, pathOut string, responseString string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido", responseString)
		return
	}

	super := Structs.NewSuperBlock()

	var path string
	partition := GetMount("REP", id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("REP", "No se encontró la partición montada con el id: "+Logged.Id, responseString)
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco", responseString)
		return
	}

	aux := strings.Split(pathOut, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos", responseString)
		return
	}
	pd := aux[0] + ".dot"

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("REP", "Error al leer el archivo", responseString)
		return
	}

	if Compare(strconv.Itoa(int(super.S_filesystem_type)), "2") {
		Error("REP", "El sistema de archivo es de tipo EXT2, no maneja la estructura de Journaling", responseString)
		return
	}

	textJ := ""
	jour := Structs.NewJournaling()
	file.Seek(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBlock{})), 0)
	for i := 0; i < int(super.S_inodes_count); i++ {
		file.Seek(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBlock{}))+int64(unsafe.Sizeof(Structs.Journaling{}))*int64(i), 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.Journaling{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &jour)
		if err_ != nil {
			Error("MKGRP", "Error al leer el archivo", responseString)
			return
		}

		pathJournaling := ""
		for k := 0; k < len(jour.Path); k++ {
			if jour.Path[k] != 0 {
				pathJournaling += string(jour.Path[k])
			}
		}

		if pathJournaling == "-" {
			break
		} else {
			pathJournaling = ""
			for k := 0; k < len(jour.Path); k++ {
				if jour.Path[k] != 0 {
					pathJournaling += string(jour.Path[k])
				}
			}

			operationJournaling := ""
			for k := 0; k < len(jour.Operation); k++ {
				if jour.Operation[k] != 0 {
					operationJournaling += string(jour.Operation[k])
				}
			}

			contentJournaling := ""
			for k := 0; k < len(jour.Content); k++ {
				if jour.Content[k] != 0 {
					contentJournaling += string(jour.Content[k])
				}
			}

			dateJournaling := ""
			for k := 0; k < len(jour.Date); k++ {
				if jour.Date[k] != 0 {
					dateJournaling += string(jour.Date[k])
				}
			}

			textJ += "<tr><td>" + operationJournaling + "</td><td>" + pathJournaling + "</td><td>" + contentJournaling + "</td><td>" + dateJournaling + "</td></tr>\n"

		}

	}

	text := "digraph Journaling{\n"
	text += "node [ shape=none fontname=Arial ]\n"
	text += "n1 [ label = <\n"
	text += "<table>\n"
	text += "<tr><td colspan=\"4\" bgcolor=\"palegreen4\"><font color=\"white\">REPORTE DE JOURNALING</font></td></tr>\n"
	text += "<tr><td>Operacion</td><td>Path</td><td>Contenido</td><td>Fecha</td></tr>\n"
	text += textJ
	text += "</table>\n"
	text += "> ]\n"
	text += "}\n"

	file.Close()

	CreateFile(pd)
	WriteFile(text, pd)
	// termination := strings.Split(pathOut, ".")
	// Execute(pathOut, pd, termination[1], responseString)
	Message("REP", "Reporte de Journaling se ha generado correctamente en"+pathOut, responseString)
}

func repCat(id string, pathOut string, file string, responseString string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido", responseString)
		return
	}

	var path string
	partition := GetMount("REP", id, &path, responseString)
	if string(partition.Part_status) == "0" {
		Error("REP", "No se encontró la partición montada con el id: "+id, responseString)
		return
	}
	content := ""
	CreateFile(pathOut + "txt")

	tmp := GetPath(file)
	content = cat(tmp, partition, path, responseString)
	WriteFile(content, pathOut+"txt")
	Message("REP", "Reporte de CAT se ha generado correctamente en"+pathOut, responseString)
}

func repLs(id string, pathOut string, responseString string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("REP", "El primer identificador no es válido", responseString)
		return
	}

	super := Structs.NewSuperBlock()
	inode := Structs.NewInodos()
	folder := Structs.NewDirectoriesBlocks()

	var path2 string
	partition := GetMount("REP", id, &path2, responseString)
	if string(partition.Part_status) == "0" {
		Error("REP", "No se encontró la partición montada con el id: "+Logged.Id, responseString)
		return
	}

	aux := strings.Split(pathOut, ".")
	if len(aux) > 2 {
		Error("REP", "No se admiten nombres de archivos que contengan puntos", responseString)
		return
	}
	pd := aux[0] + ".dot"

	file, err := os.Open(strings.ReplaceAll(path2, "\"", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco", responseString)
		return
	}

	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("REP", "Error al leer el archivo", responseString)
		return
	}

	CreateFile(pd)

	var inodes []Structs.Inodos
	inode = Structs.NewInodos()
	file.Seek(super.S_inode_start, 0)
	for i := 0; i < int(super.S_inodes_count); i++ {
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("REP", "Error al leer el archivo", responseString)
			return
		}

		if inode.I_uid == -1 {
			break
		}
		inodes = append(inodes, inode)
	}

	textL := ""
	file.Seek(super.S_inode_start, 0)
	for v := 0; v < len(inodes); v++ {
		file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*int64(v), 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("REP", "Error al leer el archivo", responseString)
			return
		}
		if inode.I_type == 0 {
			for i := 0; i < 16; i++ {
				if i < 16 {
					if inode.I_block[i] != -1 {
						folder = Structs.NewDirectoriesBlocks()
						file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))*inode.I_block[i]+int64(unsafe.Sizeof(Structs.FilesBlocks{}))*32*inode.I_block[i], 0)

						data = readBytes(file, int(unsafe.Sizeof(Structs.DirectoriesBlocks{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &folder)
						if err_ != nil {
							Error("MKDIR", "Error al leer el archivo", responseString)
							return
						}

						if folder.B_content[2].B_inodo != -1 {
							bi := folder.B_content[2].B_inodo

							name := ""
							for nam := 0; nam < len(folder.B_content[2].B_name); nam++ {
								if folder.B_content[2].B_name[nam] == 0 {
									continue
								}
								name += string(folder.B_content[2].B_name[nam])
							}

							inodeForRead := Structs.NewInodos()
							file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
							data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
							buffer = bytes.NewBuffer(data)
							err_ = binary.Read(buffer, binary.BigEndian, &inodeForRead)
							if err_ != nil {
								Error("REP", "Error al leer el archivo", responseString)
								return
							}

							tipo := ""
							if strconv.Itoa(int(inodeForRead.I_type)) == "0" {
								tipo = "Carpeta"
							} else {
								tipo = "Archivo"
							}

							fecha := ""
							for nam := 0; nam < len(inodeForRead.I_mtime); nam++ {
								if inodeForRead.I_mtime[nam] == 0 {
									continue
								}
								fecha += string(inodeForRead.I_mtime[nam])
							}

							textL += "<tr><td>" + strconv.Itoa(int(inodeForRead.I_perm)) + "</td><td>" + strconv.Itoa(int(inodeForRead.I_uid)) + "</td><td>" + strconv.Itoa(int(inodeForRead.I_gid)) + "</td><td>" + strconv.Itoa(int(inodeForRead.I_s)) + "</td><td>" + fecha + "</td><td>" + tipo + "</td><td>" + name + "</td></tr>\n"

						}

						if folder.B_content[3].B_inodo != -1 {
							bi := folder.B_content[3].B_inodo

							name := ""
							for nam := 0; nam < len(folder.B_content[3].B_name); nam++ {
								if folder.B_content[3].B_name[nam] == 0 {
									continue
								}
								name += string(folder.B_content[3].B_name[nam])
							}

							inodeForRead := Structs.NewInodos()
							file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*bi, 0)
							data = readBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
							buffer = bytes.NewBuffer(data)
							err_ = binary.Read(buffer, binary.BigEndian, &inodeForRead)
							if err_ != nil {
								Error("REP", "Error al leer el archivo", responseString)
								return
							}

							tipo := ""
							if strconv.Itoa(int(inodeForRead.I_type)) == "0" {
								tipo = "Carpeta"
							} else {
								tipo = "Archivo"
							}

							fecha := ""
							for nam := 0; nam < len(inodeForRead.I_mtime); nam++ {
								if inodeForRead.I_mtime[nam] == 0 {
									continue
								}
								fecha += string(inodeForRead.I_mtime[nam])
							}

							textL += "<tr><td>" + strconv.Itoa(int(inodeForRead.I_perm)) + "</td><td>" + strconv.Itoa(int(inodeForRead.I_uid)) + "</td><td>" + strconv.Itoa(int(inodeForRead.I_gid)) + "</td><td>" + strconv.Itoa(int(inodeForRead.I_s)) + "</td><td>" + fecha + "</td><td>" + tipo + "</td><td>" + name + "</td></tr>\n"

						}

					}
				}
			}
		} else if inode.I_type == 1 {
			continue
		} else {
			continue
		}
	}

	text := "digraph LS{\n"
	text += "node [ shape=none fontname=Arial ]\n"
	text += "n1 [ label = <\n"
	text += "<table>\n"
	text += "<tr><td colspan=\"7\" bgcolor=\"midnightblue\"><font color=\"white\">REPORTE LS</font></td></tr>\n"
	text += "<tr><td>Permisos</td><td>Propietario</td><td>Grupo</td><td>Tamaño</td><td>Fecha</td><td>Tipo</td><td>Nombre</td></tr>\n"
	text += textL
	text += "</table>\n"
	text += "> ]\n"
	text += "}\n"

	WriteFile(text, pd)
	// termination := strings.Split(pathOut, ".")
	// Execute(pathOut, pd, termination[1], responseString)
	Message("REP", "Reporte de Bloques se ha generado correctamente en"+pathOut, responseString)

}

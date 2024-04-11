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

type Transition struct {
	partition int
	start     int
	end       int
	before    int
	after     int
}

var startValue int

func DataFDISK(tokens []string, responseString string) {
	size := ""
	unit := "k"
	driveLetter := ""
	tipo := "P"
	fit := "WF"
	name := ""
	add := ""
	deleteP := ""
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "size") {
			size = tk[1]
		} else if Compare(tk[0], "unit") {
			unit = tk[1]
		} else if Compare(tk[0], "driveletter") {
			currentPath, _ := os.Getwd()
			driveLetter = currentPath + "/MIA/P2/" + tk[1] + ".dsk"
		} else if Compare(tk[0], "type") {
			tipo = tk[1]
		} else if Compare(tk[0], "fit") {
			fit = tk[1]
		} else if Compare(tk[0], "name") {
			name = tk[1]
		} else if Compare(tk[0], "delete") {
			deleteP = tk[1]
		} else if Compare(tk[0], "add") {
			add = tk[1]
		}
	}

	if driveLetter == "" || name == "" {
		Error("FDISK", "EL comando FDISK necesita parámetros obligatorios", responseString)
		return
	} else {
		if size == "" && deleteP == "" && add == "" {
			Error("FDISK", "EL comando FDISK necesita parámetros obligatorios", responseString)
			return
		} else if size != "" && deleteP == "" && add == "" {
			generatePartition(size, unit, driveLetter, tipo, fit, name, responseString)
		} else if size == "" && deleteP != "" && add == "" {
			deletePartition(deleteP, driveLetter, name, responseString)
		} else if size == "" && deleteP == "" && add != "" {
			// agregar
		}
	}
}

func generatePartition(s string, u string, d string, t string, f string, n string, responseString string) {
	startValue = 0
	i, error_ := strconv.Atoi(s)
	if error_ != nil {
		Error("FDISK", "Size debe ser un número entero", responseString)
		return
	}
	if i <= 0 {
		Error("FDISK", "Size debe ser mayor que 0", responseString)
		return
	}
	if Compare(u, "b") || Compare(u, "k") || Compare(u, "m") {
		if Compare(u, "k") {
			i = i * 1024
		} else if Compare(u, "m") {
			i = i * 1024 * 1024
		}
	} else {
		Error("FDISK", "Unit no contiene los valores esperados", responseString)
		return
	}
	if !(Compare(t, "p") || Compare(t, "e") || Compare(t, "l")) {
		Error("FDISK", "Type no contiene los valores esperados", responseString)
		return
	}
	if !(Compare(f, "bf") || Compare(f, "ff") || Compare(f, "wf")) {
		Error("FDISK", "Fit no contiene los valores esperados", responseString)
	}

	file, err := os.OpenFile(strings.ReplaceAll(d, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("FDISK", "Error al abrir el archivo", responseString)
		return
	}

	mbr := readDisk(d, responseString)

	if int(mbr.Mbr_tamano) < i {
		Error("FDISK", "EL TAMAÑO DE LA PARTICIÓN ES MAYOR QUE EL TAMAÑO DEL DISCO", responseString)
		return
	}

	partitions := GetPartitions(*mbr)
	var between []Transition
	used := 0
	ext := 0
	logic := 0
	c := 0

	base := int(unsafe.Sizeof(Structs.MBR{}))
	extended := Structs.NewPartition()

	for j := 0; j < len(partitions); j++ {
		prttn := partitions[j]
		if prttn.Part_status == '1' {
			var trn Transition
			trn.partition = c
			trn.start = int(prttn.Part_start)
			trn.end = int(prttn.Part_start + prttn.Part_s)
			trn.before = trn.start - base
			base = trn.end
			if used != 0 {
				between[used-1].after = trn.start - (between[used-1].end)
			}
			between = append(between, trn)
			used++
			if prttn.Part_type == "e"[0] || prttn.Part_type == "E"[0] {
				ext++
				extended = prttn
			}
		}
		if used == 4 && !Compare(t, "l") {
			Error("FDISK", "Límite de particioens alcanzado", responseString)
			return
		} else if ext == 1 && Compare(t, "e") {
			Error("FDISK", "Solo se puede crear una partición extendida", responseString)
			return
		}
		c++
	}
	if ext == 0 && Compare(t, "l") {
		Error("FDISK", "Aún no se han creado particiones extendidas, no se puede agregar una lógica", responseString)
		return
	}

	if used != 0 {
		between[len(between)-1].after = int(mbr.Mbr_tamano) - between[len(between)-1].end
	}

	comeBack := SearchPartitions(*mbr, n, d, responseString)
	if comeBack != nil {
		Error("FDISK", "El nombre "+n+", ya está en uso", responseString)
		return
	}

	if Compare(t, "l") {
		logic++
	}

	temporal := Structs.NewPartition()
	temporal.Part_status = '1'
	temporal.Part_s = int64(i)
	temporal.Part_type = strings.ToUpper(t)[0]
	temporal.Part_fit = strings.ToUpper(f)[0]
	copy(temporal.Part_name[:], n)
	temporal.Part_correlative = int64(used + ext + logic + 1)
	if Compare(t, "l") {
		Logic(temporal, extended, d, n, responseString)
		return
	}

	mbr = fitF(*mbr, temporal, between, partitions, used, responseString)
	if mbr == nil {
		return
	}

	file, err = os.OpenFile(strings.ReplaceAll(d, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("FDISK", "Error al abrir el archivo", responseString)
	}
	file.Seek(0, 0)
	var binary2 bytes.Buffer
	binary.Write(&binary2, binary.BigEndian, mbr)
	WritingBytes(file, binary2.Bytes())
	if Compare(t, "E") {
		ebr := Structs.NewEBR()
		ebr.Part_mount = '0'
		ebr.Part_start = int64(startValue)
		ebr.Part_s = 0
		ebr.Part_next = -1

		file.Seek(int64(startValue), 0)
		var binary3 bytes.Buffer
		binary.Write(&binary3, binary.BigEndian, ebr)
		WritingBytes(file, binary3.Bytes())
		Message("FDISK", "Partición Extendida: "+n+", creada correctamente", responseString)
		return
	}
	file.Close()
	Message("FDISK", "Partición Primaria: "+n+", creada correctamente", responseString)
}

func deletePartition(de string, d string, n string, responseString string) {
	if !Compare(de, "full") {
		Error("FDISK", "Delete no contiene los valores esperados", responseString)
		return
	}

	file, err := os.OpenFile(strings.ReplaceAll(d, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("FDISK", "Error al abrir el archivo", responseString)
		return
	}

	var partitions [4]Structs.Partition
	mbr := readDisk(d, responseString)
	partitions[0] = mbr.Mbr_partitions_1
	partitions[1] = mbr.Mbr_partitions_2
	partitions[2] = mbr.Mbr_partitions_3
	partitions[3] = mbr.Mbr_partitions_4

	file, err = os.OpenFile(strings.ReplaceAll(d, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("FDISK", "Error al abrir el archivo", responseString)
		return
	}

	founded := false
	c := 0
	ext := false
	deletedL := true
	extended := Structs.NewPartition()
	for i := 0; i < len(partitions); i++ {
		partition := partitions[i]
		if partition.Part_status == "1"[0] {
			nameP := ""
			for j := 0; j < len(partition.Part_name); j++ {
				if partition.Part_name[j] != 0 {
					nameP += string(partition.Part_name[j])
				}
			}
			if Compare(nameP, n) {
				deletedL = false
				founded = true
				c = i
				var zero int8 = 0
				size := int(partition.Part_s)
				for j := 0; j < size; j++ {
					file.Seek(partition.Part_start+int64(j), 0)
					var binaryZero bytes.Buffer
					binary.Write(&binaryZero, binary.BigEndian, zero)
					WritingBytes(file, binaryZero.Bytes())
				}
			} else if partition.Part_type == "E"[0] || partition.Part_type == "e"[0] {
				ext = true
				extended = partition
			}
		}
	}

	if ext {
		ebrs := GetLogics(extended, d, responseString)
		for i := 0; i < len(ebrs); i++ {
			ebr := ebrs[i]
			if ebr.Part_mount == '1' {
				nameE := ""
				for j := 0; j < len(ebr.Part_name); j++ {
					if ebr.Part_name[j] != 0 {
						nameE += string(ebr.Part_name[j])
					}
				}
				if Compare(nameE, n) {
					logicPartitionStart := ebr.Part_start + int64(unsafe.Sizeof(Structs.EBR{}))
					var zero int8 = 0
					newEbr := Structs.NewEBR()
					newEbr.Part_fit = '0'
					newEbr.Part_start = ebr.Part_start
					newEbr.Part_next = ebr.Part_next
					file.Seek(ebr.Part_start, 0)
					var binaryEbr bytes.Buffer
					binary.Write(&binaryEbr, binary.BigEndian, newEbr)
					WritingBytes(file, binaryEbr.Bytes())

					size := int(ebr.Part_s)
					for j := 0; j < size; j++ {
						file.Seek(logicPartitionStart+int64(j), 0)
						var binaryZero bytes.Buffer
						binary.Write(&binaryZero, binary.BigEndian, zero)
						WritingBytes(file, binaryZero.Bytes())
					}
					deletedL = false
					Message("FDISK", "Particion Lógica "+n+", eliminada correctamente", responseString)
					return
				}
			}
		}
	}

	if founded && !deletedL {
		if c == 0 {
			mbr.Mbr_partitions_1 = Structs.NewPartition()
			mbr.Mbr_partitions_1.Part_fit = '0'
			mbr.Mbr_partitions_1.Part_type = '0'
			mbr.Mbr_partitions_1.Part_status = '0'
			mbr.Mbr_partitions_1.Part_start = -1
			mbr.Mbr_partitions_1.Part_s = 0
		} else if c == 1 {
			mbr.Mbr_partitions_2 = Structs.NewPartition()
			mbr.Mbr_partitions_2.Part_fit = '0'
			mbr.Mbr_partitions_2.Part_type = '0'
			mbr.Mbr_partitions_2.Part_status = '0'
			mbr.Mbr_partitions_2.Part_start = -1
			mbr.Mbr_partitions_2.Part_s = 0
		} else if c == 2 {
			mbr.Mbr_partitions_3 = Structs.NewPartition()
			mbr.Mbr_partitions_3.Part_fit = '0'
			mbr.Mbr_partitions_3.Part_type = '0'
			mbr.Mbr_partitions_3.Part_status = '0'
			mbr.Mbr_partitions_3.Part_start = -1
			mbr.Mbr_partitions_3.Part_s = 0
		} else if c == 3 {
			mbr.Mbr_partitions_4 = Structs.NewPartition()
			mbr.Mbr_partitions_4.Part_fit = '0'
			mbr.Mbr_partitions_4.Part_type = '0'
			mbr.Mbr_partitions_4.Part_status = '0'
			mbr.Mbr_partitions_4.Part_start = -1
			mbr.Mbr_partitions_4.Part_s = 0
		}
		file.Seek(0, 0)
		var binary2 bytes.Buffer
		binary.Write(&binary2, binary.BigEndian, mbr)
		WritingBytes(file, binary2.Bytes())
		Message("FDISK", "Particion "+n+", eliminada correctamente", responseString)
	}
}

func GetPartitions(disk Structs.MBR) []Structs.Partition {
	var v []Structs.Partition
	v = append(v, disk.Mbr_partitions_1)
	v = append(v, disk.Mbr_partitions_2)
	v = append(v, disk.Mbr_partitions_3)
	v = append(v, disk.Mbr_partitions_4)
	return v
}

func SearchPartitions(mbr Structs.MBR, name string, path string, responseString string) *Structs.Partition {
	var partitions [4]Structs.Partition
	partitions[0] = mbr.Mbr_partitions_1
	partitions[1] = mbr.Mbr_partitions_2
	partitions[2] = mbr.Mbr_partitions_3
	partitions[3] = mbr.Mbr_partitions_4

	ext := false
	extended := Structs.NewPartition()
	for i := 0; i < len(partitions); i++ {
		partition := partitions[i]
		if partition.Part_status == "1"[0] {
			nameP := ""
			for j := 0; j < len(partition.Part_name); j++ {
				if partition.Part_name[j] != 0 {
					nameP += string(partition.Part_name[j])
				}
			}
			if Compare(nameP, name) {
				return &partition
			} else if partition.Part_type == "E"[0] || partition.Part_type == "e"[0] {
				ext = true
				extended = partition
			}
		}
	}

	if ext {
		ebrs := GetLogics(extended, path, responseString)
		for i := 0; i < len(ebrs); i++ {
			ebr := ebrs[i]
			if ebr.Part_mount == '1' {
				nameE := ""
				for j := 0; j < len(ebr.Part_name); j++ {
					if ebr.Part_name[j] != 0 {
						nameE += string(ebr.Part_name[j])
					}
				}
				if Compare(nameE, name) {
					tmp := Structs.NewPartition()
					tmp.Part_status = '1'
					tmp.Part_type = 'L'
					tmp.Part_fit = ebr.Part_fit
					tmp.Part_start = ebr.Part_start
					tmp.Part_s = ebr.Part_s
					copy(tmp.Part_name[:], ebr.Part_name[:])
					return &tmp
				}
			}
		}
	}
	return nil
}

func GetLogics(partition Structs.Partition, path string, responseString string) []Structs.EBR {
	var ebrs []Structs.EBR
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("FDISK", "Error al abrir el archivo", responseString)
		return nil
	}
	file.Seek(0, 0)
	tmp := Structs.NewEBR()
	file.Seek(partition.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &tmp)
	if err_ != nil {
		Error("FDISK", "Error al leer el archivo", responseString)
		return nil
	}
	for {
		if int(tmp.Part_next) != -1 && int(tmp.Part_mount) != 0 {
			ebrs = append(ebrs, tmp)
			file.Seek(tmp.Part_next, 0)
			data = readBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &tmp)
			if err_ != nil {
				Error("FDISK", "Error al leer el archivo", responseString)
				return nil
			}
		} else {
			file.Close()
			break
		}
	}
	return ebrs
}

func Logic(p Structs.Partition, e Structs.Partition, d string, n string, responseString string) {
	logic := Structs.NewEBR()
	logic.Part_mount = '1'
	logic.Part_fit = p.Part_fit
	logic.Part_s = p.Part_s
	logic.Part_next = -1
	copy(logic.Part_name[:], p.Part_name[:])

	file, err := os.Open(strings.ReplaceAll(d, "\"", ""))
	if err != nil {
		Error("FDISK", "Error al abrir el archivo del disco", responseString)
		return
	}

	file.Seek(0, 0)

	tmp := Structs.NewEBR()
	tmp.Part_mount = 0
	tmp.Part_s = 0
	tmp.Part_next = -1
	file.Seek(e.Part_start, 0)

	data := readBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &tmp)
	if err_ != nil {
		Error("FDISK", "Error al leer el archivo", responseString)
		return
	}
	if err != nil {
		Error("FDISK", "Error al abrir el archivo del disco 3", responseString)
		return
	}
	var size int64 = 0
	file.Close()
	for {
		size += int64(unsafe.Sizeof(Structs.EBR{})) + tmp.Part_s
		if (tmp.Part_s == 0 && tmp.Part_next == -1) || (tmp.Part_s == 0 && tmp.Part_next == 0) {
			file, err = os.OpenFile(strings.ReplaceAll(d, "\"", ""), os.O_WRONLY, os.ModeAppend)
			logic.Part_start = tmp.Part_start
			logic.Part_next = logic.Part_start + logic.Part_s + int64(unsafe.Sizeof(Structs.EBR{}))
			if (e.Part_s - size) <= logic.Part_s {
				Error("FDISK", "No hay espacio para más particiones lógicas", responseString)
				return
			}
			file.Seek(logic.Part_start, 0)

			var binary2 bytes.Buffer
			binary.Write(&binary2, binary.BigEndian, logic)
			WritingBytes(file, binary2.Bytes())
			nameL := ""
			for j := 0; j < len(p.Part_name); j++ {
				nameL += string(p.Part_name[j])
			}
			file.Seek(logic.Part_next, 0)
			addLogic := Structs.NewEBR()
			addLogic.Part_mount = '0'
			addLogic.Part_next = -1
			addLogic.Part_start = logic.Part_next

			file.Seek(addLogic.Part_start, 0)

			var binary3 bytes.Buffer
			binary.Write(&binary3, binary.BigEndian, addLogic)
			WritingBytes(file, binary3.Bytes())

			Message("FDISK", "Partición Lógica: "+n+", creada correctamente", responseString)
			file.Close()
			return
		}
		file, err = os.Open(strings.ReplaceAll(d, "\"", ""))
		if err != nil {
			Error("FDISK", "Error al abrir el archivo del disco", responseString)
			return
		}
		file.Seek(tmp.Part_next, 0)
		data = readBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &tmp)
		if err_ != nil {
			Error("FDISK", "Error al leer el archivo", responseString)
			return
		}
	}
}

func fitF(mbr Structs.MBR, p Structs.Partition, t []Transition, ps []Structs.Partition, u int, responseString string) *Structs.MBR {
	if u == 0 {
		p.Part_start = int64(unsafe.Sizeof(mbr))
		startValue = int(p.Part_start)
		mbr.Mbr_partitions_1 = p
		return &mbr
	} else {
		var use Transition
		c := 0
		for i := 0; i < len(t); i++ {
			tr := t[i]
			if c == 0 {
				use = tr
				c++
				continue
			}

			if Compare(string(mbr.Dsk_fit[0]), "F") {
				if int64(use.before) >= p.Part_s || int64(use.after) >= p.Part_s {
					break
				}
				use = tr
			} else if Compare(string(mbr.Dsk_fit[0]), "B") {
				if int64(tr.before) >= p.Part_s || int64(use.after) < p.Part_s {
					use = tr
				} else {
					if int64(tr.before) >= p.Part_s || int64(tr.after) >= p.Part_s {
						b1 := use.before - int(p.Part_s)
						a1 := use.after - int(p.Part_s)
						b2 := tr.before - int(p.Part_s)
						a2 := tr.after - int(p.Part_s)
						if (b1 < b2 && b1 < a2) || (a1 < b2 && a1 < a2) {
							c++
							continue
						}
						use = tr
					}
				}
			} else if Compare(string(mbr.Dsk_fit[0]), "W") {
				if int64(use.before) >= p.Part_s || int64(use.after) < p.Part_s {
					use = tr
				} else {
					if int64(tr.before) >= p.Part_s || int64(tr.after) >= p.Part_s {
						b1 := use.before - int(p.Part_s)
						a1 := use.after - int(p.Part_s)
						b2 := tr.before - int(p.Part_s)
						a2 := tr.after - int(p.Part_s)

						if (b1 > b2 && b1 > a2) || (a1 > b2 && a1 > a2) {
							c++
							continue
						}
						use = tr
					}
				}
			}
			c++
		}
		if use.before >= int(p.Part_s) || use.after >= int(p.Part_s) {
			if Compare(string(mbr.Dsk_fit[0]), "F") {
				if use.before >= int(p.Part_s) {
					p.Part_start = int64(use.start - use.before)
					startValue = int(p.Part_start)
				} else {
					p.Part_start = int64(use.end)
					startValue = int(p.Part_start)
				}
			} else if Compare(string(mbr.Dsk_fit[0]), "B") {
				b1 := use.before - int(p.Part_s)
				a1 := use.after - int(p.Part_s)

				if (use.before >= int(p.Part_s) && b1 < a1) || use.after < int(p.Part_start) {
					p.Part_start = int64(use.start - use.before)
					startValue = int(p.Part_start)
				} else {
					p.Part_start = int64(use.end)
					startValue = int(p.Part_start)
				}
			} else if Compare(string(mbr.Dsk_fit[0]), "W") {
				b1 := use.before - int(p.Part_s)
				a1 := use.after - int(p.Part_s)

				if (use.before >= int(p.Part_s) && b1 > a1) || use.after < int(p.Part_start) {
					p.Part_start = int64(use.start - use.before)
					startValue = int(p.Part_start)
				} else {
					p.Part_start = int64(use.end)
					startValue = int(p.Part_start)
				}
			}
			var partitions [4]Structs.Partition
			for i := 0; i < len(ps); i++ {
				partitions[i] = ps[i]
			}
			for i := 0; i < len(partitions); i++ {
				partition := partitions[i]
				if partition.Part_status != '1' {
					partitions[i] = p
					break
				}
			}
			mbr.Mbr_partitions_1 = partitions[0]
			mbr.Mbr_partitions_2 = partitions[1]
			mbr.Mbr_partitions_3 = partitions[2]
			mbr.Mbr_partitions_4 = partitions[3]
			return &mbr
		} else {
			Error("FDISK", "No hay espacio suficiente", responseString)
			return nil
		}
	}
}

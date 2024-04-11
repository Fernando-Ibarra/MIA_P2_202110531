package Comands

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"strings"
	"time"
	"unsafe"
)

func DataMkfs(context []string, responseString string) {
	id := ""
	tipo := "Full"
	fs := "2fs"

	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "id") {
			id = tk[1]
		} else if Compare(tk[0], "type") {
			if Compare(tk[1], "full") {
				tipo = tk[1]
			} else {
				Error("MKFS", "EL comando type debe tener calores especificos", responseString)
				return
			}
		} else if Compare(tk[0], "fs") {
			if Compare(tk[1], "2fs") || Compare(tk[1], "3fs") {
				fs = tk[1]
			} else {
				Error("MKFS", "El comando fs debe tener valores especificos", responseString)
			}
		}
	}
	if id == "" {
		Error("MKFS", "EL comando requiere el parámetro id obligatoriamente", responseString)
		return
	}
	mkfs(id, tipo, fs, responseString)
}

func mkfs(id string, t string, fs string, responseString string) {
	p := ""
	partition := GetMount("MKFS", id, &p, responseString)
	if Compare(fs, "2fs") {
		n := math.Floor(float64(partition.Part_s-int64(unsafe.Sizeof(Structs.SuperBlock{}))) / float64(4+unsafe.Sizeof(Structs.Inodos{})+3*unsafe.Sizeof(Structs.FilesBlocks{})))
		spr := Structs.NewSuperBlock()
		spr.S_magic = 0xEF53
		spr.S_inode_s = int64(unsafe.Sizeof(Structs.Inodos{}))
		spr.S_block_s = int64(unsafe.Sizeof(Structs.DirectoriesBlocks{}))
		spr.S_inodes_count = int64(n)
		spr.S_free_inodes_count = int64(n)
		spr.S_blocks_count = int64(3 * n)
		spr.S_free_blocks_count = int64(3 * n)
		dat := time.Now().String()
		copy(spr.S_mtime[:], dat)
		spr.S_mnt_count = spr.S_mnt_count + 1
		spr.S_filesystem_type = 2
		ext2(spr, partition, int64(n), p, responseString)
	} else if Compare(fs, "3fs") {
		n := math.Floor(float64(partition.Part_s-int64(unsafe.Sizeof(Structs.SuperBlock{}))) / float64(4+unsafe.Sizeof(Structs.Journaling{})+unsafe.Sizeof(Structs.Inodos{})+3*unsafe.Sizeof(Structs.FilesBlocks{})))
		spr := Structs.NewSuperBlock()
		spr.S_magic = 0xEF53
		spr.S_inode_s = int64(unsafe.Sizeof(Structs.Inodos{}))
		spr.S_block_s = int64(unsafe.Sizeof(Structs.Inodos{}))
		spr.S_inodes_count = int64(n)
		spr.S_free_inodes_count = int64(n)
		spr.S_blocks_count = int64(3 * n)
		spr.S_free_blocks_count = int64(3 * n)
		dat := time.Now().String()
		copy(spr.S_mtime[:], dat)
		spr.S_mnt_count = spr.S_mnt_count + 1
		spr.S_filesystem_type = 3
		ext3(spr, partition, int64(n), p, responseString)
	}
}

func ext2(spr Structs.SuperBlock, p Structs.Partition, n int64, path string, responseString string) {
	spr.S_bm_inode_start = p.Part_start + int64(unsafe.Sizeof(Structs.SuperBlock{}))
	spr.S_bm_block_start = spr.S_bm_inode_start + n
	spr.S_inode_start = spr.S_bm_block_start + (3 * n)
	spr.S_block_start = spr.S_bm_inode_start + (n * int64(unsafe.Sizeof(Structs.Inodos{})))

	file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco", responseString)
		return
	}
	file.Seek(p.Part_start, 0)
	var binary2 bytes.Buffer
	binary.Write(&binary2, binary.BigEndian, spr)
	WritingBytes(file, binary2.Bytes())

	zero := '0'
	file.Seek(spr.S_bm_inode_start, 0)
	for i := 0; i < int(n); i++ {
		var binaryZero bytes.Buffer
		binary.Write(&binaryZero, binary.BigEndian, zero)
		WritingBytes(file, binaryZero.Bytes())
	}

	file.Seek(spr.S_bm_block_start, 0)
	for i := 0; i < 3*int(n); i++ {
		var binaryZero bytes.Buffer
		binary.Write(&binaryZero, binary.BigEndian, zero)
		WritingBytes(file, binaryZero.Bytes())
	}

	inode := Structs.NewInodos()
	inode.I_uid = -1
	inode.I_gid = -1
	inode.I_s = -1
	for i := 0; i < len(inode.I_block); i++ {
		inode.I_block[i] = -1
	}
	inode.I_type = -1
	inode.I_perm = -1

	file.Seek(spr.S_inode_start, 0)
	for i := 0; i < int(n); i++ {
		var binaryInode bytes.Buffer
		binary.Write(&binaryInode, binary.BigEndian, inode)
		WritingBytes(file, binaryInode.Bytes())
	}

	folder := Structs.NewDirectoriesBlocks()

	for i := 0; i < len(folder.B_content); i++ {
		folder.B_content[i].B_inodo = -1
	}

	file.Seek(spr.S_block_start, 0)
	for i := 0; i < int(n); i++ {
		var binaryFolder bytes.Buffer
		binary.Write(&binaryFolder, binary.BigEndian, folder)
		WritingBytes(file, binaryFolder.Bytes())
	}
	file.Close()

	rescu := Structs.NewSuperBlock()
	file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco", responseString)
		return
	}
	file.Seek(p.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &rescu)
	if err_ != nil {
		Error("FDISK", "Error al leer el archivo", responseString)
		return
	}
	file.Close()

	inode.I_uid = 1
	inode.I_gid = 1
	inode.I_s = 0
	date := time.Now().String()
	copy(inode.I_atime[:], date)
	copy(inode.I_ctime[:], date)
	copy(inode.I_mtime[:], date)
	inode.I_type = 0
	inode.I_perm = 664
	inode.I_block[0] = 0

	fb := Structs.NewDirectoriesBlocks()
	copy(fb.B_content[0].B_name[:], ".")
	fb.B_content[0].B_inodo = 0
	copy(fb.B_content[1].B_name[:], "..")
	fb.B_content[1].B_inodo = 0
	copy(fb.B_content[2].B_name[:], "users.txt")
	fb.B_content[2].B_inodo = 1

	dataFile := "1,G,root\n1,U,root,root,123\n"
	inodetmp := Structs.NewInodos()
	inodetmp.I_uid = 1
	inodetmp.I_gid = 1
	inodetmp.I_s = int64(unsafe.Sizeof(dataFile) + unsafe.Sizeof(Structs.DirectoriesBlocks{}))

	copy(inodetmp.I_atime[:], date)
	copy(inodetmp.I_ctime[:], date)
	copy(inodetmp.I_mtime[:], date)
	inodetmp.I_type = 1
	inodetmp.I_perm = 664
	inodetmp.I_block[0] = 1

	inode.I_s = inodetmp.I_s + int64(unsafe.Sizeof(Structs.DirectoriesBlocks{})) + int64(unsafe.Sizeof(Structs.Inodos{}))

	var fileb Structs.FilesBlocks
	copy(fileb.B_content[:], dataFile)

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco", responseString)
		return
	}
	file.Seek(spr.S_bm_inode_start, 0)
	caracter := '1'

	var bin1 bytes.Buffer
	binary.Write(&bin1, binary.BigEndian, caracter)
	WritingBytes(file, bin1.Bytes())
	WritingBytes(file, bin1.Bytes())

	file.Seek(spr.S_bm_block_start, 0)
	var bin2 bytes.Buffer
	binary.Write(&bin2, binary.BigEndian, caracter)
	WritingBytes(file, bin2.Bytes())
	WritingBytes(file, bin1.Bytes())

	file.Seek(spr.S_inode_start, 0)
	var bin3 bytes.Buffer
	binary.Write(&bin3, binary.BigEndian, inode)
	WritingBytes(file, bin3.Bytes())

	file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var bin4 bytes.Buffer
	binary.Write(&bin4, binary.BigEndian, inodetmp)
	WritingBytes(file, bin4.Bytes())

	file.Seek(spr.S_block_start, 0)
	var bin5 bytes.Buffer
	binary.Write(&bin5, binary.BigEndian, fb)
	WritingBytes(file, bin5.Bytes())

	file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{})), 0)
	var bin6 bytes.Buffer
	binary.Write(&bin6, binary.BigEndian, fileb)
	WritingBytes(file, bin6.Bytes())

	file.Close()

	partitionName := ""
	for i := 0; i < len(p.Part_name); i++ {
		if p.Part_name[i] != 0 {
			partitionName += string(p.Part_name[i])
		}
	}
	Message("MKFS", "Se ha formateado la partición "+partitionName+" correctamente", responseString)
}

func ext3(spr Structs.SuperBlock, p Structs.Partition, n int64, path string, responseString string) {
	spr.S_bm_inode_start = p.Part_start + int64(unsafe.Sizeof(Structs.SuperBlock{})) + (n * int64(unsafe.Sizeof(Structs.Journaling{})))
	spr.S_bm_block_start = spr.S_bm_inode_start + n
	spr.S_inode_start = spr.S_bm_block_start + (3 * n)
	spr.S_block_start = spr.S_bm_inode_start + (n * int64(unsafe.Sizeof(Structs.Inodos{})))

	file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco", responseString)
		return
	}

	file.Seek(p.Part_start, 0)
	var binary2 bytes.Buffer
	binary.Write(&binary2, binary.BigEndian, spr)
	WritingBytes(file, binary2.Bytes())

	zero := '0'
	file.Seek(spr.S_bm_inode_start, 0)
	for i := 0; i < int(n); i++ {
		var binaryZero bytes.Buffer
		binary.Write(&binaryZero, binary.BigEndian, zero)
		WritingBytes(file, binaryZero.Bytes())
	}

	file.Seek(spr.S_bm_block_start, 0)
	for i := 0; i < 3*int(n); i++ {
		var binaryZero bytes.Buffer
		binary.Write(&binaryZero, binary.BigEndian, zero)
		WritingBytes(file, binaryZero.Bytes())
	}

	jour := Structs.NewJournaling()
	file.Seek(p.Part_start+int64(unsafe.Sizeof(Structs.SuperBlock{})), 0)
	for i := 0; i < int(n); i++ {
		var binaryZero bytes.Buffer
		binary.Write(&binaryZero, binary.BigEndian, jour)
		WritingBytes(file, binaryZero.Bytes())
	}

	inode := Structs.NewInodos()
	inode.I_uid = -1
	inode.I_gid = -1
	inode.I_s = -1
	for i := 0; i < len(inode.I_block); i++ {
		inode.I_block[i] = -1
	}
	inode.I_type = -1
	inode.I_perm = -1

	file.Seek(spr.S_inode_start, 0)
	for i := 0; i < int(n); i++ {
		var binaryInode bytes.Buffer
		binary.Write(&binaryInode, binary.BigEndian, inode)
		WritingBytes(file, binaryInode.Bytes())
	}

	folder := Structs.NewDirectoriesBlocks()

	for i := 0; i < len(folder.B_content); i++ {
		folder.B_content[i].B_inodo = -1
	}

	file.Seek(spr.S_block_start, 0)
	for i := 0; i < int(n); i++ {
		var binaryFolder bytes.Buffer
		binary.Write(&binaryFolder, binary.BigEndian, folder)
		WritingBytes(file, binaryFolder.Bytes())
	}
	file.Close()

	rescu := Structs.NewSuperBlock()

	file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco", responseString)
		return
	}

	file.Seek(p.Part_start, 0)
	data := readBytes(file, int(unsafe.Sizeof(Structs.SuperBlock{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &rescu)
	if err_ != nil {
		Error("FDISK", "Error al leer el archivo", responseString)
		return
	}
	file.Close()

	inode.I_uid = 1
	inode.I_gid = 1
	inode.I_s = 0
	date := time.Now().String()
	copy(inode.I_atime[:], date)
	copy(inode.I_ctime[:], date)
	copy(inode.I_mtime[:], date)
	inode.I_type = 0
	inode.I_perm = 664
	inode.I_block[0] = 0

	fb := Structs.NewDirectoriesBlocks()
	copy(fb.B_content[0].B_name[:], ".")
	fb.B_content[0].B_inodo = 0
	copy(fb.B_content[1].B_name[:], "..")
	fb.B_content[1].B_inodo = 0
	copy(fb.B_content[2].B_name[:], "users.txt")
	fb.B_content[2].B_inodo = 1

	dataFile := "1,G,root\n1,U,root,root,123\n"
	inodetmp := Structs.NewInodos()
	inodetmp.I_uid = 1
	inodetmp.I_gid = 1
	inodetmp.I_s = int64(unsafe.Sizeof(dataFile) + unsafe.Sizeof(Structs.DirectoriesBlocks{}))

	copy(inodetmp.I_atime[:], date)
	copy(inodetmp.I_ctime[:], date)
	copy(inodetmp.I_mtime[:], date)
	inodetmp.I_type = 1
	inodetmp.I_perm = 664
	inodetmp.I_block[0] = 1

	inode.I_s = inodetmp.I_s + int64(unsafe.Sizeof(Structs.DirectoriesBlocks{})) + int64(unsafe.Sizeof(Structs.Inodos{}))

	var fileb Structs.FilesBlocks
	copy(fileb.B_content[:], dataFile)

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco", responseString)
		return
	}

	journalingR := Structs.NewJournaling()
	operation := "mkfs"
	pathR := "/"
	contentR := "-"
	dateR := time.Now().String()
	copy(journalingR.Operation[:], operation)
	copy(journalingR.Path[:], pathR)
	copy(journalingR.Content[:], contentR)
	copy(journalingR.Date[:], dateR)

	journalingUser := Structs.NewJournaling()
	operation = "mkfs"
	pathU := "users.txt"
	contentU := dataFile
	dateU := time.Now().String()
	copy(journalingUser.Operation[:], operation)
	copy(journalingUser.Path[:], pathU)
	copy(journalingUser.Content[:], contentU)
	copy(journalingUser.Date[:], dateU)

	journalingRoot := p.Part_start + int64(unsafe.Sizeof(Structs.SuperBlock{}))
	file.Seek(journalingRoot, 0)
	var binJR bytes.Buffer
	binary.Write(&binJR, binary.BigEndian, journalingR)
	WritingBytes(file, binJR.Bytes())

	journalingU := p.Part_start + int64(unsafe.Sizeof(Structs.SuperBlock{})) + int64(unsafe.Sizeof(Structs.Journaling{}))
	file.Seek(journalingU, 0)
	var binJu bytes.Buffer
	binary.Write(&binJu, binary.BigEndian, journalingUser)
	WritingBytes(file, binJu.Bytes())

	caracter := '1'
	file.Seek(spr.S_bm_inode_start, 0)
	var bin1 bytes.Buffer
	binary.Write(&bin1, binary.BigEndian, caracter)
	WritingBytes(file, bin1.Bytes())
	WritingBytes(file, bin1.Bytes())

	file.Seek(spr.S_bm_block_start, 0)
	var bin2 bytes.Buffer
	binary.Write(&bin2, binary.BigEndian, caracter)
	WritingBytes(file, bin2.Bytes())
	WritingBytes(file, bin1.Bytes())

	file.Seek(spr.S_inode_start, 0)
	var bin3 bytes.Buffer
	binary.Write(&bin3, binary.BigEndian, inode)
	WritingBytes(file, bin3.Bytes())

	file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var bin4 bytes.Buffer
	binary.Write(&bin4, binary.BigEndian, inodetmp)
	WritingBytes(file, bin4.Bytes())

	file.Seek(spr.S_block_start, 0)
	var bin5 bytes.Buffer
	binary.Write(&bin5, binary.BigEndian, fb)
	WritingBytes(file, bin5.Bytes())

	file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.DirectoriesBlocks{})), 0)
	var bin6 bytes.Buffer
	binary.Write(&bin6, binary.BigEndian, fileb)
	WritingBytes(file, bin6.Bytes())

	partitionName := ""
	for i := 0; i < len(p.Part_name); i++ {
		if p.Part_name[i] != 0 {
			partitionName += string(p.Part_name[i])
		}
	}

	file.Close()
	Message("MKFS", "Se ha formateado la partición "+partitionName+" correctamente", responseString)
}

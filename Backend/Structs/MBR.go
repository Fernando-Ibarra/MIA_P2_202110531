package Structs

type MBR struct {
	Mbr_tamano         int64
	Mbr_fecha_creacion [16]byte
	Mbr_dsk_signature  int64
	Dsk_fit            [1]byte
	Mbr_partitions_1   Partition
	Mbr_partitions_2   Partition
	Mbr_partitions_3   Partition
	Mbr_partitions_4   Partition
}

func NewMBR() MBR {
	var mb MBR
	return mb
}

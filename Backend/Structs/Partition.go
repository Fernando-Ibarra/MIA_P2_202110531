package Structs

type Partition struct {
	Part_status      byte
	Part_type        byte
	Part_fit         byte
	Part_start       int64
	Part_s           int64
	Part_name        [16]byte
	Part_correlative int64
	Part_id          [4]byte
}

func NewPartition() Partition {
	var Part Partition
	Part.Part_status = '0'
	Part.Part_type = 'P'
	Part.Part_fit = 'F'
	Part.Part_start = -1
	Part.Part_s = 0
	Part.Part_name = [16]byte{}
	Part.Part_id = [4]byte{}
	return Part
}

package Structs

import "unsafe"

type SuperBlock struct {
	S_filesystem_type   int64
	S_inodes_count      int64
	S_blocks_count      int64
	S_free_blocks_count int64
	S_free_inodes_count int64
	S_mtime             [16]byte
	S_umtime            int64
	S_mnt_count         int64
	S_magic             int64
	S_inode_s           int64
	S_block_s           int64
	S_firts_ino         int64
	S_firts_blo         int64
	S_bm_inode_start    int64
	S_bm_block_start    int64
	S_inode_start       int64
	S_block_start       int64
}

func NewSuperBlock() SuperBlock {
	var spb SuperBlock
	spb.S_magic = 0xEF53
	spb.S_inode_s = int64(unsafe.Sizeof(Inodos{}))
	spb.S_block_s = int64(unsafe.Sizeof(DirectoriesBlocks{}))
	spb.S_firts_ino = 0
	spb.S_firts_blo = 0
	return spb
}

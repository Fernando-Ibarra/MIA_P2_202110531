package Structs

type DirectoriesBlocks struct {
	B_content [4]Content
}

func NewDirectoriesBlocks() DirectoriesBlocks {
	var db DirectoriesBlocks
	for i := 0; i < len(db.B_content); i++ {
		db.B_content[i] = NewContent()
	}
	return db
}

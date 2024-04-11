package Structs

type Journaling struct {
	Operation [15]byte
	Path      [50]byte
	Content   [50]byte
	Date      [16]byte
}

func NewJournaling() Journaling {
	var journaling Journaling
	copy(journaling.Path[:], "-")
	return journaling
}

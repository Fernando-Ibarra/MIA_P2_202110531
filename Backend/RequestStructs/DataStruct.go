package RequestStructs

type DataMkDisk struct {
	Size         string `json:"size"`
	Fit          string `json:"fit"`
	Unit         string `json:"unit"`
	CounterDisks int    `json:"counterDisks"`
}

type DataRmDisk struct {
	Driveletter string `json:"driveletter"`
}

type DataFDisk struct {
	Size        string `json:"size"`
	Driveletter string `json:"driveletter"`
	Name        string `json:"name"`
	Unit        string `json:"unit"`
	Tipo        string `json:"tipo"`
	Fit         string `json:"fit"`
	Delete      string `json:"delete"`
}

type DataMount struct {
	Driveletter string `json:"driveletter"`
	Name        string `json:"name"`
}

type DataUnMount struct {
	Id string `json:"id"`
}

type DataMkfs struct {
	Id   string `json:"id"`
	Tipo string `json:"tipo"`
	Fs   string `json:"fs"`
}

type DataLogin struct {
	Id   string `json:"id"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type DataMkUsr struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Grp  string `json:"grp"`
}

type DataRmUsr struct {
	User string `json:"user"`
}

type DataChGrp struct {
	User string `json:"user"`
	Grp  string `json:"grp"`
}

type DataMkGrp struct {
	Name string `json:"name"`
}

type DataRmGrp struct {
	Name string `json:"name"`
}

type DataMkFile struct {
	Path     string `json:"path"`
	Rboolean bool   `json:"rboolean"`
	Size     string `json:"size"`
	Cont     string `json:"cont"`
}

type DataCat struct {
	Paths []string `json:"paths"`
}

type DataRename struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type DataMkDir struct {
	Path     string `json:"path"`
	Rboolean bool   `json:"rboolean"`
}

type DataChOwn struct {
	Path     string `json:"path"`
	User     string `json:"user"`
	Rboolean bool   `json:"rboolean"`
}

type DataChMod struct {
	Path     string `json:"path"`
	Ugo      string `json:"ugo"`
	Rboolean bool   `json:"rboolean"`
}

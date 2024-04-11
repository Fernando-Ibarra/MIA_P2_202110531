package Comands

import (
	"strconv"
	"strings"
)

func DataUnMount(tokens []string, responseString string) {
	if len(tokens) > 1 {
		Error("UNMOUNT", "Solo se acepta el párametro id", responseString)
		return
	}
	id := ""
	error_ := false
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Compare(tk[0], "id") {
			if id == "" {
				id = tk[1]
			} else {
				Error("UNMOUNT", "Parámetro driveletter repetido en el comando: "+tk[0], responseString)
			}
		} else {
			Error("UNMOUNT", "No se esperaba el parámetro "+tk[0], responseString)
			error_ = false
			return
		}
	}
	if error_ {
		return
	}
	if id == "" {
		Error("UNMOUNT", "Se require el parámetro id", responseString)
		return
	} else {
		unmount(id, responseString)
	}
}

func unmount(id string, responseString string) {
	if !(id[2] == '3' && id[3] == '1') {
		Error("UNMOUNT", "El primer identificador no es válido", responseString)
		return
	}
	letter := id[0]
	j, _ := strconv.Atoi(string(id[1] - 1))
	if j < 0 {
		Error("UNMOUNT", "El primer identificador no es válido", responseString)
		return
	}
	for i := 0; i < 99; i++ {
		if DiskMount[i].Partitions[j].State == 1 {
			if DiskMount[i].Partitions[j].Letter == letter {
				DiskMount[i].Partitions[j].State = 0
				Message("UNMOUNT", "Se ha realizado correctamente el unmount -id="+id, responseString)
				return
			} else {
				Error("UNMOUNT", "No se ha podido realizar correctamente el unmount -id="+id, responseString)
				return
			}
		}
	}

}

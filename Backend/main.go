package main

import (
	"Backend/Comands"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"strings"
)

var CounterDisk = 1

var logued = false

var ResponseString = ""

type DataReq struct {
	ComandsReq string `json:"comands_req"`
}

func main() {
	// ROUTER
	// router := mux.NewRouter()
	mux := http.NewServeMux()

	// ROUTES
	mux.HandleFunc("/", initServer)
	mux.HandleFunc("/makeMagic", makeMagic)
	mux.HandleFunc("/file-system", makeFileSystem)
	mux.HandleFunc("/reports", makeReports)

	// CORS allowCORS(router)
	handler := cors.Default().Handler(mux)

	// start server listen

	// SERVER
	fmt.Println("Server on port 3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
}

func allowCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		handler.ServeHTTP(w, r)
	})
}

func initServer(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "<h1>Servidor Activo</h1>")
	if err != nil {
		return
	}
}

func makeMagic(w http.ResponseWriter, r *http.Request) {
	CounterDisk = 1

	var data DataReq
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// get current comands
	com := data.ComandsReq

	// current path
	currentPath, _ := os.Getwd()
	path := currentPath + "/comandos.txt"

	// crate file
	var _, errF = os.Stat(path)
	if os.IsNotExist(errF) {
		var file, err2 = os.Create(path)
		if err2 != nil {
			return
		}
		defer file.Close()
	}

	// write file
	var file, err3 = os.OpenFile(path, os.O_RDWR, 0644)
	if err3 != nil {
		return
	}
	defer file.Close()
	_, err3 = file.WriteString(com)
	if err3 != nil {
		return
	}
	err3 = file.Sync()
	if err3 != nil {
		return
	}

	// read file to execute comands
	fileRead, errRead := os.Open(path)
	if errRead != nil {
		log.Fatalf("Error al abrir el archivo: %s", err)
	}
	fileScanner := bufio.NewScanner(fileRead)
	for fileScanner.Scan() {
		text := fileScanner.Text()
		text = strings.TrimSpace(text)
		tk := Comand(text)
		if text != "" {
			if Comands.Compare(tk, "pause") {
				ResponseString += ">>>>>>>>>> COMANDD PAUSA <<<<<<<<<<<<<<<<<<<<\n"
				Comands.Message("PAUSE", "PROGRAMDA PAUSADO", ResponseString)
				continue
			} else if string(text[0]) == "#" {
				ResponseString += ">>>>>>>>>> COMENTARIO <<<<<<<<<<<<<<<<<<<<\n"
				Comands.Message("COMENTARIO", text, ResponseString)
				continue
			}
			text = strings.TrimLeft(text, tk)
			tokens := Separatorokens(text)
			functions(tk, tokens)
		}
	}
	if err = fileScanner.Err(); err != nil {
		log.Fatalf(
			"Error al leer el archivo: %s",
			err,
		)
	}
	defer r.Body.Close()
	fmt.Fprint(w, ResponseString)
}

func makeFileSystem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Comands.MakeJson()
	defer r.Body.Close()
	_, err := fmt.Fprintf(w, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func makeReports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Comands.MakeRep()
	defer r.Body.Close()
	_, err := fmt.Fprintf(w, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Comand(text string) string {
	var tkn string
	finished := false
	for i := 0; i < len(text); i++ {
		if finished {
			if string(text[i]) == " " || string(text[i]) == "-" {
				break
			}
			tkn += string(text[i])
		} else if string(text[i]) != " " && !finished {
			if string(text[i]) == "#" {
				tkn = text
			} else {
				tkn += string(text[i])
				finished = true
			}
		}
	}
	return tkn
}

func Separatorokens(text string) []string {
	var tokens []string
	if text == "" {
		return tokens
	}
	text += " "
	var token string
	state := 0
	for i := 0; i < len(text); i++ {
		c := string(text[i])
		if state == 0 && c == "-" {
			state = 1
		} else if state == 0 && c == "#" {
			continue
		} else if state != 0 {
			if state == 1 {
				if c == "=" {
					state = 2
				} else if c == " " {
					continue
				} else if (c == "P" || c == "p") && string(text[i+1]) == " " && string(text[i-1]) == "-" {
					state = 0
					tokens = append(tokens, c)
					token = ""
					continue
				} else if (c == "R" || c == "r") && string(text[i+1]) == " " && string(text[i-1]) == "-" {
					state = 0
					tokens = append(tokens, c)
					token = ""
					continue
				}
			} else if state == 2 {
				if c == " " {
					continue
				}
				if c == "\"" {
					state = 3
					continue
				} else {
					state = 4
				}
			} else if state == 3 {
				if c == "\"" {
					state = 4
					continue
				}
			} else if state == 4 && c == "\"" {
				tokens = []string{}
				continue
			} else if state == 4 && c == " " {
				state = 0
				tokens = append(tokens, token)
				token = ""
				continue
			}
			token += c
		}
	}
	return tokens
}

func functions(token string, tks []string) {
	if token != "" {
		if Comands.Compare(token, "MKDISK") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO MKDISK <-------------------\n"
			Comands.DataMKDISK(tks, CounterDisk, &CounterDisk, ResponseString)
			CounterDisk++
		} else if Comands.Compare(token, "RMDISK") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO RMDISK <-------------------\n"
			Comands.RMDISK(tks, ResponseString)
		} else if Comands.Compare(token, "FDISK") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO FDISK <-------------------\n"
			Comands.DataFDISK(tks, ResponseString)
		} else if Comands.Compare(token, "MOUNT") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO MOUNT <-------------------\n"
			Comands.DataMount(tks, ResponseString)
		} else if Comands.Compare(token, "UNMOUNT") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO UNMOUNT <-------------------\n"
			Comands.DataUnMount(tks, ResponseString)
		} else if Comands.Compare(token, "MKFS") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO MKFS <-------------------\n"
			Comands.DataMkfs(tks, ResponseString)
		} else if Comands.Compare(token, "LOGIN") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO LOGIN <-------------------\n"
			if logued {
				Comands.Error("LOGIN", "Ya hay un usuario en linea.", ResponseString)
				return
			} else {
				logued = Comands.DataUserLogin(tks, ResponseString)
			}
		} else if Comands.Compare(token, "LOGOUT") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO LOGOUT <-------------------\n"
			if !logued {
				Comands.Error("LOGOUT", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				logued = Comands.LogOut(ResponseString)
			}
		} else if Comands.Compare(token, "MKGRP") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO MKGRP <-------------------\n"
			if !logued {
				Comands.Error("MKGRP", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				Comands.DataGroup(tks, "MK", ResponseString)
			}
		} else if Comands.Compare(token, "RMGRP") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO RMGRP <-------------------\n"
			if !logued {
				Comands.Error("RMGRP", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				Comands.DataGroup(tks, "RM", ResponseString)
			}
		} else if Comands.Compare(token, "CHGRP") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO CHGRP <-------------------\n"
			if !logued {
				Comands.Error("CHGRP", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				Comands.DataUser(tks, "CH", ResponseString)
			}
		} else if Comands.Compare(token, "MKUSR") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO MKUSR <-------------------\n"
			if !logued {
				Comands.Error("MKUSER", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				Comands.DataUser(tks, "MK", ResponseString)
			}
		} else if Comands.Compare(token, "RMUSR") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO RMUSR <-------------------\n"
			if !logued {
				Comands.Error("RMUSER", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				Comands.DataUser(tks, "RM", ResponseString)
			}
		} else if Comands.Compare(token, "MKDIR") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO MKDIR <-------------------\n"
			if !logued {
				Comands.Error("MKDIR", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("MKDIR", Comands.Logged.Id, &p, ResponseString)
				Comands.DataDir(tks, partition, p, ResponseString)
			}
		} else if Comands.Compare(token, "MKFILE") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO MKFILE <-------------------\n"
			if !logued {
				Comands.Error("MKDIR", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("MKDIR", Comands.Logged.Id, &p, ResponseString)
				Comands.DataFile(tks, partition, p, ResponseString)
			}
		} else if Comands.Compare(token, "CAT") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO CAT <-------------------\n"
			if !logued {
				Comands.Error("CAT", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("CAT", Comands.Logged.Id, &p, ResponseString)
				Comands.DataCat(tks, partition, p, ResponseString)
			}
		} else if Comands.Compare(token, "CHMOD") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO CHMOD <-------------------\n"
			if !logued {
				Comands.Error("CHMOD", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("CHMOD", Comands.Logged.Id, &p, ResponseString)
				Comands.DataChmod(tks, partition, p, ResponseString)
			}
		} else if Comands.Compare(token, "CHOWN") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO CHOWN <-------------------\n"
			if !logued {
				Comands.Error("CHOWN", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("CHOWN", Comands.Logged.Id, &p, ResponseString)
				Comands.DataChown(tks, partition, p, ResponseString)
			}
		} else if Comands.Compare(token, "RENAME") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO RENAME <-------------------\n"
			if !logued {
				Comands.Error("RENAME", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("CHOWN", Comands.Logged.Id, &p, ResponseString)
				Comands.DataRename(tks, partition, p, ResponseString)
			}
		} else if Comands.Compare(token, "MOVE") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO MOVE <-------------------\n"
			if !logued {
				Comands.Error("MOVE", "Aún no se ha iniciado sesión", ResponseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("MOVE", Comands.Logged.Id, &p, ResponseString)
				Comands.DataMove(tks, partition, p, ResponseString)
			}
		} else if Comands.Compare(token, "REP") {
			ResponseString += ""
			ResponseString += "-------------------> COMANDO REP <-------------------\n"
			Comands.DataRep(tks, ResponseString)
		} else {
			Comands.Error("ANALIZADOR", "NO se reconoce el comando \" "+token+"\" ", ResponseString)
		}
	}
}

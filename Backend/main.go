package main

import (
	"Backend/Comands"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
)

var CounterDisk = 1

var logued = false

var responseString = ""

type DataReq struct {
	ComandsReq string `json:"comands_req"`
}

func main() {
	// ROUTER
	router := mux.NewRouter()

	// ROUTES
	router.HandleFunc("/", initServer).Methods("GET")
	router.HandleFunc("/makeMagic", makeMagic).Methods("POST")

	// CORS
	handler := allowCORS(router)

	// SERVER
	fmt.Println("Server on port 3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
}

func allowCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
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
	var datar DataReq
	err := json.NewDecoder(r.Body).Decode(&datar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get current comands
	com := datar.ComandsReq

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
				responseString += ">>>>>>>>>> COMANDD PAUSA <<<<<<<<<<<<<<<<<<<<\n"
				Comands.Message("PAUSE", "PROGRAMDA PAUSADO", responseString)
				continue
			} else if string(text[0]) == "#" {
				responseString += ">>>>>>>>>> COMENTARIO <<<<<<<<<<<<<<<<<<<<\n"
				Comands.Message("COMENTARIO", text, responseString)
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
	fmt.Fprint(w, responseString)
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
			responseString += ""
			responseString += "-------------------> COMANDO MKDISK <-------------------\n"
			Comands.DataMKDISK(tks, CounterDisk, &CounterDisk, responseString)
			CounterDisk++
		} else if Comands.Compare(token, "RMDISK") {
			responseString += ""
			responseString += "-------------------> COMANDO RMDISK <-------------------\n"
			Comands.RMDISK(tks, responseString)
		} else if Comands.Compare(token, "FDISK") {
			responseString += ""
			responseString += "-------------------> COMANDO FDISK <-------------------\n"
			Comands.DataFDISK(tks, responseString)
		} else if Comands.Compare(token, "MOUNT") {
			responseString += ""
			responseString += "-------------------> COMANDO MOUNT <-------------------\n"
			Comands.DataMount(tks, responseString)
		} else if Comands.Compare(token, "UNMOUNT") {
			responseString += ""
			responseString += "-------------------> COMANDO UNMOUNT <-------------------\n"
			Comands.DataUnMount(tks, responseString)
		} else if Comands.Compare(token, "MKFS") {
			responseString += ""
			responseString += "-------------------> COMANDO MKFS <-------------------\n"
			Comands.DataMkfs(tks, responseString)
		} else if Comands.Compare(token, "LOGIN") {
			responseString += ""
			responseString += "-------------------> COMANDO LOGIN <-------------------\n"
			if logued {
				Comands.Error("LOGIN", "Ya hay un usuario en linea.", responseString)
				return
			} else {
				logued = Comands.DataUserLogin(tks, responseString)
			}
		} else if Comands.Compare(token, "LOGOUT") {
			responseString += ""
			responseString += "-------------------> COMANDO LOGOUT <-------------------\n"
			if !logued {
				Comands.Error("LOGOUT", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				logued = Comands.LogOut(responseString)
			}
		} else if Comands.Compare(token, "MKGRP") {
			responseString += ""
			responseString += "-------------------> COMANDO MKGRP <-------------------\n"
			if !logued {
				Comands.Error("MKGRP", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				Comands.DataGroup(tks, "MK", responseString)
			}
		} else if Comands.Compare(token, "RMGRP") {
			responseString += ""
			responseString += "-------------------> COMANDO RMGRP <-------------------\n"
			if !logued {
				Comands.Error("RMGRP", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				Comands.DataGroup(tks, "RM", responseString)
			}
		} else if Comands.Compare(token, "CHGRP") {
			responseString += ""
			responseString += "-------------------> COMANDO CHGRP <-------------------\n"
			if !logued {
				Comands.Error("CHGRP", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				Comands.DataUser(tks, "CH", responseString)
			}
		} else if Comands.Compare(token, "MKUSR") {
			responseString += ""
			responseString += "-------------------> COMANDO MKUSR <-------------------\n"
			if !logued {
				Comands.Error("MKUSER", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				Comands.DataUser(tks, "MK", responseString)
			}
		} else if Comands.Compare(token, "RMUSR") {
			responseString += ""
			responseString += "-------------------> COMANDO RMUSR <-------------------\n"
			if !logued {
				Comands.Error("RMUSER", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				Comands.DataUser(tks, "RM", responseString)
			}
		} else if Comands.Compare(token, "MKDIR") {
			responseString += ""
			responseString += "-------------------> COMANDO MKDIR <-------------------\n"
			if !logued {
				Comands.Error("MKDIR", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("MKDIR", Comands.Logged.Id, &p, responseString)
				Comands.DataDir(tks, partition, p, responseString)
			}
		} else if Comands.Compare(token, "MKFILE") {
			responseString += ""
			responseString += "-------------------> COMANDO MKFILE <-------------------\n"
			if !logued {
				Comands.Error("MKDIR", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("MKDIR", Comands.Logged.Id, &p, responseString)
				Comands.DataFile(tks, partition, p, responseString)
			}
		} else if Comands.Compare(token, "CAT") {
			responseString += ""
			responseString += "-------------------> COMANDO CAT <-------------------\n"
			if !logued {
				Comands.Error("CAT", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("CAT", Comands.Logged.Id, &p, responseString)
				Comands.DataCat(tks, partition, p, responseString)
			}
		} else if Comands.Compare(token, "CHMOD") {
			responseString += ""
			responseString += "-------------------> COMANDO CHMOD <-------------------\n"
			if !logued {
				Comands.Error("CHMOD", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("CHMOD", Comands.Logged.Id, &p, responseString)
				Comands.DataChmod(tks, partition, p, responseString)
			}
		} else if Comands.Compare(token, "CHOWN") {
			responseString += ""
			responseString += "-------------------> COMANDO CHOWN <-------------------\n"
			if !logued {
				Comands.Error("CHOWN", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("CHOWN", Comands.Logged.Id, &p, responseString)
				Comands.DataChown(tks, partition, p, responseString)
			}
		} else if Comands.Compare(token, "RENAME") {
			responseString += ""
			responseString += "-------------------> COMANDO RENAME <-------------------\n"
			if !logued {
				Comands.Error("RENAME", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("CHOWN", Comands.Logged.Id, &p, responseString)
				Comands.DataRename(tks, partition, p, responseString)
			}
		} else if Comands.Compare(token, "MOVE") {
			responseString += ""
			responseString += "-------------------> COMANDO MOVE <-------------------\n"
			if !logued {
				Comands.Error("MOVE", "Aún no se ha iniciado sesión", responseString)
				return
			} else {
				var p string
				partition := Comands.GetMount("MOVE", Comands.Logged.Id, &p, responseString)
				Comands.DataMove(tks, partition, p, responseString)
			}
		} else if Comands.Compare(token, "REP") {
			responseString += ""
			responseString += "-------------------> COMANDO REP <-------------------\n"
			Comands.DataRep(tks, responseString)
		} else {
			Comands.Error("ANALIZADOR", "NO se reconoce el comando \" "+token+"\" ", responseString)
		}
	}
}

func setResponseString(text string) {
	responseString += text + "\n"
}

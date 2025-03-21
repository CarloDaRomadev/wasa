package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Funzione che serve a creare un gruppo dato un nome e una foto
func (rt *_router) CreaGruppo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Nome         string `json:"nome"`
		PercorsoFoto string `json:"foto"`
	}
	UtenteChiamante := ps.ByName("utente")

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	if len(input.Nome) == 0 {
		http.Error(w, "Il nome è obbligatorio", http.StatusBadRequest)
		return
	}
	if len(input.PercorsoFoto) == 0 {
		http.Error(w, "La foto è obbligatoria", http.StatusBadRequest)
		return
	}

	var fileFoto []byte

	if len(input.PercorsoFoto) > 0 {
		fileFoto, err = ReadImageFile(input.PercorsoFoto)
		if err != nil {
			http.Error(w, "Errore durante la lettura della foto: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	idFoto, err := rt.db.CreaFoto(input.PercorsoFoto, fileFoto)
	if err != nil {
		http.Error(w, "Errore durante l'inserimento della foto profilo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = rt.db.CreaGruppoDB(UtenteChiamante, input.Nome, idFoto)
	if err != nil {
		http.Error(w, "Errore durante la creazione dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Gruppo creato con successo "))
}

// test
/*
curl -X POST http://localhost:3000/wasachat/:utente/gruppi \
-H "Content-Type: application/json" \
-d '{
  "nome": "Gruppo1",
  "foto": "/home/carlo/Scrivania/wasachat/immagini/prova.png"
}'
*/

// funzione che serve ad aggiungere un utente ad un gruppo, se l'utente è già presente nel gruppo non verrà aggiunto
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Utente string `json:"utente"`
	}
	UtenteChiamante := ps.ByName("utente")
	idConversazioneStr := ps.ByName("chat")

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	if len(input.Utente) == 0 {
		http.Error(w, "Il nome è obbligatorio", http.StatusBadRequest)
		return
	}

	idConversazione, err := strconv.Atoi(idConversazioneStr)
	if err != nil {
		http.Error(w, "ID della conversazione non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.AggiungiAGruppoDB(idConversazione, UtenteChiamante, input.Utente)
	if err != nil {
		http.Error(w, "Errore durante l'aggiunta dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Utente aggiunto con successo "))
}

// test
/*
curl -X PUT http://localhost:3000/wasachat/:utente/chats/gruppi/:chat/aggiungi \
-H "Content-Type: application/json" \
-d '{
  "utente": "Luigi"
}'
*/

// Funzione che serve a lasciare un gruppo, l'utente deve essere presente nel gruppo
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	idConversazioneStr := ps.ByName("chat")
	idConversazione, err := strconv.Atoi(idConversazioneStr)
	if err != nil {
		http.Error(w, "ID della conversazione non valido", http.StatusBadRequest)
		return
	}
	err = rt.db.LasciaGruppo(idConversazione, UtenteChiamante)
	if err != nil {
		http.Error(w, "Errore durante la rimozione dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Utente rimosso con successo "))
}

// Funzione per impostare una nuova foto al gruppo
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idConversazioneStr := ps.ByName("chat")
	var input struct {
		PercorsoFoto string `json:"foto"`
	}
	UtenteChiamante := ps.ByName("utente")
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}
	if len(input.PercorsoFoto) == 0 {
		http.Error(w, "La foto è obbligatoria", http.StatusBadRequest)
		return
	}

	var fileFoto []byte

	if len(input.PercorsoFoto) > 0 {
		fileFoto, err = ReadImageFile(input.PercorsoFoto)
		if err != nil {
			http.Error(w, "Errore durante la lettura della foto: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	idFoto, err := rt.db.CreaFoto(input.PercorsoFoto, fileFoto)
	if err != nil {
		http.Error(w, "Errore durante l'inserimento della foto profilo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	idConversazione, err := strconv.Atoi(idConversazioneStr)
	if err != nil {
		http.Error(w, "ID della conversazione non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.ImpostaFotoGruppo(UtenteChiamante, idFoto, idConversazione)
	if err != nil {
		http.Error(w, "Errore durante la creazione dell'utente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Gruppo modificato con successo "))
}

// Funzione per impostare un nome ad un gruppo
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	UtenteChiamante := ps.ByName("utente")
	idConversazioneStr := ps.ByName("chat")

	var input struct {
		Nome string `json:"nome"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Formato della richiesta non valido", http.StatusBadRequest)
		return
	}

	idConversazione, err := strconv.Atoi(idConversazioneStr)
	if err != nil {
		http.Error(w, "ID della conversazione non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.ImpostaNomeGruppo(UtenteChiamante, input.Nome, idConversazione)
	if err != nil {
		http.Error(w, "Errore durante l'aggiornamento del nome del gruppo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Gruppo modificato con successo "))
}

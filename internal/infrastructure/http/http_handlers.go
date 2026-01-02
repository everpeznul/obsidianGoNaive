package http

import (
	"encoding/json"
	"net/http"
	"obsidianGoNaive/internal/domain"
	"obsidianGoNaive/internal/use_case"

	"github.com/google/uuid"
)

// /
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	obsiLog.Debug("Home handler")
	w.Write([]byte("Home page"))
}

// /notes/{id}
func NotesUUIDHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	tempId := r.PathValue("id")
	id, err := uuid.Parse(tempId)

	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)

		obsiLog.Debug("NotesUUIDHandler error", "error", "invalid uuid", "id", id)
		return
	}

	switch r.Method {

	case http.MethodGet:
		note, err := domain.Repo.GetByID(ctx, id)

		if writeError(w, err) {
			return
		}

		obsiLog.Debug("NotesUUIDHandler GET", "note", note)

		data, err := json.Marshal(nm.DomainToHTTP(note))
		if err != nil {

			obsiLog.Debug("NotesUUIDHandler error", "error", "marshal")
			http.Error(w, "Marshal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)

		obsiLog.Debug("NotesUUIDHandler GET OK")
		return

	case http.MethodPut:
		var note httpNote

		// Парсим JSON из request body в структуру
		err := json.NewDecoder(r.Body).Decode(&note)
		if err != nil {

			http.Error(w, "Invalid JSON", http.StatusBadRequest)

			obsiLog.Debug("NotesUUIDHandler PUT Error", "error", "invalid json")
			return
		}

		obsiLog.Debug("NotesUUIDHandler PUT", "httpNote", note)

		domainNote := nm.HTTPToDomain(note)
		if id != domainNote.Id {
			http.Error(w, "Note id != id", http.StatusBadRequest)

			obsiLog.Debug("NotesUUIDHandler PUT Error", "error", "strange id", "id", id, "note id", domainNote.Id)
			return
		}
		domainNote.Id = id

		obsiLog.Debug("NotesUUIDHandler PUT", "domainNote", domainNote)

		err = use_case.Updtr.Update(ctx, domainNote)

		if writeError(w, err) {
			return
		}
		w.WriteHeader(http.StatusOK)
		obsiLog.Debug("NotesUUIDHandler PUT OK")
		return

	case http.MethodDelete:
		err := domain.Repo.DeleteById(ctx, id)

		if writeError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
		obsiLog.Debug("NotesUUIDHandler DELETE OK", "id", id)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		obsiLog.Debug("NotesUUIDHandler default")

		return
	}

}

// /notes&title=key&ancestor=key
func NotesHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	switch r.Method {

	// поиск по имени или предку
	case http.MethodGet:
		params := r.URL.Query()

		if len(params) == 0 {
			notes, err := domain.Repo.GetAll(ctx)
			if writeError(w, err) {
				return
			}

			data, err := json.Marshal(nm.DomainToHTTPSlice(notes))
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)

				obsiLog.Debug("NotesHandler GET all error", "error", "marshal")
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(data)

			obsiLog.Debug("NotesHandler GET all OK")
			return
		} else if len(params) == 1 {

			if _, exist := params["name"]; exist {
				note, err := domain.Repo.FindByName(ctx, params.Get("name"))

				if writeError(w, err) {
					return
				}

				data, err := json.Marshal(nm.DomainToHTTP(note))
				if err != nil {
					http.Error(w, "Internal error", http.StatusInternalServerError)

					obsiLog.Debug("NotesHandler GET name error", "error", "marshal")
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write(data)

				obsiLog.Debug("NotesHandler GET note OK", "note", note)
				return
			}

			if _, exist := params["ancestor"]; exist {
				notes, err := domain.Repo.FindByAncestor(ctx, params.Get("ancestor"))
				if writeError(w, err) {
					return
				}

				data, err := json.Marshal(nm.DomainToHTTPSlice(notes))
				if err != nil {
					http.Error(w, "Internal error", http.StatusInternalServerError)

					obsiLog.Debug("NotesHandler GET ancestor error", "error", "marshal")
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write(data)

				obsiLog.Debug("NotesHandler GET ancestor OK", "notes", notes)
				return

			}
			http.Error(w, "Unknown parameter", http.StatusBadRequest)
			return
		} else if len(params) == 2 {
			http.Error(w, "Two invalid parameters", http.StatusBadRequest)

			obsiLog.Debug("NotesHandler GET error", "error", "two param")
			return
		}

		http.Error(w, "Too many parameters", http.StatusBadRequest)

		obsiLog.Debug("NotesHandler GET error", "error", "to many param")
		return

	// создание новой заметки
	case http.MethodPost:
		n := httpNote{}

		// Парсим JSON из request body в структуру
		err := json.NewDecoder(r.Body).Decode(&n)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)

			obsiLog.Debug("NotesHandler POST error", "error", "marshal")
			return
		}

		note := nm.HTTPToDomain(n)
		_, err = domain.Repo.Insert(ctx, note)
		if writeError(w, err) {
			return
		}

		w.WriteHeader(http.StatusCreated)

		obsiLog.Debug("NotesHandler POST OK", "note", note)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		obsiLog.Debug("NotesHandler default error", "error", "wrong status")
		return
	}
}

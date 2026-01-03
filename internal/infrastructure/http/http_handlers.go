package http

import (
	"encoding/json"
	"net/http"
	"obsidianGoNaive/internal/domain"
	"obsidianGoNaive/internal/use_case"

	"github.com/google/uuid"
)

// HomeHandler handle home request /
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	obsiLog.Debug("Home handler")
	_, err := w.Write([]byte("Home page"))

	if err != nil {

		obsiLog.Error("HomeHandler Write ERROR", "error", err)
		return
	}
}

// NotesUUIDHandler handle notes with id request /notes/{id}
func NotesUUIDHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	tempId := r.PathValue("id")
	id, err := uuid.Parse(tempId)
	if err != nil {

		obsiLog.Error("NotesUUIDHandler parse uuid ERROR", "error", err, "id", id)
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
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

			obsiLog.Error("NotesUUIDHandler ERROR", "error", err)
			http.Error(w, "Marshal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(data)
		if err != nil {

			obsiLog.Error("NotesUUIDHandler GET Write ERROR", "error", err)
			return
		}

		obsiLog.Info("NotesUUIDHandler GET OK")
		return

	case http.MethodPut:

		note := httpNote{}

		err := json.NewDecoder(r.Body).Decode(&note)
		if err != nil {

			obsiLog.Error("NotesUUIDHandler PUT ERROR", "error", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		obsiLog.Debug("NotesUUIDHandler PUT", "httpNote", note)

		domainNote := nm.HTTPToDomain(note)
		if id != domainNote.Id {

			obsiLog.Error("NotesUUIDHandler PUT Error", "error", "strange id", "id", id, "note id", domainNote.Id)
			http.Error(w, "Note id != id", http.StatusBadRequest)
			return
		}
		domainNote.Id = id

		obsiLog.Debug("NotesUUIDHandler PUT", "domainNote", domainNote)

		err = use_case.Updtr.Update(ctx, domainNote)
		if writeError(w, err) {
			return
		}

		w.WriteHeader(http.StatusOK)

		obsiLog.Info("NotesUUIDHandler PUT OK")
		return

	case http.MethodDelete:

		err := domain.Repo.DeleteById(ctx, id)

		if writeError(w, err) {
			return
		}

		w.WriteHeader(http.StatusNoContent)

		obsiLog.Info("NotesUUIDHandler DELETE OK", "id", id)
		return

	default:

		obsiLog.Error("NotesUUIDHandler ERROR", "error", "method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// NotesHandler handle note with name/ancestor request /note?name&?ancestor.
// It is not good to have both of params at one request
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

				obsiLog.Error("NotesHandler GET All ERROR", "error", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(data)
			if err != nil {

				obsiLog.Error("Notes GET All Write ERROR", "error", err)
				return
			}

			obsiLog.Info("Notes GET All OK")
			return

		} else if len(params) == 1 {

			if _, exist := params["name"]; exist {

				note, err := domain.Repo.FindByName(ctx, params.Get("name"))
				if writeError(w, err) {
					return
				}

				obsiLog.Debug("NotesHandler GET Name", "domainNote", note)

				httpNote := nm.DomainToHTTP(note)

				obsiLog.Debug("NotesHandler GET Name", "httpNote", httpNote)

				data, err := json.Marshal(httpNote)
				if err != nil {

					obsiLog.Error("NotesHandler GET Name ERROR", "error", err)
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				_, err = w.Write(data)
				if err != nil {

					obsiLog.Error("Notes GET Name Write ERROR", "error", err)
					return
				}

				obsiLog.Info("NotesHandler GET Name Note OK")
				return
			}

			if _, exist := params["ancestor"]; exist {

				notes, err := domain.Repo.FindByAncestor(ctx, params.Get("ancestor"))
				if writeError(w, err) {
					return
				}

				obsiLog.Debug("NotesHandler GET Ancestor", "domainNotes", notes)

				httpNotes := nm.DomainToHTTPSlice(notes)

				obsiLog.Debug("NotesHandler GET Ancestor", "httpNotes", httpNotes)

				data, err := json.Marshal(httpNotes)
				if err != nil {

					obsiLog.Error("NotesHandler GET Ancestor ERROR", "error", err)
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				_, err = w.Write(data)
				if err != nil {

					obsiLog.Error("NotesHandler GET Ancestor Write ERROR", "error", err)
					return
				}

				obsiLog.Info("NotesHandler GET Ancestor OK")
				return

			}

			obsiLog.Error("NotesHandler GET ERROR", "error", "unknown parameters")
			http.Error(w, "Unknown parameter", http.StatusBadRequest)
			return

		} else if len(params) == 2 {

			obsiLog.Error("NotesHandler GET ERROR", "error", "two param")
			http.Error(w, "Two parameters", http.StatusBadRequest)
			return
		}

		obsiLog.Error("NotesHandler GET ERROR", "error", "to many param")
		http.Error(w, "Too many parameters", http.StatusBadRequest)
		return

	// создание новой заметки
	case http.MethodPost:

		n := httpNote{}

		err := json.NewDecoder(r.Body).Decode(&n)
		if err != nil {

			obsiLog.Error("NotesHandler POST ERROR", "error", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		obsiLog.Debug("NotesHandler POST", "httpNote", n)

		note := nm.HTTPToDomain(n)
		_, err = domain.Repo.Insert(ctx, note)
		if writeError(w, err) {
			return
		}

		obsiLog.Debug("NotesHandler POST", "domainNote", note)

		w.WriteHeader(http.StatusCreated)

		obsiLog.Info("NotesHandler POST OK")
		return

	default:

		obsiLog.Error("NotesHandler Default ERROR", "error", "Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

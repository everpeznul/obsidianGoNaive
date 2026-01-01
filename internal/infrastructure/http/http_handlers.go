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

	w.Write([]byte("Home page"))
}

// /notes/{id}
func NotesUUIDHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	tempId := r.PathValue("id")
	id, err := uuid.Parse(tempId)

	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case "GET":
		note, err := domain.Repo.GetByID(ctx, id)

		if writeError(w, err) {
			return
		}

		data, err := json.Marshal(nm.DomainToHTTP(note))
		if err != nil {

			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return

	case "PUT":
		var note httpNote

		// Парсим JSON из request body в структуру
		err := json.NewDecoder(r.Body).Decode(&note)
		if err != nil {

			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		domainNote := nm.HTTPToDomain(note)
		if id != domainNote.Id {
			http.Error(w, "Note id != id", http.StatusBadRequest)
			return
		}
		domainNote.Id = id

		err = use_case.Updtr.Update(ctx, domainNote)

		if writeError(w, err) {
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	case "DELETE":
		err := domain.Repo.DeleteById(ctx, id)

		if writeError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

}

// /notes&title=key&ancestor=key
func NotesHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	switch r.Method {

	// поиск по имени или предку
	case "GET":
		params := r.URL.Query()

		if len(params) == 0 {
			notes, err := domain.Repo.GetAll(ctx)
			if writeError(w, err) {
				return
			}

			data, err := json.Marshal(nm.DomainToHTTPSlice(notes))
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)

				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(data)
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

					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write(data)
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

					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write(data)
				return

			}
			http.Error(w, "Unknown parameter", http.StatusBadRequest)
			return
		} else if len(params) == 2 {
			http.Error(w, "Two invalid parameters", http.StatusBadRequest)

			return
		}

		http.Error(w, "Too many parameters", http.StatusBadRequest)
		return

	// создание новой заметки
	case "POST":
		n := httpNote{}

		// Парсим JSON из request body в структуру
		err := json.NewDecoder(r.Body).Decode(&n)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		_, err = domain.Repo.Insert(ctx, nm.HTTPToDomain(n))
		if writeError(w, err) {
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}
}

package http

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"obsidianGoNaive/internal/domain"
)

var nt = noteTranslater{}

// /notes/{id}
func NotesUUIDHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	uuid, _ := uuid.Parse(id)

	switch r.Method {

	case "GET":
		note, err := domain.Repo.GetByID(uuid)

		if err != nil {
			http.Error(w, "Note wasn't found", http.StatusNotFound)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(nt.DomainToHTTP(note))
		}

	case "PUT":
		var n note

		// Парсим JSON из request body в структуру
		err := json.NewDecoder(r.Body).Decode(&n)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusMethodNotAllowed)
			return
		}

		err = domain.Repo.UpdateById(nt.HTTPToDomain(n))

		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}

	case "DELETE":
		err := domain.Repo.DeleteByID(uuid)

		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}

	default:
		http.Error(w, "Invalid JSON", http.StatusMethodNotAllowed)
	}

}

// /notes&title=key&ancestor=key
func NotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		params := r.URL.Query()

		if len(params) == 1 {

			if _, exist := params["title"]; exist {
				note, err := domain.Repo.FindByTitle(params.Get("title"))

				if err != nil {
					http.Error(w, "Internal Error", http.StatusInternalServerError)
				} else {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode((nt.DomainToHTTP(note)))
				}
			}

			if _, exist := params["ancestor"]; exist {
				notes, err := domain.Repo.FindByAncestor(params.Get("ancestor"))
				if err != nil {
					http.Error(w, "Internal Error", http.StatusInternalServerError)
				} else {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode((notes))
				}
			}
		} else if len(params) == 2 {
			http.Error(w, "Two invalid parameters", http.StatusBadRequest)
		} else {

			notes, err := domain.Repo.GetAll()
			if err != nil {
				http.Error(w, "Internal Error", http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode((notes))
			}
		}

	case "POST":
		var n note

		// Парсим JSON из request body в структуру
		err := json.NewDecoder(r.Body).Decode(&n)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusMethodNotAllowed)
			return
		}

		domain.Repo.Insert(nt.HTTPToDomain(n))
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	default:
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}
}

package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"obsidianGoNaive/internal/domain"
	"obsidianGoNaive/internal/use_case"
	"strings"

	"github.com/google/uuid"
)

var nt = noteTranslater{}

// /
// домашняя страница
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Home page"))
}

// /notes/{id}
// получение заметки по айди
func NotesUUIDHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	uuid, _ := uuid.Parse(id)

	switch r.Method {

	case "GET":
		note, err := domain.Repo.GetByID(uuid)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
				http.Error(w, "Note wasn't found", http.StatusNotFound)

				return
			} else {
				http.Error(w, "Internal error", http.StatusInternalServerError)

				return
			}

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

		err = use_case.Updtr.Update(nt.HTTPToDomain(n))

		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)

			return
		} else {
			w.WriteHeader(http.StatusOK)
		}

	case "DELETE":
		err := domain.Repo.DeleteByID(uuid)

		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)

			return
		}

	default:
		http.Error(w, "Invalid JSON", http.StatusMethodNotAllowed)

		return
	}

}

// /notes&title=key&ancestor=key
func NotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	// поиск по имени или предку
	case "GET":
		params := r.URL.Query()

		if len(params) == 1 {

			if _, exist := params["title"]; exist {
				note, err := domain.Repo.FindByTitle(params.Get("title"))

				if err != nil {
					http.Error(w, "Internal Error", http.StatusInternalServerError)

					return
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

					return
				} else {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode((notes))
				}
			}
		} else if len(params) == 2 {
			http.Error(w, "Two invalid parameters", http.StatusBadRequest)

			return
		} else {

			notes, err := domain.Repo.GetAll()
			if err != nil {
				http.Error(w, "Internal Error", http.StatusInternalServerError)

				return
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode((notes))
			}
		}

	// создание новой заметки
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

			return
		} else {
			w.WriteHeader(http.StatusOK)
		}
	default:
		http.Error(w, "Invalid JSON", http.StatusBadRequest)

		return
	}
}

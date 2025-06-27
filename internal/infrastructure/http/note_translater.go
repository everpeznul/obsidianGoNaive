package http

import "obsidianGoNaive/internal/domain"

type noteTranslater struct{}

func (nt noteTranslater) HTTPToDomain(n note) domain.Note {

	return domain.Note{
		Id:         n.Id,
		Title:      n.Title,
		Path:       n.Path,
		Class:      n.Class,
		Tags:       n.Tags,
		Links:      n.Links,
		Content:    n.Content, // Используем безопасную версию
		CreateTime: n.CreateTime,
		UpdateTime: n.UpdateTime,
	}
}

func (nt *noteTranslater) DomainToHTTP(n domain.Note) note {

	return note{n.Id,
		n.Title,
		n.Path,
		n.Class,
		n.Tags,
		n.Links,
		n.Content,
		n.CreateTime,
		n.UpdateTime,
	}
}

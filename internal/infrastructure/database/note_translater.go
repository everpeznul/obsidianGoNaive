package database

import "obsidianGoNaive/internal/domain"

type noteTranslater struct{}

func (nt *noteTranslater) DatabaseToDomain(n note) domain.Note {

	return domain.Note{n.Id,
		n.Title,
		n.Path,
		n.Class,
		n.Tags,
		n.Links,
		n.Content[0],
		n.CreateTime,
		n.UpdateTime,
	}
}

func (nt *noteTranslater) DomainToDatabase(n domain.Note) note {

	return note{n.Id,
		n.Title,
		n.Path,
		n.Class,
		n.Tags,
		n.Links,
		[]string{n.Content},
		n.CreateTime,
		n.UpdateTime,
	}
}

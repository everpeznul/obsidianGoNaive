package database

import "obsidianGoNaive/internal/domain"

type noteTranslater struct{}

// internal/infrastructure/database/note_translater.go
func (nt noteTranslater) DatabaseToDomain(n note) domain.Note {
	// Проверяем, что Content не пустой, иначе создаем массив с пустой строкой
	var content string
	if len(n.Content) > 0 {
		content = n.Content[0]
	}
	// content автоматически будет пустой строкой, если массив пустой

	return domain.Note{
		Id:         n.Id,
		Title:      n.Title,
		Path:       n.Path,
		Class:      n.Class,
		Tags:       n.Tags,
		Links:      n.Links,
		Content:    content, // Используем безопасную версию
		CreateTime: n.CreateTime,
		UpdateTime: n.UpdateTime,
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

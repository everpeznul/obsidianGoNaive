package transport

import (
	"context"
	"fmt"
	cmn "obsidianGoNaive/pkg/protos/gen/common"
	pbn "obsidianGoNaive/pkg/protos/gen/notes"

	"github.com/google/uuid"
)

func (u *UpdaterService) Exist(ctx context.Context, title string) (bool, error) {

	_, err := u.Client.Find(ctx, &pbn.FindRequest{Name: title})

	if err != nil {

		return false, fmt.Errorf("note=%s not found: %w", title, err)
	}

	return true, nil
}

func (u *UpdaterService) Create(ctx context.Context, title string) (uuid.UUID, error) {

	resp, err := u.Client.Create(ctx, &pbn.CreateRequest{Note: &cmn.Note{}})
	note := resp.Note
	id := note.Id

	if err != nil {

		return uuid.Nil, fmt.Errorf("create ERROR: %w", err)
	}

	uid, _ := uuid.Parse(id)

	return uid, nil
}

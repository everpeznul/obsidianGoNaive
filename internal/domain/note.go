package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	Id         uuid.UUID
	Title      string
	Path       string
	Class      string
	Tags       []string
	Links      []string
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}

func (n *Note) FindFounder() string {

	founder := strings.Split(n.Title, ".")[0]
	if Exists(founder) {

	}

	return founder
}
func (n *Note) FindAncestor() string {
	a := strings.Split(n.Title, ".")
	var ancestor string
	for i := len(a) - 2; i >= 0; i-- {

		if !strings.Contains(a[i], "%") {
			ancestor = strings.Join(a, ".")
		}
	}

	if Exists(ancestor) {
	}
	return ancestor

}
func (n *Note) FindFather() string {

	f := strings.Split(n.Title, ".")
	father := strings.Join(f[:len(f)-1], ".")
	if Exists(father) {

	}
	return father
}

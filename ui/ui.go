package ui

import (
	"github.com/dontang97/ui/pg"
)

type Source interface {
	Query() ([]pg.User, error)
}

type UI struct {
	src Source
}

func New(src Source) *UI {
	return &UI{
		src: src,
	}
}

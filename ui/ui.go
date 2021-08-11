package ui

import (
	"github.com/dontang97/ui/pg"
)

type UI struct {
	pg.PG
}

func New() *UI {
	return &UI{}
}

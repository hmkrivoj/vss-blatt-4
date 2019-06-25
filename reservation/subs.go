package main

import (
	"context"
	protoShowing "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_showing/proto"
)

type showingDeletedHandler struct {
	db *dataBase
}

func newCinemaShowingDeletedHandler(db *dataBase) *showingDeletedHandler {
	return &showingDeletedHandler{db: db}
}

func (handler *showingDeletedHandler) ShowingDeleted(
	ctx context.Context, // required for signature
	event *protoShowing.DeleteCinemaShowingResponse,
) error {
	_, err := handler.db.removeAllWhereShowingID(event.Showing.Id)
	ctx.Done() // do something with context so the linter will shut up
	return err
}

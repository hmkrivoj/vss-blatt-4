package main

import (
	"context"

	"github.com/micro/go-micro"
	protoCinemaHall "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
	protoCinemaShowing "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_showing/proto"
	protoMovie "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/movie/proto"
)

type cinemaHallDeletedHandler struct {
	db  *dataBase
	pub micro.Publisher
}

func newCinemaHallDeletedHandler(publisher micro.Publisher, db *dataBase) *cinemaHallDeletedHandler {
	handler := &cinemaHallDeletedHandler{db: db}
	handler.pub = publisher
	return handler
}

type movieDeletedHandler struct {
	db  *dataBase
	pub micro.Publisher
}

func newMovieDeletedHandler(publisher micro.Publisher, db *dataBase) *movieDeletedHandler {
	handler := &movieDeletedHandler{db: db}
	handler.pub = publisher
	return handler
}

func (handler *cinemaHallDeletedHandler) CinemaHallDeleted(
	ctx context.Context, // required sub signature
	event *protoCinemaHall.DeleteCinemaHallResponse,
) error {
	deleted, err := handler.db.removeAllWhereCinemaHallID(event.Hall.Id)
	for _, showing := range deleted {
		res := &protoCinemaShowing.DeleteCinemaShowingResponse{}
		res.Showing = &protoCinemaShowing.CinemaShowing{
			Id:         showing.id,
			Movie:      showing.movie,
			CinemaHall: showing.cinemaHall,
		}
		_ = handler.pub.Publish(context.Background(), res)
	}
	ctx.Done() // do something with context so the linter will shut up
	return err
}

func (handler *movieDeletedHandler) MovieDeleted(
	ctx context.Context, // required sub signature
	event *protoMovie.DeleteMovieResponse,
) error {
	deleted, err := handler.db.removeAllWhereMovieID(event.Movie.Id)
	for _, showing := range deleted {
		res := &protoCinemaShowing.DeleteCinemaShowingResponse{}
		res.Showing = &protoCinemaShowing.CinemaShowing{
			Id:         showing.id,
			Movie:      showing.movie,
			CinemaHall: showing.cinemaHall,
		}
		_ = handler.pub.Publish(context.Background(), res)
	}
	ctx.Done() // do something with context so the linter will shut up
	return err
}

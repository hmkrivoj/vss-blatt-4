package main

import (
	"context"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	protoCinemaShowing "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_showing/proto"
)

type serviceHandler struct {
	db  *dataBase
	pub micro.Publisher
}

func newCinemaShowingHandler(publisher micro.Publisher, db *dataBase) *serviceHandler {
	handler := &serviceHandler{}
	handler.pub = publisher
	handler.db = db
	return handler
}

func (handler *serviceHandler) Create(
	ctx context.Context,
	req *protoCinemaShowing.CreateCinemaShowingRequest,
	res *protoCinemaShowing.CreateCinemaShowingResponse,
) error {
	// map proto showing into showing
	showing := cinemaShowing{
		movie:      req.Movie,
		cinemaHall: req.CinemaHall,
	}
	showing = handler.db.create(showing)
	// map showing into proto showing
	res.Showing = &protoCinemaShowing.CinemaShowing{
		Id:         showing.id,
		Movie:      showing.movie,
		CinemaHall: showing.cinemaHall,
	}
	return nil
}

func (handler *serviceHandler) Delete(
	ctx context.Context,
	req *protoCinemaShowing.DeleteCinemaShowingRequest,
	res *protoCinemaShowing.DeleteCinemaShowingResponse,
) error {
	showing, err := handler.db.remove(req.Id)
	if err != nil {
		return err
	}
	// map showing into proto showing
	res.Showing = &protoCinemaShowing.CinemaShowing{
		Id:         showing.id,
		Movie:      showing.movie,
		CinemaHall: showing.cinemaHall,
	}
	// propagate deletion
	err = handler.pub.Publish(context.Background(), res)
	return err
}

func (handler *serviceHandler) FindAll(
	ctx context.Context,
	req *protoCinemaShowing.FindAllCinemaShowingsRequest,
	res *protoCinemaShowing.FindAllCinemaShowingsResponse,
) error {
	showings := handler.db.findAll()
	// map showings into proto showings
	pShowings := make([]*protoCinemaShowing.CinemaShowing, 0)
	for _, showing := range showings {
		pShowings = append(pShowings, &protoCinemaShowing.CinemaShowing{
			Id:         showing.id,
			Movie:      showing.movie,
			CinemaHall: showing.cinemaHall,
		})
	}
	res.Showings = pShowings
	return nil
}

func (handler *serviceHandler) Find(
	cxt context.Context,
	req *protoCinemaShowing.FindCinemaShowingRequest,
	res *protoCinemaShowing.FindCinemaShowingResponse,
) error {
	showing := handler.db.find(req.Id)
	res.Showing = &protoCinemaShowing.CinemaShowing{
		Id:         showing.id,
		CinemaHall: showing.cinemaHall,
		Movie:      showing.movie,
	}
	return nil
}

func main() {
	service := micro.NewService(micro.Name("cinema.cinema_showing.service"))
	service.Init()

	// init dependencies
	publisher := micro.NewPublisher("cinema.cinema_showing.deleted", service.Client())
	db := newDataBase()

	// inject dependencies
	handler := newCinemaShowingHandler(publisher, db)
	deletedHallHandler := newCinemaHallDeletedHandler(publisher, db)
	deletedMovieHandler := newMovieDeletedHandler(publisher, db)

	// register handlers
	err := micro.RegisterSubscriber("cinema.cinema_hall.deleted", service.Server(), deletedHallHandler)
	if err != nil {
		log.Fatal(err)
	}
	err = micro.RegisterSubscriber("cinema.movie.deleted", service.Server(), deletedMovieHandler)
	if err != nil {
		log.Fatal(err)
	}
	err = protoCinemaShowing.RegisterCinemaShowingServiceHandler(service.Server(), handler)
	if err != nil {
		log.Fatal(err)
	}

	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}

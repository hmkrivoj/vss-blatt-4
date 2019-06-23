package main

import (
	"context"
	"errors"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	protoCinemaHall "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
	protoCinemaShowing "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_showing/proto"
	protoMovie "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/movie/proto"
	"sync"
)

type CinemaHallDeletedHandler struct{}

func NewCinemaHallDeletedHandler() *CinemaHallDeletedHandler {
	return &CinemaHallDeletedHandler{}
}

type MovieDeletedHandler struct{}

func NewMovieDeletedHandler() *MovieDeletedHandler {
	return &MovieDeletedHandler{}
}

type CinemaShowingHandler struct {
	mutex          sync.Mutex
	idCounter      int64
	cinemaShowings map[int64]cinemaShowing
	pub            micro.Publisher
}

type cinemaShowing struct {
	id         int64
	movie      int64
	cinemaHall int64
}

func NewCinemaShowingHandler(publisher micro.Publisher) *CinemaShowingHandler {
	handler := &CinemaShowingHandler{}
	handler.idCounter = 1
	handler.cinemaShowings = make(map[int64]cinemaShowing)
	handler.pub = publisher
	return handler
}

func (handler *CinemaShowingHandler) Create(ctx context.Context, req *protoCinemaShowing.CreateCinemaShowingRequest, res *protoCinemaShowing.CreateCinemaShowingResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	id := handler.idCounter
	handler.idCounter++

	showing := cinemaShowing{
		id:         id,
		movie:      req.Movie,
		cinemaHall: req.CinemaHall,
	}
	handler.cinemaShowings[id] = showing

	res.Showing = &protoCinemaShowing.CinemaShowing{
		Id:         showing.id,
		Movie:      showing.movie,
		CinemaHall: showing.cinemaHall,
	}

	return nil
}

func (handler *CinemaShowingHandler) Delete(ctx context.Context, req *protoCinemaShowing.DeleteCinemaShowingRequest, res *protoCinemaShowing.DeleteCinemaShowingResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	if _, ok := handler.cinemaShowings[req.Id]; !ok {
		return errors.New("no such id")
	}

	showing := handler.cinemaShowings[req.Id]
	delete(handler.cinemaShowings, req.Id)

	res.Showing = &protoCinemaShowing.CinemaShowing{
		Id:         showing.id,
		Movie:      showing.movie,
		CinemaHall: showing.cinemaHall,
	}
	err := handler.pub.Publish(context.Background(), res)

	return err
}

func (handler *CinemaShowingHandler) FindAll(ctx context.Context, req *protoCinemaShowing.FindAllCinemaShowingsRequest, res *protoCinemaShowing.FindAllCinemaShowingsResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	showings := make([]*protoCinemaShowing.CinemaShowing, 0)
	for _, showing := range handler.cinemaShowings {
		showings = append(showings, &protoCinemaShowing.CinemaShowing{
			Id:         showing.id,
			Movie:      showing.movie,
			CinemaHall: showing.cinemaHall,
		})
	}
	res.Showings = showings

	return nil
}

func (handler *MovieDeletedHandler) CinemaHallDeleted(context.Context, *protoCinemaHall.DeleteCinemaHallResponse) error {
	log.Logf("Received hall deleted")
	return nil
}

func (handler *CinemaHallDeletedHandler) MovieDeleted(context.Context, *protoMovie.DeleteMovieResponse) error {
	log.Logf("Received movie deleted")
	return nil
}

func main() {
	service := micro.NewService(micro.Name("cinema.cinema_showing.service"))
	service.Init()

	publisher := micro.NewPublisher("cinema.cinema_showing.deleted", service.Client())
	handler := NewCinemaShowingHandler(publisher)
	deletedHallHandler := NewCinemaHallDeletedHandler()
	deletedMovieHandler := NewMovieDeletedHandler()

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

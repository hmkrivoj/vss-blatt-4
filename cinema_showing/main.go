package main

import (
	"context"
	"errors"
	"sync"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	protoCinemaHall "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
	protoCinemaShowing "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_showing/proto"
	protoMovie "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/movie/proto"
)

type cinemaShowing struct {
	id         int64
	movie      int64
	cinemaHall int64
}

type dataBase struct {
	mutex          sync.Mutex
	idCounter      int64
	cinemaShowings map[int64]cinemaShowing
}

func newDataBase() *dataBase {
	db := &dataBase{}
	db.idCounter = 1
	db.cinemaShowings = make(map[int64]cinemaShowing)

	return db
}

func (db *dataBase) create(showing cinemaShowing) cinemaShowing {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	id := db.idCounter
	db.idCounter++

	created := cinemaShowing{
		id:         id,
		movie:      showing.movie,
		cinemaHall: showing.cinemaHall,
	}
	db.cinemaShowings[id] = created

	return created
}

func (db *dataBase) remove(id int64) (cinemaShowing, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if _, ok := db.cinemaShowings[id]; !ok {
		return cinemaShowing{}, errors.New("no such id")
	}
	showing := db.cinemaShowings[id]
	delete(db.cinemaShowings, id)
	return showing, nil
}

func (db *dataBase) removeAllWhereMovieID(movie int64) ([]cinemaShowing, error) {
	showings := db.findAll()
	toBeRemoved := make([]int64, 0)
	for _, showing := range showings {
		if showing.movie == movie {
			toBeRemoved = append(toBeRemoved, showing.id)
		}
	}
	removed := make([]cinemaShowing, 0)
	for _, id := range toBeRemoved {
		showing, err := db.remove(id)
		if err != nil {
			return removed, err
		}
		removed = append(removed, showing)
	}
	return removed, nil
}

func (db *dataBase) removeAllWhereCinemaHallID(cinemaHall int64) ([]cinemaShowing, error) {
	showings := db.findAll()
	toBeRemoved := make([]int64, 0)
	for _, showing := range showings {
		if showing.cinemaHall == cinemaHall {
			toBeRemoved = append(toBeRemoved, showing.id)
		}
	}
	removed := make([]cinemaShowing, 0)
	for _, id := range toBeRemoved {
		showing, err := db.remove(id)
		if err != nil {
			return removed, err
		}
		removed = append(removed, showing)
	}
	return removed, nil
}

func (db *dataBase) findAll() []cinemaShowing {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	showings := make([]cinemaShowing, 0)
	for _, showing := range db.cinemaShowings {
		showings = append(showings, showing)
	}
	return showings
}

func (db *dataBase) find(id int64) cinemaShowing {
	return db.cinemaShowings[id]
}

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

type serviceHandler struct {
	db  *dataBase
	pub micro.Publisher
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
	showing := cinemaShowing{
		movie:      req.Movie,
		cinemaHall: req.CinemaHall,
	}
	showing = handler.db.create(showing)

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
	res.Showing = &protoCinemaShowing.CinemaShowing{
		Id:         showing.id,
		Movie:      showing.movie,
		CinemaHall: showing.cinemaHall,
	}
	err = handler.pub.Publish(context.Background(), res)
	return err
}

func (handler *serviceHandler) FindAll(
	ctx context.Context,
	req *protoCinemaShowing.FindAllCinemaShowingsRequest,
	res *protoCinemaShowing.FindAllCinemaShowingsResponse,
) error {
	showings := handler.db.findAll()
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

func main() {
	service := micro.NewService(micro.Name("cinema.cinema_showing.service"))
	service.Init()

	publisher := micro.NewPublisher("cinema.cinema_showing.deleted", service.Client())
	db := newDataBase()
	handler := newCinemaShowingHandler(publisher, db)
	deletedHallHandler := newCinemaHallDeletedHandler(publisher, db)
	deletedMovieHandler := newMovieDeletedHandler(publisher, db)

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

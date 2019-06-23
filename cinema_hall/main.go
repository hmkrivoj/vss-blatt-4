package main

import (
	"context"
	"errors"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
	"sync"
)

type CinemaHallHandler struct {
	mutex       sync.Mutex
	idCounter   int64
	cinemaHalls map[int64]cinemaHall
	pub         micro.Publisher
}

func NewCinemaHallHandler(publisher micro.Publisher) *CinemaHallHandler {
	handler := &CinemaHallHandler{}
	handler.idCounter = 1
	handler.cinemaHalls = make(map[int64]cinemaHall)
	handler.pub = publisher
	return handler
}

type cinemaHall struct {
	id   int64
	name string
	rows int32
	cols int32
}

func (handler *CinemaHallHandler) Create(ctx context.Context, req *proto.CreateCinemaHallRequest, res *proto.CreateCinemaHallResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	id := handler.idCounter
	handler.idCounter++

	hall := cinemaHall{
		id:   id,
		name: req.Name,
		rows: req.Rows,
		cols: req.Cols,
	}
	handler.cinemaHalls[id] = hall

	res.Hall = &proto.CinemaHall{
		Id:   hall.id,
		Name: hall.name,
		Rows: hall.rows,
		Cols: hall.cols,
	}

	return nil
}

func (handler *CinemaHallHandler) Delete(ctx context.Context, req *proto.DeleteCinemaHallRequest, res *proto.DeleteCinemaHallResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	if _, ok := handler.cinemaHalls[req.Id]; !ok {
		return errors.New("no such id")
	}

	hall := handler.cinemaHalls[req.Id]
	delete(handler.cinemaHalls, req.Id)

	res.Hall = &proto.CinemaHall{
		Id:   hall.id,
		Rows: hall.rows,
		Cols: hall.cols,
		Name: hall.name,
	}
	err := handler.pub.Publish(context.Background(), res)

	return err
}

func (handler *CinemaHallHandler) FindAll(ctx context.Context, req *proto.FindAllCinemaHallsRequest, res *proto.FindAllCinemaHallsResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	halls := make([]*proto.CinemaHall, 0)
	for _, hall := range handler.cinemaHalls {
		halls = append(halls, &proto.CinemaHall{
			Id:   hall.id,
			Name: hall.name,
			Rows: hall.rows,
			Cols: hall.cols,
		})
	}
	res.Halls = halls

	return nil
}

func main() {
	service := micro.NewService(micro.Name("cinema.cinema_hall.service"))
	service.Init()

	publisher := micro.NewPublisher("cinema.cinema_hall.deleted", service.Client())
	handler := NewCinemaHallHandler(publisher)

	err := proto.RegisterCinemaHallServiceHandler(service.Server(), handler)
	if err != nil {
		log.Fatal(err)
	}

	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}

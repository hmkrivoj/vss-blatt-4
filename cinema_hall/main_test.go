package main

import (
	"context"
	"errors"
	"testing"

	"github.com/micro/go-micro/client"

	proto "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
	"github.com/stretchr/testify/assert"
)

type mockPub struct {
	noCalls int32
}

func newMockPub() *mockPub {
	return &mockPub{noCalls: 0}
}

func (m *mockPub) Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error {
	m.noCalls++
	return nil
}

type mockPubError struct {
	noCalls int32
}

func newMockPubError() *mockPubError {
	return &mockPubError{noCalls: 0}
}

func (m *mockPubError) Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error {
	m.noCalls++
	return errors.New("")
}

func TestCinemaHallHandler_Create(t *testing.T) {
	handler := NewCinemaHallHandler(nil)
	assert.Nil(t, handler.pub)
	assert.Equal(t, 1, int(handler.idCounter))
	assert.Equal(t, 0, len(handler.cinemaHalls))
	res := &proto.CreateCinemaHallResponse{}
	err := handler.Create(context.TODO(), &proto.CreateCinemaHallRequest{Name: "Alpha"}, res)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(handler.cinemaHalls))
	assert.Equal(t, 1, int(res.Hall.Id))
	assert.Equal(t, "Alpha", res.Hall.Name)

	res = &proto.CreateCinemaHallResponse{}
	err = handler.Create(context.TODO(), &proto.CreateCinemaHallRequest{Name: "Gamma"}, res)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(handler.cinemaHalls))
	assert.Equal(t, 2, int(res.Hall.Id))
	assert.Equal(t, "Gamma", res.Hall.Name)
}

func TestCinemaHallHandler_Delete(t *testing.T) {
	mock := newMockPub()
	handler := NewCinemaHallHandler(mock)
	err := handler.Create(
		context.TODO(),
		&proto.CreateCinemaHallRequest{Name: "Alpha"},
		&proto.CreateCinemaHallResponse{},
	)
	assert.Nil(t, err)
	err = handler.Create(
		context.TODO(),
		&proto.CreateCinemaHallRequest{Name: "Beta"},
		&proto.CreateCinemaHallResponse{},
	)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(handler.cinemaHalls))
	// delete non existing cinema hall
	assert.Equal(t, 0, int(mock.noCalls))
	err = handler.Delete(context.TODO(), &proto.DeleteCinemaHallRequest{Id: 42}, &proto.DeleteCinemaHallResponse{})
	assert.NotNil(t, err)
	assert.Equal(t, 0, int(mock.noCalls))
	assert.Equal(t, 0, int(mock.noCalls))
	err = handler.Delete(context.TODO(), &proto.DeleteCinemaHallRequest{Id: 1337}, &proto.DeleteCinemaHallResponse{})
	assert.NotNil(t, err)
	assert.Equal(t, 0, int(mock.noCalls))

	// delete existing cinema hall
	res := &proto.DeleteCinemaHallResponse{}
	err = handler.Delete(context.TODO(), &proto.DeleteCinemaHallRequest{Id: 1}, res)
	assert.Nil(t, err)
	assert.Equal(t, 3, int(handler.idCounter))
	assert.Equal(t, 1, len(handler.cinemaHalls))
	assert.Equal(t, 1, int(mock.noCalls))
	assert.Equal(t, 1, int(res.Hall.Id))
	assert.Equal(t, "Alpha", res.Hall.Name)

	// simulate error from publisher
	mockError := newMockPubError()
	handlerError := NewCinemaHallHandler(mockError)
	_ = handlerError.Create(
		context.TODO(),
		&proto.CreateCinemaHallRequest{Name: "Gamma"},
		&proto.CreateCinemaHallResponse{},
	)
	// Delete existing cinema hall, expect error from publisher
	res = &proto.DeleteCinemaHallResponse{}
	assert.Equal(t, 0, int(mockError.noCalls))
	err = handlerError.Delete(context.TODO(), &proto.DeleteCinemaHallRequest{Id: 1}, res)
	assert.NotNil(t, err)
	assert.Equal(t, 1, int(mockError.noCalls))
}

func TestCinemaHallHandler_FindAll(t *testing.T) {
	handler := NewCinemaHallHandler(nil)
	assert.Equal(t, 0, len(handler.cinemaHalls))
	res := &proto.FindAllCinemaHallsResponse{}
	err := handler.FindAll(context.TODO(), &proto.FindAllCinemaHallsRequest{}, res)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(res.Halls))
	err = handler.Create(
		context.TODO(),
		&proto.CreateCinemaHallRequest{Name: "Alpha"},
		&proto.CreateCinemaHallResponse{},
	)
	assert.Nil(t, err)
	err = handler.Create(
		context.TODO(),
		&proto.CreateCinemaHallRequest{Name: "Beta"},
		&proto.CreateCinemaHallResponse{},
	)
	assert.Nil(t, err)
	err = handler.Create(
		context.TODO(),
		&proto.CreateCinemaHallRequest{Name: "Gamma"},
		&proto.CreateCinemaHallResponse{},
	)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(handler.cinemaHalls))
	res = &proto.FindAllCinemaHallsResponse{}
	err = handler.FindAll(context.TODO(), &proto.FindAllCinemaHallsRequest{}, res)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(res.Halls))
}

func TestCinemaHallHandler_Find(t *testing.T) {
	handler := NewCinemaHallHandler(nil)
	assert.Nil(t, handler.pub)
	assert.Equal(t, 1, int(handler.idCounter))
	assert.Equal(t, 0, len(handler.cinemaHalls))
	_ = handler.Create(context.TODO(),
		&proto.CreateCinemaHallRequest{Name: "Alpha"},
		&proto.CreateCinemaHallResponse{},
	)

	res := &proto.FindCinemaHallResponse{}
	err := handler.Find(context.TODO(), &proto.FindCinemaHallRequest{Id: 1}, res)
	assert.Nil(t, err)
	assert.Equal(t, 1, int(res.Hall.Id))
	assert.Equal(t, "Alpha", res.Hall.Name)
}

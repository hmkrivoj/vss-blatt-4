package main

import (
	"context"
	"errors"
	"testing"

	"github.com/micro/go-micro/client"

	proto "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/movie/proto"
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

func TestMovieHandler_Create(t *testing.T) {
	handler := NewMovieHandler(nil)
	assert.Nil(t, handler.pub)
	assert.Equal(t, 1, int(handler.idCounter))
	assert.Equal(t, 0, int(len(handler.movies)))
	res := &proto.CreateMovieResponse{}
	err := handler.Create(context.TODO(), &proto.CreateMovieRequest{Title: "Life of Brian"}, res)
	assert.Nil(t, err)
	assert.Equal(t, 1, int(len(handler.movies)))
	assert.Equal(t, 1, int(res.Movie.Id))
	assert.Equal(t, "Life of Brian", res.Movie.Title)

	res = &proto.CreateMovieResponse{}
	err = handler.Create(context.TODO(), &proto.CreateMovieRequest{Title: "Monty Python and the holy grail"}, res)
	assert.Nil(t, err)
	assert.Equal(t, 2, int(len(handler.movies)))
	assert.Equal(t, 2, int(res.Movie.Id))
	assert.Equal(t, "Monty Python and the holy grail", res.Movie.Title)
}

func TestMovieHandler_Delete(t *testing.T) {
	mock := newMockPub()
	handler := NewMovieHandler(mock)
	err := handler.Create(
		context.TODO(),
		&proto.CreateMovieRequest{Title: "Life of Brian"},
		&proto.CreateMovieResponse{},
	)
	assert.Nil(t, err)
	err = handler.Create(
		context.TODO(),
		&proto.CreateMovieRequest{Title: "Flying Circus"},
		&proto.CreateMovieResponse{},
	)
	assert.Nil(t, err)
	assert.Equal(t, 2, int(len(handler.movies)))
	// delete non existing movie
	assert.Equal(t, 0, int(mock.noCalls))
	err = handler.Delete(context.TODO(), &proto.DeleteMovieRequest{Id: 42}, &proto.DeleteMovieResponse{})
	assert.NotNil(t, err)
	assert.Equal(t, 0, int(mock.noCalls))

	// delete existing movie
	res := &proto.DeleteMovieResponse{}
	err = handler.Delete(context.TODO(), &proto.DeleteMovieRequest{Id: 1}, res)
	assert.Nil(t, err)
	assert.Equal(t, 3, int(handler.idCounter))
	assert.Equal(t, 1, int(len(handler.movies)))
	assert.Equal(t, 1, int(mock.noCalls))
	assert.Equal(t, 1, int(res.Movie.Id))
	assert.Equal(t, "Life of Brian", res.Movie.Title)

	// simulate error from publisher
	mockError := newMockPubError()
	handlerError := NewMovieHandler(mockError)
	err = handlerError.Create(
		context.TODO(),
		&proto.CreateMovieRequest{Title: "Life of Brian"},
		&proto.CreateMovieResponse{},
	)
	// Delete existing movie, expect error from publisher
	res = &proto.DeleteMovieResponse{}
	assert.Equal(t, 0, int(mockError.noCalls))
	err = handlerError.Delete(context.TODO(), &proto.DeleteMovieRequest{Id: 1}, res)
	assert.NotNil(t, err)
	assert.Equal(t, 1, int(mockError.noCalls))
}

func TestMovieHandler_FindAll(t *testing.T) {
	handler := NewMovieHandler(nil)
	assert.Equal(t, 0, int(len(handler.movies)))
	res := &proto.FindAllMoviesResponse{}
	err := handler.FindAll(context.TODO(), &proto.FindAllMoviesRequest{}, res)
	assert.Nil(t, err)
	assert.Equal(t, 0, int(len(res.Movies)))
	err = handler.Create(
		context.TODO(),
		&proto.CreateMovieRequest{Title: "Life of Brian"},
		&proto.CreateMovieResponse{},
	)
	assert.Nil(t, err)
	err = handler.Create(
		context.TODO(),
		&proto.CreateMovieRequest{Title: "Flying Circus"},
		&proto.CreateMovieResponse{},
	)
	assert.Nil(t, err)
	assert.Equal(t, 2, int(len(handler.movies)))
	res = &proto.FindAllMoviesResponse{}
	err = handler.FindAll(context.TODO(), &proto.FindAllMoviesRequest{}, res)
	assert.Nil(t, err)
	assert.Equal(t, 2, int(len(res.Movies)))
}

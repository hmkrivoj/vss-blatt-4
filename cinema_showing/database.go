package main

import (
	"errors"
	"sync"
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
	// lock to avoid trouble with findall
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
	// lock to avoid trouble with findall
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
	// determine which showings have to be deleted
	for _, showing := range showings {
		if showing.movie == movie {
			toBeRemoved = append(toBeRemoved, showing.id)
		}
	}
	// delete showings and track them
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
	// determine which showings have to be deleted
	toBeRemoved := make([]int64, 0)
	for _, showing := range showings {
		if showing.cinemaHall == cinemaHall {
			toBeRemoved = append(toBeRemoved, showing.id)
		}
	}
	// delete showings and track them
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
	// lock to avoid trouble while iterating
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

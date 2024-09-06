package films

import (
	"minimal-http-server/database"
	"strconv"
	"sync"
)

func GetFilms() []database.Film {
	db := database.Database()

	return db.Films
}

func GetFilm(id string) (database.Film, error) {
	db := database.Database()

	for _, film := range db.Films {
		if film.Id == id {
			return film, nil
		}
	}

	return database.Film{}, nil
}

var mu sync.Mutex

func AddFilm(film database.Film) {
	mu.Lock()
	defer mu.Unlock()

	db := database.Database()

	film.Id = strconv.Itoa(len(db.Films) + 1)

	db.Films = append(db.Films, film)

	db.Save()
}

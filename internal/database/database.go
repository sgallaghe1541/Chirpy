package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	Path string
	mtx  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		Path: path,
		mtx:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) loadDB() (DBStructure, error) {
	file, err := os.ReadFile(db.Path)
	if err != nil {
		return DBStructure{}, err
	}

	data := DBStructure{}

	err = json.Unmarshal(file, &data.Chirps)
	if err != nil {
		return DBStructure{}, err
	}

	return data, nil
}

func (db *DB) writeDB(dbdata DBStructure) error {
	filecontents, err := json.Marshal(dbdata.Chirps)
	if err != nil {
		return err
	}

	return os.WriteFile(db.Path, filecontents, 0644)
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mtx.RLock()

	postedChirps, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := getNextID(postedChirps.Chirps)

	newChirp := Chirp{
		Body: body,
		Id:   id,
	}

	postedChirps.Chirps[id] = newChirp

	err = db.writeDB(postedChirps)
	if err != nil {
		return Chirp{}, err
	}

	db.mtx.RUnlock()

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mtx.RLock()

	postedChirps, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := make([]Chirp, 0, len(postedChirps.Chirps))
	for _, chirp := range postedChirps.Chirps {
		chirps = append(chirps, chirp)
	}

	err = db.writeDB(postedChirps)
	if err != nil {
		return []Chirp{}, err
	}
	db.mtx.RUnlock()

	return chirps, nil
}

func (db *DB) ensureDB() error {
	_, err := os.Stat(db.Path)
	if !errors.Is(err, os.ErrNotExist) {
		return nil
	}

	file, err := os.Create(db.Path)
	if err != nil {
		return err
	}

	defer file.Close()
	return nil
}

func getNextID(data map[int]Chirp) int {
	var maxID int
	for maxID = range data {
		break
	}

	for n := range data {
		if n > maxID {
			maxID = n
		}
	}

	maxID++

	return maxID
}

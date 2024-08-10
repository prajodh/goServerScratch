package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct{
	path string
	max *sync.RWMutex
}

type Chirp string

var index int;

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	// Try to open the file with O_CREATE and O_WRONLY flags
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Initialize the DB struct
	db := &DB{
		path: path,
		max:  &sync.RWMutex{},
	}

	return db, nil
}


// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error){
	data, err := db.loadDB()
	insertData := Chirp(body)
	if err != nil{
		return insertData, err
	}
	index += 1
	mu := db.max
	mu.Lock()
	defer mu.Unlock()
	data.Chirps[index] = insertData
	db.writeDB(data)
	return insertData, nil

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() (map[int]Chirp, error){
	dbstruct, err := db.loadDB()
	if err != nil{
		return nil, err
	}
	// for _,v := range(dbstruct.Chirps){
	// 	result = append(result, string(v))
	// }
	return dbstruct.Chirps, nil	
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error){
	data, err := os.ReadFile(db.path)
	dbb := DBStructure{Chirps: map[int]Chirp{}}
	if err != nil{
		return dbb, err
	}
	err_ := json.Unmarshal(data,&dbb)
	if err_ != nil{
		return dbb, err
	}
	return dbb, err
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error{
	data, err := json.Marshal(dbStructure)
	if err != nil{
		return err
	}
	os.WriteFile(db.path, data, 0777)
	return nil
}


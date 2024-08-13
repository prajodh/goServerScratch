package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type DB struct{
	path string
	max *sync.RWMutex
}

type Chirp string
type Email string
type Password string


var index int;
var userIndex int;

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[Email]Password
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

func (db *DB)CreateUsers(email string, password string) (string, error){
	data, err := db.loadDB()
	if err != nil{
		return "", err
	}
	mu := db.max
	mu.Lock()
	defer mu.Unlock()
	data.Users[Email(email)] = Password(password)
	db.writeDB(data)
	return strconv.Itoa(userIndex)+" : "+ email, nil

}

func (db *DB)GetUsers()(map[Email]Password, error){
	data, err := db.loadDB()
	if err!= nil{
		return nil, err
	}
	return data.Users, nil

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
	dbb := DBStructure{
    Chirps: map[int]Chirp{},
    Users: map[Email]Password{},
	}
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
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(string(data))
	if err != nil{
		return err
	}
	os.WriteFile(db.path, data, 0777)
	return nil
}


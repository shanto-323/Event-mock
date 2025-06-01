package main

import (
	"fmt"
	"log"
)

type Repository interface {
	MockGetData(m *GetModel)
	MockCreateData(m *CreateModel)
}

type dbRepository struct{}

func NewDbRepository() Repository {
	return &dbRepository{}
}

func (d *dbRepository) MockGetData(m *GetModel) {
	log.Println(fmt.Sprint("getting data from database: ", m))
}

func (d *dbRepository) MockCreateData(m *CreateModel) {
	log.Println(fmt.Sprint("data created in database: ", m))
}

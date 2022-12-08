package models

import "time"

type Temperaturi struct {
	Id        int       `bson:"id, omitempty" json:"id"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Valoare   float32   `bson:"valoare" json:"valoare" validate:"required"`
	Id_oras   *int      `bson:"id_oras" json:"idOras" validate:"required"`
}

type Temperaturi_Response struct {
	Id        int       `bson:"id, omitempty" json:"id"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Valoare   float32   `bson:"valoare" json:"valoare" validate:"required"`
}

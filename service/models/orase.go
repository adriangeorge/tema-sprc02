package models

type Orase struct {
	Id          int     `bson:"id, omitempty" json:"id"`
	Id_tara     *int    `bson:"id_tara" json:"idTara" validate:"required"`
	Nume_oras   string  `bson:"nume_oras" json:"nume" validate:"required"`
	Latitudine  float64 `bson:"latitudine" json:"lat" validate:"required"`
	Longitudine float64 `bson:"longitudine" json:"lon" validate:"required"`
}

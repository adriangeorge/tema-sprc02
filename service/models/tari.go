package models

type Tara struct {
	Id          int     `bson:"id, omitempty" json:"id"`
	Nume_tara   string  `bson:"nume_tara" json:"nume" validate:"required"`
	Latitudine  float64 `bson:"latitudine" json:"lat" validate:"required"`
	Longitudine float64 `bson:"longitudine" json:"lon" validate:"required"`
}

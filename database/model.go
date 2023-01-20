package database

type Model interface {
	GetID() int
	Create() error
	FindByID(id int) (Model, error)
	Update() error
	Delete() error
}

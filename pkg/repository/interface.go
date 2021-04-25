package repository

type Repository interface {
	Save(string) error
}

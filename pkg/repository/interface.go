package repository

type Repository interface {
	Save(string) (string, error)
}

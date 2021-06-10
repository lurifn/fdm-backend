package repository

type Repository interface {
	Save(string) (string, error)
	Update(string2 string) (string, error)
}

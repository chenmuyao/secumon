package dao

type LogDAO interface{}

type GORMLogDAO struct{}

func NewLogDAO() LogDAO {
	return &GORMLogDAO{}
}

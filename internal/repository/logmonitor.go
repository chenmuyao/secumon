package repository

type LogRepo interface{}

type logRepo struct{}

func NewLogRepo() LogRepo {
	return &logRepo{}
}

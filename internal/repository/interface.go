package repository

import (
	"blog_service2/internal/model"
)

type Storage interface {
	Insert(post model.Record) error
	Remove(id int) error
	Update(post model.Record) (int64, error)
	ReadOne(id int) (model.Record, error)
	Read(str string) ([]model.Record, error)
}

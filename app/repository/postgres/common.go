package postgres_repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func whereBy[T uuid.UUID | string | int | bool](str T, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", field), str)
	}
}

func paginate(limit, offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit).Offset(offset)
	}
}

func whereByText(text, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		searchText := "%" + text + "%"
		return db.Where(fmt.Sprintf("%s LIKE ?", field), searchText)
	}
}

func whereBySlice[T []uuid.UUID | []string | []int](data T, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IN ?", field), data)
	}
}

func orBy[T uuid.UUID | string | int | bool](str T, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Or(fmt.Sprintf("%s = ?", field), str)
	}
}

func orBySlice[T []uuid.UUID | []string | []int](data T, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Or(fmt.Sprintf("%s IN ?", field), data)
	}
}

func orByText(text, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		searchText := "%" + text + "%"
		return db.Or(fmt.Sprintf("%s LIKE ?", field), searchText)
	}
}

package model

import (
	"github.com/jinzhu/gorm"
)

type Chemical struct {
	gorm.Model
	Species       string
	Name          string
	CriticalValue string
}

func (c Chemical) String() string {
	return c.Name
}

type Species struct {
	label string
	Name  []string
}

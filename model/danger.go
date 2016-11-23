package model

import (
	"github.com/jinzhu/gorm"
)

type Danger struct {
	gorm.Model
	Species       string
	Name          string
	CriticalValue string
}

func (c Danger) String() string {
	return c.Name
}

type Species struct {
	label string
	Name  []string
}

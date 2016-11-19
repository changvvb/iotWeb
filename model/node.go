package model

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type Node struct {
	gorm.Model       `json:"-"`
	Species          string `gorm:"type:varchar(100)"`
	MaxValue         float64
	MinValue         float64
	Describe         string
	X                int
	Y                int
	PositionDescribe string `json:"-"`
	Data             []Data `gorm:"ForeignKey:NodeRefer" json:"-"`
	ParkRefer        uint
	Number           int
}

type Data struct {
	Time      time.Time
	Data      float64
	NodeRefer uint
}

type Point struct {
	X int
	Y int
}

type Points struct {
	Points []Point
}

func (n *Node) GetData() {
	db, err := opendb()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	db.Where("node_refer=?", n.ID).Find(&n.Data)
}

func (n *Node) InsertData(v float64) {
	db, err := opendb()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	db.Create(&Data{Time: time.Now(), Data: v, NodeRefer: n.ID})
}

package model

import (
	"github.com/jinzhu/gorm"
)

type Park struct {
	gorm.Model
	Name    string
	Node    []Node `gorm:"ForeignKey:ParkRefer" json:"-"`
	Tel     string
	Address string
}

func (p *Park) GetNodes() []Node {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
		return nil
	}
	defer db.Close()
	db.Where("park_refer=?", p.ID).Find(&p.Node)
	// db.Model(p).Related(&p.Node, "nodes")
	printlnLog(p.Node)
	return p.Node
}

func (p *Park) AddNode(n *Node) {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
		return
	}
	defer db.Close()
	db.Create(n)
}

package model

import (
	"github.com/jinzhu/gorm"
)

type Park struct {
	gorm.Model
	Name    string
	Tel     string
	Address string
	// Species string
	Nodes []Node `gorm:"ForeignKey:ParkRefer"`
}

func (p *Park) GetNodes() []Node {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
		return nil
	}
	defer db.Close()

	db.Where("park_refer=?", p.ID).Find(&p.Nodes)
	s := make(map[string]bool)
	for index, _ := range p.Nodes {
		p.Nodes[index].DangerID = p.Nodes[index].DangerID
		db.Where("id=?", p.Nodes[index].DangerID).Find(&p.Nodes[index].Danger)
		s[p.Nodes[index].Danger.Species] = true
	}

	// for index, _ := range s {
	//     p.Species += (" " + index)
	//     // if s[index] == true {
	//     //     p.Species = append(p.Species, index)
	//     // }
	// }

	return p.Nodes
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

func (p *Park) GetSpecies() {

}

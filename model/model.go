package model

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	// "time"
)

type User struct {
	gorm.Model
}

func init() {
	db, err := opendb()
	defer db.Close()
	if err != nil {
		printlnLog(err)
		return
	}
	// datas := []Data{{Time: time.Now(), Data: 34}, {Time: time.Now(), Data: 56}}

	// node := Node{Data: datas, Species: "my node", MaxValue: 100, MinValue: 0, Describe: "haha"}
	db.AutoMigrate(&Node{}, &Data{})
	//db.Create(&node)
}

func GetNodes() []Node {
	db, err := opendb()
	if err != nil {
		printlnLog("model:", err)
	}

	var nodes []Node
	db.Find(&nodes)

	defer db.Close()
	return nodes
}

func GetIdByPosition(x int, y int) uint {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
	}
	defer db.Close()
	node := Node{}
	db.Where("x=? AND y=?", x, y).First(&node)
	return node.ID
}

func GetNodeByID(ID uint) *Node {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
	}
	defer db.Close()
	node := Node{}
	db = db.Model(&node).Where("id=?", ID)
	if db.Find(&node).RecordNotFound() {
		log.Println("NotFound")
		return nil
	}
	return &node
}

//打开数据库
func opendb() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:root@/iot?charset=utf8&parseTime=True&loc=Local")
	return db, err
}

func printlnLog(v ...interface{}) {
	log.Println("model:", v)
}

func UpdateNode(n *Node, id uint) {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
		return
	}
	defer db.Close()
	db.Model(n).Where("id=?", id).Update(n)
}

func DeleteNode(id uint) {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
		return
	}
	defer db.Close()
	n := Node{}
	db.Model(&n).Where("id=?", id).Delete(&n)
}

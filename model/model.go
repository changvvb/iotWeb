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

	db.CreateTable(&Node{})
	db.AutoMigrate(&Node{}, &Data{}, &Park{}, &Danger{})
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
	node.GetDanger()
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
	db.Model(n).Where("id=?", id).Update(n).Update("x", n.X).Update("y", n.Y).Update("min_value", n.MinValue).Update("max_value", n.MaxValue)
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

func DeletePark(id uint) {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
		return
	}
	defer db.Close()
	p := &Park{}
	p.ID = id
	db.Model(p).Delete(p)
}

func AddNode(n *Node) {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
		return
	}
	defer db.Close()
	db.Create(n)
}

func GetParks() []Park {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
		return nil
	}
	defer db.Close()
	var parks []Park
	db.Find(&parks)
	return parks
}

func GetParkByID(id uint) *Park {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
		return nil
	}
	defer db.Close()
	park := &Park{}
	db.Where("id=?", id).Find(park)
	return park
}

func AddPark(p *Park) {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
		return
	}
	defer db.Close()
	db.Create(p)
}

func GetDangers() map[string][]string {

	dangerMap := make(map[string][]string)
	var d []Danger
	db, err := opendb()
	if err != nil {
		printlnLog(err)
		return nil
	}
	defer db.Close()

	db.Find(&d)

	for _, v := range d {
		dangerMap[v.Species] = append(dangerMap[v.Species], v.Name)
	}

	return dangerMap
}

func GetDangerIDByString(name string) uint {
	db, err := opendb()
	if err != nil {
		printlnLog(err)
		return 0
	}
	defer db.Close()

	d := Danger{}
	db.Where("name=?", name).Find(&d)
	return d.ID
}

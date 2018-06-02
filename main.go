package main

import (
	"fmt"
	"iotWeb/model"
	"iotWeb/mqtt"
	"log"
	"strconv"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

var server *iris.Application

var (
	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
)

func main() {
	server = serverNew()
	serverSetup()
	mqttSetup()
	serverRun()
}

func serverNew() *iris.Application {
	return iris.New()
}

func serverSetup() {
	server.StaticServe("static", "/static")
	server.StaticServe("/js", "./static/js")

	temp := iris.HTML("./template", ".html")
	temp.Layout("base.html")
	temp.Reload(true)
	server.RegisterView(temp.Binary(nil, nil))

	server.Use(func(ctx iris.Context) {
		// path := ctx.PathString()
		path := ctx.Path()
		log.Println("request path:", path)
		if path != "/index" && path != "/logout" && path != "/login" && path != "/auth" && path != "/" && path != "/getallnodes" && path != "/parkinfo" && path != "/parklist" && path != "/parknodes" {
			session := sess.Start(ctx)
			if session.GetString("username") != "" {
				ctx.Next()
			} else {
				ctx.Redirect("/login")
			}
			return
		} else if 1 == 1 {

		}
		ctx.Next()
	})

	server.Get("/", func(ctx iris.Context) {
		ctx.Redirect("/index")
	})

	//主页
	server.Get("/index", func(ctx iris.Context) {
		// ctx.Render("index.html", struct{ Index bool }{true})
		ctx.ViewData("Index", true)
		ctx.View("index.html")
	})

	//管理员界面
	server.Get("/admin", func(ctx iris.Context) {
		ctx.ViewData("List", model.GetDangerSpeciesList())
		err := ctx.View("admin.html")
		checkError(err)
	})

	//管理员界面的json数据
	server.Get("/adminjson", func(ctx iris.Context) {
		ps := model.GetParks()
		for i, _ := range ps {
			ps[i].GetNodes()
		}
		ctx.JSON(ps)
	})

	//添加危险源
	server.Post("/adddanger", func(ctx iris.Context) {
		s := ctx.FormValue("species")
		n := ctx.FormValue("name")
		if s == "other" {
			s = ctx.FormValue("otherspecies")
		}
		log.Println(s, n)
		d := model.Danger{Species: s, Name: n}
		model.AddDanger(&d)
		ctx.Redirect("/admin")
	})

	//进入对应园区管理界面
	server.Get("/park/:id", func(ctx iris.Context) {
		// id, err := ctx.URLParamInt("id")
		id, err := ctx.Params().GetInt("id")
		checkError(err)

		park := model.GetParkByID(uint(id))
		type Park struct {
			model.Park
			OffLineNode   []*model.Node
			OnLineNodeMap []*model.Node
			Dangers       map[string][]string
		}

		p := &Park{}
		p.Park = *park
		p.OnLineNodeMap, p.OffLineNode = mqtt.GetNodes(park)
		p.Dangers = model.GetDangers()
		// ctx.ViewData("Park", p)
		// ctx.ViewData("Name", p.Name)
		// ctx.ViewData("", p)
		ctx.ViewData("OnLineNodeMap", p.OnLineNodeMap)
		ctx.ViewData("ID", p.ID)
		ctx.ViewData("Name", p.Name)
		ctx.ViewData("OffLineNode", p.OffLineNode)
		ctx.ViewData("Dangers", p.Dangers)
		log.Println(ctx.View("park.html"))
	})

	//获得所有节点
	server.Get("/getallnodes", func(ctx iris.Context) {
		nodes := model.GetNodes()
		ctx.JSON(nodes)
	})

	//登陆
	server.Get("/login", func(ctx iris.Context) {
		ctx.View("login.html")
	})

	server.Get("logout", func(ctx iris.Context) {
		session := sess.Start(ctx)
		session.Clear()
		ctx.Redirect("/index")
	})

	server.Post("/login", func(ctx iris.Context) {
		username := ctx.FormValue("username")
		password := ctx.FormValue("password")
		if (username == "changvvb" && password == "changvvb") || (username == "123456" && password == "123456") {
			log.Println("login success")
			sess.Start(ctx).Set("username", username)
			ctx.Redirect("/admin")
		} else {
			ctx.ViewData("LoginError", true)
			ctx.View("login.html")
		}
	})

	server.Get("/nodexy/:x/:y/:park", func(ctx iris.Context) {
		x, _ := ctx.Params().GetInt("x")
		y, _ := ctx.Params().GetInt("y")
		park, _ := ctx.Params().GetInt("park")
		log.Println(x, y)
		park = park
		id := model.GetIdByPosition(x, y, uint(park))
		// ctx.DirectTo("/node/" + fmt.Sprint(id))
		ctx.Redirect("/node/" + fmt.Sprint(id))
	})

	server.Get("/node/:id", func(ctx iris.Context) {
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			log.Println(err)
			return
		}
		type Node struct {
			model.Node
			ParkName string
			Dangers  map[string][]string
		}
		node := Node{}
		n := model.GetNodeByID(uint(id))
		if n == nil {
			log.Println("not found")
			ctx.StatusCode(iris.StatusNotFound)
			ctx.View("404.html")
			return
		}
		node.Node = *n
		node.ParkName = model.GetParkByID(node.ParkRefer).Name
		node.Dangers = model.GetDangers()

		ctx.ViewData("Node", node)
		ctx.ViewData("Danger", node.Danger)
		ctx.ViewData("ParkRefer", node.ParkRefer)
		ctx.ViewData("ParkName", node.ParkName)
		ctx.ViewData("MinValue", node.MinValue)
		ctx.ViewData("MaxValue", node.MaxValue)
		ctx.ViewData("X", node.X)
		ctx.ViewData("Y", node.Y)
		ctx.ViewData("Describe", node.Describe)
		ctx.ViewData("ID", node.ID)
		log.Println(ctx.View("nodeview.html"))
	})

	//节点历史界面
	server.Get("/nodehistory/:id", func(ctx iris.Context) {
		// id, err := ctx.URLParamInt("id")
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			return
		}
		node := model.GetNodeByID(uint(id))
		node.GetData()
		// log.Println(ctx.Render("nodehistory.html", node))
		log.Println(ctx.View("nodehistory.html"))
	})

	//节点历史json数据
	server.Get("/nodehistoryjson/:id", func(ctx iris.Context) {
		// id, err := ctx.URLParamInt("id")
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			return
		}
		node := model.GetNodeByID(uint(id))
		node.GetData()
		ctx.JSON(node.Data)
	})

	server.Get("/nodeseries/:id", func(ctx iris.Context) {
		ctx.Header("Content-type", "text/json")
		idInt, err := ctx.Params().GetInt("id")
		if err != nil {
			log.Println(err)
			return
		}
		id := uint(idInt)

		var r [2]interface{}
		r[0] = time.Now().Unix() * 1000

		if mqtt.OnLineNodeMap[id] == nil {
			r[1] = nil
			ctx.JSON(r)
			return
		}

		r[1] = mqtt.OnLineNodeMap[id].FreshData.Data
		ctx.JSON(r)
	})

	//修改节点
	server.Post("nodemodify/:id", func(ctx iris.Context) {
		danger := ctx.FormValue("danger")
		max := ctx.FormValue("max")
		min := ctx.FormValue("min")
		describe := ctx.FormValue("describe")
		x := ctx.FormValue("X")
		y := ctx.FormValue("Y")

		// ID, err := ctx.URLParamInt("id")
		id, err := ctx.Params().GetInt("id")
		checkError(err)
		Max, err := strconv.ParseFloat(max, 10)
		checkError(err)
		Min, err := strconv.ParseFloat(min, 10)
		checkError(err)
		X, err := strconv.ParseInt(x, 10, 64)
		log.Println("X:", int(X))
		checkError(err)
		Y, err := strconv.ParseInt(y, 10, 64)
		checkError(err)

		node := model.Node{
			// Species:  species,
			MaxValue: Max,
			MinValue: Min,
			Describe: describe,
			X:        int(X),
			Y:        int(Y),
			DangerID: model.GetDangerIDByString(danger),
		}
		model.UpdateNode(&node, uint(id))
		if n := mqtt.OnLineNodeMap[uint(id)]; n != nil {
			n.UpdateNode()
		}
		ctx.Redirect("/node/" + fmt.Sprint(id))

	})

	//增加一个节点
	server.Post("/nodeadd/:parkid", func(ctx iris.Context) {
		// danger := ctx.FormValueString("danger")
		// danger := ctx.FormValues("danger")
		danger := ctx.FormValue("danger")
		max := ctx.FormValue("max")
		min := ctx.FormValue("min")
		describe := ctx.FormValue("describe")
		x := ctx.FormValue("X")
		y := ctx.FormValue("Y")
		// id, err := ctx.URLParamInt("parkid")
		id, err := ctx.Params().GetInt("id")
		checkError(err)

		Max, err := strconv.ParseFloat(max, 10)
		checkError(err)
		Min, err := strconv.ParseFloat(min, 10)
		checkError(err)
		X, err := strconv.ParseInt(x, 10, 64)
		log.Println("X:", int(X))
		checkError(err)
		Y, err := strconv.ParseInt(y, 10, 64)
		checkError(err)

		node := model.Node{
			MaxValue:  Max,
			MinValue:  Min,
			Describe:  describe,
			X:         int(X),
			Y:         int(Y),
			ParkRefer: uint(id),
			DangerID:  model.GetDangerIDByString(danger),
		}
		model.GetParkByID(uint(id)).AddNode(&node)
		ctx.Redirect(fmt.Sprintf("/park/%d", id))
	})

	//删除一个节点
	server.Post("/delete/:id", func(ctx iris.Context) {
		// id, err := ctx.URLParamInt("id")
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("delete", id)
		pid := model.DeleteNode(uint(id))
		ctx.Redirect(fmt.Sprintf("/park/%d", pid))
	})

	server.Post("/deletepark/:id", func(ctx iris.Context) {
		// id, err := ctx.URLParamInt("id")
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Delete park", id)
		model.DeletePark(uint(id))
		ctx.Redirect("/admin")
	})

	server.Post("/addpark", func(ctx iris.Context) {
		name := ctx.FormValue("name")
		address := ctx.FormValue("address")
		tel := ctx.FormValue("tel")

		p := &model.Park{
			Name:    name,
			Address: address,
			Tel:     tel,
		}

		model.AddPark(p)
	})

	//给手机的,返回所有园区列表
	server.Get("/parklist", func(ctx iris.Context) {
		parks := model.GetParks()
		ctx.JSON(parks)
	})
	//给手机的,返回某个园区信息
	server.Get("/parkinfo", func(ctx iris.Context) {
		// idint, err := ctx.URLParamInt("id")
		idint, err := ctx.Params().GetInt("id")
		checkError(err)
		id := uint(idint)
		park := model.GetParkByID(id)
		park.GetNodes()
		ctx.JSON(park)
	})
	//给手机的返回节点信息
	server.Get("/parknodes", func(ctx iris.Context) {
		// idint, err := ctx.URLParamInt("id")
		idint, err := ctx.Params().GetInt("id")
		checkError(err)
		id := uint(idint)
		park := model.GetParkByID(id)
		park.GetNodes()
		ctx.JSON(park.Nodes)
	})
	//给手机的验证密码
	server.Get("/auth", func(ctx iris.Context) {
		name := ctx.URLParam("username")
		pass := ctx.URLParam("password")
		log.Println("/auth", name, pass)
		if name == "123456" && pass == "123456" {
			ctx.StatusCode(iris.StatusOK)
		} else {
			ctx.StatusCode(iris.StatusUnauthorized)
		}
	})
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
		//	panic(err)
	}
}

func serverRun() {
	server.Run(iris.Addr(":7070"))
}

func mqttSetup() {
	// mqtt.Subscribe("haha", 0)
}

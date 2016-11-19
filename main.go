package main

import (
	"fmt"
	"iotWeb/model"
	"iotWeb/mqtt"
	"log"
	"strconv"
	"time"

	"github.com/kataras/go-template/html"
	"github.com/kataras/iris"
)

var server *iris.Framework

func main() {
	server = serverNew()
	serverSetup()
	websocketsSetup()
	mqttSetup()
	serverRun()
}

func serverNew() *iris.Framework {
	return iris.New()
}

func serverSetup() {
	server.Config.IsDevelopment = true

	server.StaticServe("static", "/static")
	server.Static("/js", "./static/js", 1)
	server.UseTemplate(html.New(html.Config{
		Layout: "base.html",
	})).Directory("template", ".html")
	/*  server.UseFunc(func(ctx *iris.Context) { */
	//     path := ctx.PathString()
	//     if path != "/index" && path != "/logout" && path != "/login" && path != "/" && path != "/getallnodes" {
	//         if ctx.Session().GetString("username") != "" {
	//             ctx.Next()
	//         } else {
	//             ctx.Redirect("/login")
	//         }
	//         return
	//     }
	//     ctx.Next()
	// })

	server.Get("/", func(ctx *iris.Context) {
		ctx.Redirect("/index")
	})

	//主页
	server.Get("/index", func(ctx *iris.Context) {
		// ctx.MustRender("base.html", nil)
		ctx.Render("index.html", struct{ Index bool }{true})
	})

	//管理员界面
	server.Get("/admin", func(ctx *iris.Context) {
		Park := model.GetParks()
		for _, p := range Park {
			log.Println(p.Name)
		}
		err := ctx.Render("admin.html", struct{ Park []model.Park }{model.GetParks()})
		checkError(err)
	})

	//进入对应园区管理界面
	server.Get("/park/:id", func(ctx *iris.Context) {
		id, err := ctx.ParamInt("id")
		checkError(err)

		park := model.GetParkByID(uint(id))
		type Park struct {
			model.Park
			OffLineNode   []*model.Node
			OnLineNodeMap []*model.Node
		}

		p := &Park{}
		p.Park = *park
		p.OnLineNodeMap, p.OffLineNode = mqtt.GetNodes(park)
		log.Println(ctx.Render("park.html", p))
	})

	server.Get("/getallnodes", func(ctx *iris.Context) {
		nodes := model.GetNodes()
		ctx.JSON(iris.StatusOK, nodes)
	})

	server.Get("/login", func(ctx *iris.Context) {
		ctx.MustRender("login.html", nil)
	})

	server.Get("logout", func(ctx *iris.Context) {
		ctx.Session().Clear()
		ctx.Redirect("/index")
	})

	server.Post("/login", func(ctx *iris.Context) {
		username := ctx.FormValueString("username")
		password := ctx.FormValueString("password")
		if (username == "changvvb" && password == "changvvb") || (username == "123456" && password == "123456") {
			log.Println("login success")
			ctx.Session().Set("username", username)
			ctx.Redirect("/admin")
		} else {
			ctx.Render("login.html", struct{ LoginError bool }{true})
		}
	})

	server.Get("/nodexy/:x/:y", func(ctx *iris.Context) {
		x, _ := ctx.ParamInt("x")
		y, _ := ctx.ParamInt("y")
		log.Println(x, y)
		id := model.GetIdByPosition(x, y)
		// ctx.DirectTo("/node/" + fmt.Sprint(id))
		ctx.Redirect("/node/" + fmt.Sprint(id))
	})

	server.Get("/node/:id", func(ctx *iris.Context) {
		id, err := ctx.ParamInt("id")
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(id)
		node := model.GetNodeByID(uint(id))
		if node == nil {
			log.Println("not found")
			ctx.RenderWithStatus(iris.StatusNotFound, "404.html", nil)
			return
		}
		log.Println(ctx.Render("nodeview.html", node))
	})

	server.Get("/nodehistory/:id", func(ctx *iris.Context) {
		id, err := ctx.ParamInt("id")
		if err != nil {
			return
		}
		node := model.GetNodeByID(uint(id))
		node.GetData()
		log.Println(ctx.Render("nodehistory.html", node))
	})

	server.Get("/nodeseries/:id", func(ctx *iris.Context) {
		ctx.SetHeader("Content-type", "text/json")
		idInt, err := ctx.ParamInt("id")
		if err != nil {
			log.Println(err)
			return
		}
		id := uint(idInt)

		var r [2]interface{}
		r[0] = time.Now().Unix() * 1000

		if mqtt.OnLineNodeMap[id] == nil {
			r[1] = nil
			ctx.JSON(iris.StatusOK, r)
			return
		}

		r[1] = mqtt.OnLineNodeMap[id].FreshData.Data
		ctx.JSON(iris.StatusOK, r)
	})
	server.Post("nodemodify/:id", func(ctx *iris.Context) {
		species := ctx.FormValueString("species")
		max := ctx.FormValueString("max")
		min := ctx.FormValueString("min")
		describe := ctx.FormValueString("describe")
		x := ctx.FormValueString("X")
		y := ctx.FormValueString("Y")

		ID, err := ctx.ParamInt("id")
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
			Species:  species,
			MaxValue: Max,
			MinValue: Min,
			Describe: describe,
			X:        int(X),
			Y:        int(Y),
		}
		model.UpdateNode(&node, uint(ID))
		if n := mqtt.OnLineNodeMap[uint(ID)]; n != nil {
			n.UpdateNode()
		}
		ctx.Redirect("/node/" + fmt.Sprint(ID))

	})

	server.Post("/nodeadd/:parkid", func(ctx *iris.Context) {
		species := ctx.FormValueString("species")
		max := ctx.FormValueString("max")
		min := ctx.FormValueString("min")
		describe := ctx.FormValueString("describe")
		x := ctx.FormValueString("X")
		y := ctx.FormValueString("Y")
		id, err := ctx.ParamInt("parkid")
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
			Species:   species,
			MaxValue:  Max,
			MinValue:  Min,
			Describe:  describe,
			X:         int(X),
			Y:         int(Y),
			ParkRefer: uint(id),
		}
		model.GetParkByID(uint(id)).AddNode(&node)
		ctx.Redirect("/admin")
	})

	server.Post("/delete/:id", func(ctx *iris.Context) {
		id, err := ctx.ParamInt("id")
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("delete", id)
		model.DeleteNode(uint(id))
		ctx.Redirect("/admin")
	})

	server.Post("/deletepark/:id", func(ctx *iris.Context) {
		id, err := ctx.ParamInt("id")
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Delete park", id)
		model.DeletePark(uint(id))
		ctx.Redirect("/admin")
	})
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
		//	panic(err)
	}
}

var messageRoom = "001"

func websocketsSetup() {

	server.Config.Websocket.Endpoint = "/endpoint"
	server.Websocket.OnConnection(func(c iris.WebsocketConnection) {
		c.Join(messageRoom)
		c.On("chat", func(message string) {
			c.To(messageRoom).Emit("chat", "From: "+c.ID()+": "+message)
			log.Println(message)
		})
		c.OnDisconnect(func() {
			log.Printf("\nConnection with ID: %s	has	beendiscon	nected!", c.ID())
		})
	})
}

func serverRun() {
	server.Listen(":7070")
}

func mqttSetup() {
	// mqtt.Subscribe("haha", 0)
}

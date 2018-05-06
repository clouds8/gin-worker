package main

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/go-nats"
)

type Mission struct {
	Name     string `json:"mission"`
	Desc     string `json:"description"`
	Priority int    `json:"priority"`
}

const NATS = nats.DefaultURL // "nats://localhost:4222"

func main() {
	nc, err := nats.Connect(NATS)
	if err != nil {
		log.Println("err: ", err)
		return
	}
	nc.Subscribe("mission", func(m *nats.Msg) {
		var mission *Mission
		log.Println("[]byte string", string(m.Data))
		json.Unmarshal(m.Data, &mission)
		log.Println("receive a mission topic. ", mission)
	})

	r := gin.Default()
	r.GET("nats", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"nats": NATS,
		})
	})
	r.POST("mission", func(c *gin.Context) {
		//test the publishing
		var mission *Mission
		err := c.Bind(&mission)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "invalid input format",
			})
		}
		missionByte, err := json.Marshal(mission)
		err = nc.Publish("mission", missionByte)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "ok",
			})
		}

	})
	r.Run(":7000")
}

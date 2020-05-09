package procon_data

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

type Msg struct {
	Jwt string `json:"jwt"`
	Type string `json:"type"`
	Data string	`json:"data"`
}
//from interface
type TryUser struct {
	User     string     `json:"user"`
	Password string     `json:"password"`
}
//a user struct for testing user and pass.. not to send over network...

type Usertc struct {
	gorm.Model
	FullName string `gorm:"unique;not null;type:varchar(100);default:null"`
	User string `gorm:"unique;not null;type:varchar(100);default:null"`
	Coins uint64
	Role string
	Password string `gorm:"unique;not null;type:varchar(100);default:null"`
}

//Websocket Connections
type Client struct {
	Id string
	Conn *websocket.Conn
}
type PubSub struct {
	Clients []Client
}
var Ps = &PubSub{}
func (ps *PubSub) AddClient(client Client) (*PubSub) {

	ps.Clients = append(ps.Clients, client );

	return ps
}
func (ps *PubSub) RemoveClient(client Client) (*PubSub) {
	for index, c := range ps.Clients {
		if c.Id == client.Id {
			ps.Clients = append(ps.Clients[:index], ps.Clients[index+1:]...)
		}
	}
	return ps
}
func (ps *PubSub) PublishAll(jd string) {
	for index, c := range ps.Clients {
		_ = index
		m := Msg{"^vAr^","client-list", jd};
		if err := c.Conn.WriteJSON(m); err != nil {
			fmt.Println(err)
		}
	}
}
//SSE Data
type SseUpdate struct {
	Type string
	Data []map[string]interface{}
}

var Su = &SseUpdate{
	Type: "SseUpdate",
	Data: make([]map[string]interface{}, 0, 0),
}
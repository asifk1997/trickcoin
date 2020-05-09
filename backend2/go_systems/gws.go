package main

import (
	"./procon_asyncq"
	"./procon_config"
	"./procon_db"
	"./procon_jwt"
	"./procon_utils"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

var addr = flag.String("addr", "0.0.0.0:1200", "http service address")
var upgrader = websocket.Upgrader{}

//the jwt for messages out can be used for anything or obscure data...
//It is in the msg struct for protected messages coming in containing a jwt token
type msg struct {
	Jwt string `json:"jwt"`
	Type string `json:"type"`
	Data string	`json:"data"`
}


func sendMsg(j string, t string, d string, c *websocket.Conn) {
	m := msg{j, t, d};
	if err := c.WriteJSON(m); err != nil {
		fmt.Println(err)
	}
	//mm, _ := json.Marshal(m);
	//fmt.Println(string(mm));
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request ) bool { return true}
	c, err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		fmt.Print("WTF @HandleAPI Ws Upgrade Error> ", err)
		return
	}
Loop:
	for {
		//fmt.Printf("for loop running")
		in := msg{}

		err := c.ReadJSON(&in)
		if err != nil {
			//fmt.Println("Error reading json.", err)
			c.Close()

			break Loop

		}
		switch (in.Type) {
		case "get-jwt-token":
			fmt.Println(in.Data); //DEBUG REDACTED
			usr, pwd, err := procon_utils.B64DecodeTryUser(in.Data);
			if err != nil { fmt.Println(err);  } else {  fmt.Println(string(usr), string(pwd))  }

			upv, auser, err := procon_db.MongoTryUser(usr,pwd)
			//if err != nil {
			//	fmt.Println(err)
			//} else {
			//	fmt.Println(upv)
			//	fmt.Println(auser)
			//	jauser, err := json.Marshal(auser);
			//	if err != nil {
			//		fmt.Println("error marshaling AUser.")
			//	} else {
			//		//jwt, err := procon_jwt.GenerateJWT(procon_config.PrivKeyFile);
			//		fmt.Println(string(jauser))
			//		_ = err
			//	}
			//}


			if err != nil { fmt.Println(err); sendMsg("noop", "invalid-credentials","noop", c); } else {
				if upv == true { fmt.Println("A user has logged in."); }
				auser.Password = "F00"
				jauser,err := json.Marshal(auser); if err != nil { fmt.Println("error marshaling AUser.") } else {
					jwt, err := procon_jwt.GenerateJWT(procon_config.PrivKeyFile);
					if err != nil { fmt.Println(err);  } else  { sendMsg(jwt, "jwt-token", string(jauser), c);  }
				}
			}

		case "verify-jwt-token": fallthrough
		case "validate-stored-jwt-token":
			valid, err := procon_jwt.ValidateJWT(procon_config.PubKeyFile,in.Jwt)
			if err != nil { fmt.Println(err); sendMsg("^vAr^", "jwt-token-invalid",err.Error(), c) } else if (err == nil && valid ) {
				if in.Type == "verify-jwt-token" { sendMsg("^vAr^", "jwt-token-valid","noop", c) }
				if in.Type == "validate-stored-jwt-token" {  sendMsg("^vAr^", "stored-jwt-token-valid","noop", c) }
			}
			break;
		case "rapid-test-user-avail":
			tobj := procon_db.NewRapidTestUserAvailTask(in.Data, c);
			procon_asyncq.TaskQueue <- tobj
			break;
		case "create-user":
			tobj := procon_db.NewCreateUserTask(in.Data, c);
			procon_asyncq.TaskQueue <- tobj
			break;

		default:
			break

		}
	}
}

func main() {

	db:=procon_db.GetDB()
	procon_asyncq.StartTaskDispatcher(9)
	flag.Parse()
	log.SetFlags(0)

	r:= mux.NewRouter()
	// websocket api
	r.HandleFunc("/api",handleAPI)

	http.ListenAndServe(*addr, r)



	defer db.Close()
}

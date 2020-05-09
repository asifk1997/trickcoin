package procon_db

import (
	"../procon_data"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	//"go.mongodb.org/mongo-driver/bson"
	"../procon_utils"
	//"go.mongodb.org/mongo-driver/bson"
	"os"
)

var db *gorm.DB //database


func init() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")


	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&procon_data.Usertc{}) //Database migration
}


//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}



func MongoTryUser(u []byte, p []byte) (bool,*procon_data.Usertc,error) {
	var xdoc procon_data.Usertc
	err:=db.Table("usertcs").Where(`"user" = ?`, string(u)).First(&xdoc).Error
	if err!=nil{
		return false,nil,err
	}else{
		bres,err := procon_utils.ValidateUserPassword(p,[]byte(xdoc.Password))
		if err!=nil { return false,nil,err} else { return bres,&xdoc,nil }
	}
}

/* Create User Code */
type CreateUserTask  struct {
	jsonuserstr string
	ws *websocket.Conn
}
func NewCreateUserTask(jsonuserstr string, ws *websocket.Conn)  *CreateUserTask {
	return	&CreateUserTask {jsonuserstr, ws}
}
func (cut *CreateUserTask) Perform() {
	user := procon_data.Usertc{}
	err := json.Unmarshal([]byte(cut.jsonuserstr), &user);
	if err != nil {  fmt.Println(err) }else {
		//fmt.Printf("%+v\n", user) DEBUG REDACTED

		//collection := client.Database("api").Collection("users");
		//check if user already exists

		var xdoc interface{}
		//filter := bson.D{{"user", user.User }}
		//err := collection.FindOne(ctx, filter).Decode(&xdoc);
		err:=db.Table("usertcs").Where(`"user" = ?`, user.User).First(&xdoc).Error
		if (err != nil && xdoc == nil) {
			fmt.Println("User Available", err);

			hp := procon_utils.GenerateUserPassword(user.Password);
			user.Password = hp;
			user.Coins =100;
			//insertResult, err := collection.InsertOne(ctx, &user)
			err:=db.Create(&user).Error
			if err != nil { fmt.Println("Error Inserting Document"); } else {
				fmt.Println("Inserted a single User: ", &user.User)
				procon_utils.SendMsg("vAr","toast-success", "User Created Successfully", cut.ws);
			}
		} else {
			//shouldn't get here but it means some how rapid test didn't catch this
			//modal is still open so just display modal error...
			procon_utils.SendMsg("vAr","rapid-test-user-avail-fail", "User Already Exists!", cut.ws);
		}
	}
}

type RapidTestUserAvailTask struct {
	rtu string
	ws *websocket.Conn
}
func NewRapidTestUserAvailTask(rtu string, ws *websocket.Conn)  *RapidTestUserAvailTask {
	return	&RapidTestUserAvailTask{rtu, ws}
}

func (rtuat *RapidTestUserAvailTask) Perform() {
	var xdoc procon_data.Usertc

	fmt.Println("rapidtestuser",rtuat.rtu);
	err:=db.Table("usertcs").Where(`"user" = ?`, rtuat.rtu).First(&xdoc).Error
	fmt.Println(xdoc)
	if (err != nil && xdoc.User != rtuat.rtu) {
		procon_utils.SendMsg("vAr","rapid-test-user-avail-success", "noop", rtuat.ws);
	} else {
		procon_utils.SendMsg("vAr","rapid-test-user-avail-fail", "noop", rtuat.ws);
	}
}

package procon_config

import (
	"../../go_systems/procon_fs"
	"crypto/rsa"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	//"procon_fs"
	//"awesomeProject/go_systems/procon_jwt"
)

var (
	PubKeyFile	*rsa.PublicKey
	PrivKeyFile *rsa.PrivateKey
)
//these can also be set...
const (
	PKPWD = "SOMEHARDPASSWORD"
	DirPath = "D:\\Asif\\Old PRojects\\wab_dev\\backend2\\"
	KeyCertPath = DirPath+"keycertz\\"
	PrivKeyPath = DirPath+"keycertz\\mykey.pem"
	PubKeyPath  = DirPath+"keycertz\\mykey.pub"

	//dont forget to escape characters like @ w/ %40
	MongoHost = "127.0.0.1"
	MongoUser = "mongod"
	MongoPassword = "SOMEHARDPASSWORD"
	MongoDb = "admin"

	RedisRP = "SOMELOGASSPASSWORD HERE"

	MysqlPass = "ANOTHER-HARD-PASSOWRD"
)


func init() {
	f,ok,err := procon_fs.ReadFile(PubKeyPath)
	if (!ok || err != nil) { fmt.Println(err) } else {
		//PubKeyFile, err = procon_jwt.ParseRSAPublicKeyFromPEM(f)
		PubKeyFile, err = jwtgo.ParseRSAPublicKeyFromPEM(f)
		if err != nil { fmt.Println(err) }
	}
	f,ok,err = procon_fs.ReadFile(PrivKeyPath)
	if (!ok || err != nil) { fmt.Println(err) } else {
		//PrivKeyFile, err = procon_jwt.ParseRSAPrivateKeyFromPEMWithPassword(f, PKPWD)
		PrivKeyFile, err = jwtgo.ParseRSAPrivateKeyFromPEMWithPassword(f, PKPWD)
		if err != nil { fmt.Println(err) }
	}
}


package lib

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoClient struct {
	IP      string
	Port    int
	Session *mgo.Session
}

var (
	g_MongoClient *MongoClient = nil
)

func InitMongo(ip string, port int) {
	session, err := mgo.Dial(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	fmt.Printf("Connect Mgo %s:%d Succ\n", ip, port)
	g_MongoClient = new(MongoClient)
	g_MongoClient.IP = ip
	g_MongoClient.Port = port
	g_MongoClient.Session = session
}

func GetMongoConn() *mgo.Session {
	if g_MongoClient == nil {
		panic("Mgo Client Not inited")
	}
	return g_MongoClient.Session
}

func GetMongoDB(db string) *mgo.Database {
	obj := GetMongoConn()
	return obj.DB(db)
}

func GetMongoCollect(db string, c string) *mgo.Collection {
	obj := GetMongoConn()
	return obj.DB(db).C(c)
}

func NewObjectIdHex() string {
	return bson.NewObjectId().Hex()
}

func NewObjectIdObj() bson.ObjectId {
	return bson.NewObjectId()
}

package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	insert()
}

func testConnection() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://192.168.10.4:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}

func insert() {
	// 示例 MongoDB 连接信息
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://192.168.10.4:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// 选择数据库和集合
	collection := client.Database("minigame").Collection("test_minigame_game_user")

	// 示例的键和 Unix 时间戳字符串
	key := "ttlTime"
	unixString := "1715913470"

	// 将 Unix 时间戳字符串转换为 int64 类型
	ttlValue, err := strconv.ParseInt(unixString, 10, 64)
	if err != nil {
		logrus.Error("string to int64 error")
	}

	// 将 int64 类型的时间戳转换为 time.Time 类型
	ttlTime := time.Unix(ttlValue, 0)

	// 创建要插入的 BSON 文档
	bsonDoc := bson.M{}
	bsonDoc[key] = ttlTime

	// 插入文档到 MongoDB
	_, err = collection.InsertOne(ctx, bsonDoc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted document with TTL time:", bsonDoc)
}

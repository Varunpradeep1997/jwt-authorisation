package database

import(
	"github.com/joho/godotenv"
	"fmt"
	"log"
	"time"
	"os"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
	err:=godotenv.Load(".env")           //here load func is used to boot up the environment
	if err != nil{
		log.Fatal("error loading .env file")
	}

	MongoDb :=os.Getenv("MONGODB_URL")   //initialized a variable called mongodb.
client,err:=	mongo.NewClient(options.Client().ApplyURI(MongoDb))
if err != nil{
	log.Fatal(err)
}
ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
defer cancel()
err=client.Connect(ctx)
if err != nil{
	log.Fatal(err)
}
fmt.Println("connected to Mongodb!") 
return client

}
var Client *mongo.Client=DBinstance()

func OpenCollection(client *mongo.Client,collectionName string)*mongo.Collection{
	var collection *mongo.Collection=client.Database("cluster0").Collection(collectionName)
	return collection

}
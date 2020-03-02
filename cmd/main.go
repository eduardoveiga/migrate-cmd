package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"

	_ "../migrations"
	migrate "github.com/xakep666/mongo-migrate"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if len(os.Args) == 1 {
		logrus.Fatal("Missing options: up or down")
	}
	option := os.Args[1]
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// Connect to MongoDB
	client, err := mgo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.Info(err.Error())
	}
	err = nil
	migrate.SetDatabase(client.Database("main"))
	migrate.SetMigrationsCollection("migrations")

	switch option {
	case "new":
		if len(os.Args) != 3 {
			fmt.Println(len(os.Args))
			logrus.Fatal("Should be: new <description of migration>")
		}
		fName := fmt.Sprintf("../migrations/%s_%s.go", time.Now().Format("200601021504"), os.Args[2])
		from, err := os.Open("../migrations/template.go")
		if err != nil {
			logrus.Fatal("Should be: new <description of migration>")
		}
		defer from.Close()
		to, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			logrus.Fatal(err.Error())
		}

		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			logrus.Fatal(err.Error())
		}
		logrus.WithFields(logrus.Fields{
			"file": fName,
		}).Info("New migration created")
	case "up":
		err = migrate.Up(migrate.AllAvailable)

	case "down":
		err = migrate.Down(migrate.AllAvailable)

	}
	if err != nil {
		logrus.Fatal(err.Error())
	}
}

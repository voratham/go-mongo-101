package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `json:"name"`
	Email     string             `json:email`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

func main() {
	client := dbConnection()
	person := Person{
		Name:      "Voratham",
		Email:     "voratham_sir@dev.io",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	insertedPerson1 := CreatePerson(client, person)
	oid, err := primitive.ObjectIDFromHex(insertedPerson1)
	if err != nil {
		log.Fatal((err))
	}

	updatedPerson1 := UpdateOnePerson(client, bson.M{"email": "voratham_sir@true-e-logistics.com"}, bson.M{"_id": oid})
	log.Println(updatedPerson1)
	personUpdated := GetPersonById(client, bson.M{"_id": oid})
	log.Println("personUpdate ::", personUpdated)

	people := GetPersonAll(client, bson.M{})
	for _, person := range people {
		fmt.Println("ID :: ", person.ID, " Name :: ", person.Name, " Email :: ", person.Email)
	}

	personLast := people[len(people)-1]
	DeletePersonById(client, bson.M{"_id": personLast.ID})
	personDeleted := GetPersonById(client, bson.M{"_id": personLast.ID})
	fmt.Println("personDeleted :: ", personDeleted)

}

func dbConnection() *mongo.Client {
	clinetOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clinetOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB !!!!!!!")
	return client
}

func CreatePerson(client *mongo.Client, person Person) string {
	collection := client.Database("user").Collection("user")
	res, err := collection.InsertOne(context.TODO(), person)
	if err != nil {
		log.Fatalln("Error on inserted new person", err)
	}
	return res.InsertedID.(primitive.ObjectID).Hex()
}

func GetPersonAll(client *mongo.Client, filter bson.M) []*Person {
	var people []*Person
	collection := client.Database("user").Collection("user")
	res, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	for res.Next(context.TODO()) {
		var person Person
		err := res.Decode(&person)
		if err != nil {
			log.Fatal("Error on Deconding the document ", err)
		}
		people = append(people, &person)
	}
	return people
}

func GetPersonById(client *mongo.Client, filter bson.M) Person {
	var person Person
	collection := client.Database("user").Collection("user")
	collection.FindOne(context.TODO(), filter).Decode(&person)
	return person

}

func DeletePersonById(client *mongo.Client, filter bson.M) int64 {
	collection := client.Database("user").Collection("user")
	deleteRes, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error on deleteing one person :: ", err)
	}
	return deleteRes.DeletedCount
}

func UpdateOnePerson(client *mongo.Client, updateData interface{}, filter bson.M) int64 {
	collection := client.Database("user").Collection("user")
	docUpdate := bson.D{{Key: "$set", Value: updateData}}
	log.Println()
	updateRes, err := collection.UpdateOne(context.TODO(), filter, docUpdate)
	if err != nil {
		log.Fatal(err)
	}
	return updateRes.ModifiedCount

}

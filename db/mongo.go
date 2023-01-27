package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	// BLOB
	Image string `json:"img"`
}

func Connect() (*mongo.Client, context.CancelFunc) {
	// Établir une connexion à la base de données MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println("Début connexion")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Fin connexion")
	return client, cancel
}

func GetLivresCollection(client *mongo.Client) *mongo.Collection {
	// Obtenir une référence à la collection "categories"
	return client.Database("mydb").Collection("livres")
}

func PostLivre(ctx context.Context, client *mongo.Client, livre Book) error {

	collection := GetLivresCollection(client)

	//newId := bsongen.NewObjectId()
	// Insérez le document dans la collection

	_, got := GetLivre(ctx, client, livre.ID)
	if got == false {
		_, err := collection.InsertOne(ctx, livre)
		if err != nil {
			panic(err)
		}
		return nil

	}

	return fmt.Errorf("already inserted")

}

// sortir un livre en particulier de la collection

func GetLivres(ctx context.Context, client *mongo.Client) ([]Book, error) {
	//collection := GetLivresCollection(client)
	coll := client.Database("mydb").Collection("livres")

	var result []Book

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for cursor.Next(ctx) {
		if err := cursor.All(ctx, &result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func GetLivre(ctx context.Context, client *mongo.Client, id string) (*Book, bool) {
	//collection := GetLivresCollection(client)
	coll := client.Database("mydb").Collection("livres")

	filter := bson.M{"id": id}

	var result []Book
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var livre bson.M
		if err = cursor.Decode(&livre); err != nil {
			log.Fatal(err)
		}
		fmt.Println(livre)

		var b Book
		bsonBytes, _ := bson.Marshal(livre)
		bson.Unmarshal(bsonBytes, &b)
		result = append(result, b)
	}

	if len(result) < 1 {
		return nil, false
	}

	return &result[0], true
}

func DeleteBook(ctx context.Context, id string, client *mongo.Client) error {

	// Construction du filtre pour supprimer l'objet par son ID
	collection := GetLivresCollection(client)
	filter := bson.M{"id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Deleted %d ", result.DeletedCount))

	return nil
}

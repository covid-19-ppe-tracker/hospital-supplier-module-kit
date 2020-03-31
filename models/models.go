package models

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	//"bitbucket.org/arorankit/getdone-kit/getdone-log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrNoMongoConnection is returned when there is no mongo connection
var ErrNoMongoConnection = errors.New("Mongo client is not connected")

type mongoclient struct {
	sync.Mutex
	c           *mongo.Client
	isConnected bool
	isStarted   bool
}

var client = mongoclient{}

//Connect connects and sets a connection.
//More info on mongo connection string
// https://docs.mongodb.com/manual/reference/connection-string/
func Connect(connectionURI string) error {
	if !client.isStarted {
		client.Lock()
		defer client.Unlock()
		if !client.isStarted {
			go func() {
				for {
					isConnectedOld := client.isConnected
					if !client.isConnected {
						connect(connectionURI)
					}
					if client.c != nil {
						err := ping()
						//logger, _ := getdone_log.GetLogger()
						if err != nil {
							client.isConnected = false
							client.c.Disconnect(context.Background())
							client.c = nil
							if isConnectedOld != client.isConnected {
								//log only when connection status changes
								//logger.Errorf("cannot create mongo session: %s\n", err)
							}
						} else {
							client.isConnected = true
							if isConnectedOld != client.isConnected {
								//log only when connection status changes
								//logger.Info("Connected to mongo at " + connectionURI)
							}
						}
					} else {
						client.isConnected = false
					}
					time.Sleep(5 * time.Second)
				}
			}()
			client.isStarted = true
		}
	}
	return nil
}

func ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return client.c.Ping(ctx, readpref.Primary())
}

func connect(connectionURI string) error {
	if !client.isConnected {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		c, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
		if err != nil {
			client.isConnected = false
			return err
		}
		client.c = c
	}
	return nil
}

//Disconnect mongo client connection.
//Returns ErrNoMongoConnection as error when no mongo connection
func Disconnect() error {
	if !client.isConnected {
		return ErrNoMongoConnection
	}
	client.c.Disconnect(context.Background())
	client.isConnected = false
	return nil
}

//GetClient returns the connected mongo client
//Returns ErrNoMongoConnection as error when no mongo connection
func GetClient() (*mongo.Client, error) {
	if !client.isConnected {
		return nil, ErrNoMongoConnection
	}
	return client.c, nil
}

type JSONSerializer interface {
	SerializeToJSON(w http.ResponseWriter) error
}

type JSONBodyDeserializer interface {
	DeserializeFromJSONInBody(r *http.Request) JSONBodyDeserializer
}

type collectionNamer interface {
	collectionName() string
}

type databaseNamer interface {
	databaseName() string
}

type collectionDatabaseNamer interface {
	collectionNamer
	databaseNamer
}

type authdb struct {
	databaseNamer
}

func (a authdb) databaseName() string {
	return "auth"
}

type b2cdb struct {
	databaseNamer
}

func (b b2cdb) databaseName() string {
	return "b2c"
}

//FindOne finds one entry in a collection based on mongo query
//Returns ErrNoMongoConnection as error when no mongo connection
func FindOne(m collectionDatabaseNamer, i interface{}, d bson.D) error {
	if !client.isConnected {
		return ErrNoMongoConnection
	}
	c := client.c.Database(m.databaseName()).Collection(m.collectionName())
	result := c.FindOne(nil, d)
	err := result.Decode(i)
	if err != nil {
		return err
	}
	return nil
}

//UpdateOne finds one entry in a collection based on mongo query and updates it
//Returns ErrNoMongoConnection as error when no mongo connection
func UpdateOne(m collectionDatabaseNamer, filter bson.D, update bson.D) (*mongo.UpdateResult, error) {
	if !client.isConnected {
		return nil, ErrNoMongoConnection
	}
	c := client.c.Database(m.databaseName()).Collection(m.collectionName())
	result, err := c.UpdateOne(nil, filter, bson.D{{"$set", update}})
	if err != nil {
		return nil, err
	}
	return result, nil
}

//DeleteOne delete one entry in a collection based on mongo query
//Returns ErrNoMongoConnection as error when no mongo connection
func DeleteOne(m collectionDatabaseNamer, d bson.D) error {
	if !client.isConnected {
		return ErrNoMongoConnection
	}
	c := client.c.Database(m.databaseName()).Collection(m.collectionName())
	dr, err := c.DeleteOne(nil, d)
	if err != nil {
		return err
	}
	if dr.DeletedCount == 1 {
		return nil
	}
	return errors.New("did not delete exactly one document")
}

//InsertOne inserts one entry in a collection based on mongo query.
//Returns ErrNoMongoConnection as error when no mongo connection
func InsertOne(m collectionDatabaseNamer, i interface{}) (interface{}, error) {
	if !client.isConnected {
		return nil, ErrNoMongoConnection
	}
	c := client.c.Database(m.databaseName()).Collection(m.collectionName())
	ir, err := c.InsertOne(nil, i)
	if err != nil {
		return nil, err
	}
	return ir.InsertedID, nil
}

//AuthenticateWithJWT authenticates user with jwt token
func AuthenticateWithJWT(id interface{}, r *http.Request) (*User, bool, error) {
	_id := id.(string)
	objID, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, false, err
	}

	user := &User{}
	err = FindOne(userModel, user, bson.D{{Key: "_id", Value: objID}})
	if err != nil {
		return nil, false, err
	}
	return user, true, nil
}

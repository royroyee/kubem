package k8s

import (
	cm "github.com/royroyee/kubem/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
	"time"
)

func (kh K8sHandler) DBSession() {

	log.Println("Success to Create DB Session")
	// already created db session (main/initHandlers)
	defer kh.session.Close()

}

// Create MongoDB Session
func GetDBSession() *mgo.Session {
	log.Println("Create DB Session .. ")
	//session, err := mgo.Dial("mongodb://db-service:27017") // db-service is name of mongodb service(kubernetes)
	session, err := mgo.Dial("mongodb://localhost:27017")

	if err != nil {
		panic(err)
	}
	return session
}

func (kh K8sHandler) GetEvents(eventLevel string, page int, perPage int) ([]cm.Event, error) {
	var result []cm.Event
	collection := kh.session.DB("kargos").C("event")

	skip := (page - 1) * perPage
	limit := perPage
	filter := bson.M{"eventlevel": strings.Title(eventLevel)}
	if eventLevel == "" {
		filter = bson.M{}
	}
	err := collection.Find(filter).Skip(skip).Limit(limit).Sort("-created").All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (kh K8sHandler) StoreEventInDB(event cm.Event) {

	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()

	collection := cloneSession.DB("kargos").C("event")

	err := collection.Insert(event)
	if err != nil {
		log.Println(err)
		return
	}
}

// Delete all event data older than 24 hours
func (kh K8sHandler) deleteEventFromDB() {
	collection := kh.session.DB("kargos").C("event")

	cutoff := time.Now().Add(-24 * time.Minute)
	_, err := collection.RemoveAll(bson.M{"timestamp": bson.M{"$lte": cutoff}})
	if err != nil {
		log.Println(err)
		return
	}
}

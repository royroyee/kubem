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
	collection := kh.session.DB("kubem").C("event")

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

func (kh K8sHandler) NumberOfEvents(eventLevel string) (cm.Count, error) {
	var result cm.Count
	filter := bson.M{}
	collection := kh.session.DB("kubem").C("event")

	if eventLevel != "" {
		filter = bson.M{"eventlevel": strings.Title(eventLevel)}
	}
	count, err := collection.Find(filter).Count()
	if err != nil {
		log.Println(err)
		return result, err
	}
	result.Count = count
	return result, nil
}

func (kh K8sHandler) StoreEventInDB(event cm.Event) {

	// Use its own session to avoid any concurrent use issues
	cloneSession := kh.session.Clone()

	collection := cloneSession.DB("kubem").C("event")

	err := collection.Insert(event)
	if err != nil {
		log.Println(err)
		return
	}
}

// Delete all event data older than 24 hours
func (kh K8sHandler) deleteEventFromDB() {
	collection := kh.session.DB("kubem").C("event")

	cutoff := time.Now().Add(-24 * time.Minute)
	_, err := collection.RemoveAll(bson.M{"timestamp": bson.M{"$lte": cutoff}})
	if err != nil {
		log.Println(err)
		return
	}
}

func (kh K8sHandler) GetNodeUsageAvg() (cm.NodeUsage, error) {
	var result cm.NodeUsage

	collection := kh.session.DB("kubem").C("node")

	// Aggregate the average value of cpuusage and ramusage per minute
	pipeline := collection.Pipe([]bson.M{
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"minute": bson.M{"$minute": bson.M{"$toDate": "$timestamp"}},
				},
				"avgCpuUsage": bson.M{"$avg": "$cpuusage"},
				"avgRamUsage": bson.M{"$avg": "$ramusage"},
			},
		},
		{"$limit": 24},
	})

	// Extract the result
	var getUsage []bson.M
	err := pipeline.All(&getUsage)
	if err != nil {
		log.Println(err)
		return result, err
	}

	for _, usage := range getUsage {
		avgCpuUsage := int(usage["avgCpuUsage"].(float64))
		result.CpuUsage = append(result.CpuUsage, avgCpuUsage)

		avgRamUsage := int(usage["avgRamUsage"].(float64))
		result.RamUsage = append(result.RamUsage, avgRamUsage)

	}

	return result, nil
}

func (kh K8sHandler) GetNodeOverview(page int, perPage int) ([]cm.NodeOverview, error) {
	var result []cm.NodeOverview

	// Get a reference to the "node" collection
	collection := kh.session.DB("kubem").C("node")

	// Define the pipeline stages
	pipeline := []bson.M{
		{"$sort": bson.M{"timestamp": -1}},
		{"$group": bson.M{
			"_id":      "$name",
			"name":     bson.M{"$first": "$name"},
			"cpuusage": bson.M{"$first": "$cpuusage"},
			"ramusage": bson.M{"$first": "$ramusage"},
			"ip":       bson.M{"$first": "$ip"},
			"status":   bson.M{"$first": "$status"},
		}},
		{"$skip": (page - 1) * perPage},
		{"$limit": perPage},
	}

	// Execute the query and get the results
	//err := collection.Find(query).Sort(sort).Skip(skip).Limit(limit).All(&result)
	err := collection.Pipe(pipeline).All(&result)
	if err != nil {
		log.Printf("error querying database: %s", err)
		return result, err
	}

	return result, nil
}

func (kh K8sHandler) GetNodeUsage(nodeName string) (cm.NodeUsage, error) {
	var result cm.NodeUsage

	collection := kh.session.DB("kubem").C("node")

	pipeline := collection.Pipe([]bson.M{
		{"$match": bson.M{"name": nodeName}},
		{"$limit": 24},
		{"$project": bson.M{
			"_id":      nil,
			"cpuusage": 1,
			"ramusage": 1,
		}},
	})

	// Extract the result
	var getUsage []bson.M
	err := pipeline.All(&getUsage)
	if err != nil {
		log.Println(err)
		return result, err
	}
	for _, usage := range getUsage {
		CpuUsage := int(usage["cpuusage"].(float64))
		result.CpuUsage = append(result.CpuUsage, CpuUsage)

		RamUsage := int(usage["ramusage"].(float64))
		result.RamUsage = append(result.RamUsage, RamUsage)
	}
	return result, nil
}

func (kh K8sHandler) GetVolumesOfController(namespace string, name string) (cm.ControllerDetail, error) {
	var result cm.ControllerDetail
	collection := kh.session.DB("kubem").C("controller")

	filter := bson.M{"namespace": namespace, "name": name}

	err := collection.Find(filter).One(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (kh K8sHandler) GetControllersByFilter(namespace string, controller string, page int, perPage int) ([]cm.ControllerOverview, error) {
	var result []cm.ControllerOverview
	collection := kh.session.DB("kubem").C("controller")

	skip := (page - 1) * perPage
	limit := perPage
	var filter bson.M
	if namespace != "" && controller != "" {
		filter = bson.M{
			"namespace": namespace,
			"type":      controller,
		}
	} else if namespace != "" {
		filter = bson.M{
			"namespace": namespace,
		}
	} else if controller != "" {
		filter = bson.M{
			"type": controller,
		}
	}
	err := collection.Find(filter).Skip(skip).Limit(limit).All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func (kh K8sHandler) NumberOfControllers(namespace string, controllerType string) (cm.Count, error) {
	var result cm.Count
	filter := bson.M{}
	collection := kh.session.DB("kubem").C("controller")

	if namespace != "" && controllerType == "" {
		filter = bson.M{"namespace": namespace}
	} else if namespace != "" && controllerType != "" {
		filter = bson.M{"namespace": namespace, "type": controllerType}
	}
	count, err := collection.Find(filter).Count()
	if err != nil {
		log.Println(err)
		return result, err
	}
	result.Count = count
	return result, nil
}

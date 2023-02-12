package input

import (
	"context"
	"sync"
	"time"

	"git.fiblab.net/sim/sidecar/core/downloader"
	agentv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/agent/v2"
	mapv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/map/v2"
	inputv1 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/traffic/input/v1"
	traffic_lightv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/traffic_light/v2"
	"git.fiblab.net/sim/simulet-go/utils/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 下载数据
func Init(config config.Config, cacheDir string) (res *inputv1.Input) {
	useCache := preCheckCache(cacheDir)

	// 连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Mongo.URI))
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		client.Disconnect(context.Background())
	}()
	// 初始化返回值
	res = &inputv1.Input{
		Persons:       &agentv2.Agents{},
		TrafficLights: &traffic_lightv2.TrafficLights{},
	}

	var wg sync.WaitGroup

	mapCol := downloader.GetMongoCol(client, config.Mongo.Map)
	if useCache {
		res.Map = mustReadCacheOrDownloadAndSave(
			cacheDir, config.Mongo.Map,
			func() *mapv2.Map {
				m := downloadMap(mapCol)
				return m
			},
		)
	} else {
		res.Map = downloadMap(mapCol)
	}
	log.Printf("finish fetching map from %s.%s at %s",
		config.Mongo.Map.DB, config.Mongo.Map.Col, config.Mongo.URI)

	personCol := downloader.GetMongoCol(client, config.Mongo.Person)
	if useCache {
		res.Persons = mustReadCacheOrDownloadAndSave(
			cacheDir, config.Mongo.Person,
			func() *agentv2.Agents {
				return downloadAgents(personCol)
			},
		)
	} else {
		res.Persons = downloadAgents(personCol)
	}
	log.Printf("finish fetching persons from %s.%s at %s",
		config.Mongo.Person.DB, config.Mongo.Person.Col, config.Mongo.URI)
	if len(res.Persons.Agents) == 0 {
		log.Fatal("no valid agents to simulate")
	}

	if config.Mongo.TrafficLight != nil {
		tlCol := downloader.GetMongoCol(client, config.Mongo.TrafficLight)
		wg.Add(1)
		go func() {
			defer wg.Done()
			res.TrafficLights = downloadTrafficLights(tlCol)
			log.Printf("finish fetching traffic_lights from %s.%s at %s",
				config.Mongo.TrafficLight.DB, config.Mongo.TrafficLight.Col, config.Mongo.URI)
		}()
	}
	wg.Wait()
	return
}

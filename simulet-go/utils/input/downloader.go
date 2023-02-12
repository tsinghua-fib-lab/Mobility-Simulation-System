package input

import (
	"context"

	"git.fiblab.net/sim/sidecar/core/downloader"
	agentv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/agent/v2"
	mapv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/map/v2"
	traffic_lightv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/traffic_light/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 下载地图，返回pb格式的地图数据与通过external得到的额外信息
func downloadMap(col *mongo.Collection) *mapv2.Map {
	pb := &mapv2.Map{}
	// header
	pb.Header = &mapv2.Header{}
	raw, err := col.FindOne(context.Background(), bson.M{"class": "header"}).DecodeBytes()
	if err != nil {
		log.Fatalln(err)
	}
	if err := downloader.UnmarshalBson(raw, pb.Header); err != nil {
		log.Fatalln(err)
	}
	// lanes
	pb.Lanes = downloader.DownloadPbsFromMongo(
		col, "lane",
		func(pb *mapv2.Lane, rawBson bson.Raw) error { return nil },
	)
	// roads
	pb.Roads = downloader.DownloadPbsFromMongo(
		col, "road",
		func(pb *mapv2.Road, rawBson bson.Raw) error {
			return nil
		},
	)
	// junctions
	pb.Junctions = downloader.DownloadPbsFromMongo(
		col, "junction",
		func(pb *mapv2.Junction, rawBson bson.Raw) error { return nil },
	)
	// aois
	pb.Aois = downloader.DownloadPbsFromMongo(
		col, "aoi",
		func(pb *mapv2.Aoi, rawBson bson.Raw) error {
			return nil
		},
	)
	return pb
}

// 下载agents，并通过检查所获取到的agent是否合法（不合法将被过滤）
func downloadAgents(col *mongo.Collection) *agentv2.Agents {
	pb := &agentv2.Agents{}
	pb.Agents = downloader.DownloadPbsFromMongo(
		col, "agent",
		func(agent *agentv2.Agent, rawBson bson.Raw) error {
			return nil
		},
	)
	return pb
}

// 下载traffic lights，并检查所获取到的是否合法（不合法将被过滤）
func downloadTrafficLights(col *mongo.Collection) *traffic_lightv2.TrafficLights {
	pb := &traffic_lightv2.TrafficLights{}
	pb.TrafficLights = downloader.DownloadPbsFromMongo(
		col, "traffic_light",
		func(tl *traffic_lightv2.TrafficLight, rawBson bson.Raw) error {
			return nil
		},
	)
	return pb
}

package main

import (
	"sync"

	"git.fiblab.net/sim/simulet-go/entity"
	"git.fiblab.net/sim/simulet-go/entity/aoi"
	"git.fiblab.net/sim/simulet-go/entity/junction"
	"git.fiblab.net/sim/simulet-go/entity/lane"
	"git.fiblab.net/sim/simulet-go/entity/person"
	"git.fiblab.net/sim/simulet-go/entity/person/route"
	"git.fiblab.net/sim/simulet-go/rpc/client"
	"git.fiblab.net/sim/simulet-go/utils/clock"
	"git.fiblab.net/sim/simulet-go/utils/config"
	"git.fiblab.net/sim/simulet-go/utils/input"

	traffic_lightv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/traffic_light/v2"
	"github.com/sirupsen/logrus"
)

const (
	STEP_LOG_INTERVAL = 100
)

var (
	// log
	LOG_LEVEL = map[string]logrus.Level{
		"trace":    logrus.TraceLevel,
		"debug":    logrus.DebugLevel,
		"info":     logrus.InfoLevel,
		"warn":     logrus.WarnLevel,
		"error":    logrus.ErrorLevel,
		"critical": logrus.FatalLevel,
		"off":      logrus.PanicLevel,
	}
)

func Init(c config.Config) {

	// input
	initRes := input.Init(c, *cacheDir)

	clock.Init(c)
	config.Init(c)

	// 数据加载
	mapData := initRes.Map
	persons := initRes.Persons.Agents
	tls := initRes.TrafficLights.TrafficLights

	log.Infof("Lane: %v", len(mapData.Lanes))
	log.Infof("Road: %v", len(mapData.Roads))
	log.Infof("Junction: %v", len(mapData.Junctions))
	log.Infof("AOI: %v", len(mapData.Aois))
	log.Infof("Person: %v", len(persons))

	// entity manager
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		person.Manager.Init(persons)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		aoi.Manager.Init(mapData.Aois)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		junction.Manager.Init(mapData.Junctions)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		lane.Manager.Init(mapData.Lanes)
	}()
	wg.Wait()

	// 初始化数据集合
	wg.Add(1)
	go func() {
		defer wg.Done()
		junction.Manager.InitLanes(lane.Manager)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		aoi.Manager.InitLanes(lane.Manager)
	}()
	for _, p := range person.Manager.Data() {
		home := p.Home()
		wg.Add(1)
		go func(p *person.Person) {
			defer wg.Done()
			// 假设所有的人一开始都在AOI里
			if home.AoiPosition == nil {
				log.Panicf("person %d has no home aoi position", p.ID())
			}
			aoiID := home.AoiPosition.AoiId
			if aoi, err := aoi.Manager.Get(aoiID); err != nil {
				log.Panic(err)
			} else {
				aoi.Add(
					p,
					entity.AoiMoveType_INIT,
					entity.AoiMoveType_SLEEP,
					-1,
				)
			}
		}(p)
	}
	wg.Wait()
	for _, tl := range tls {
		wg.Add(1)
		go func(tl *traffic_lightv2.TrafficLight) {
			defer wg.Done()
			junction, err := junction.Manager.Get(tl.GetJunctionId())
			if err != nil {
				log.Panic(err)
			}
			if err := junction.SetTrafficLight(tl); err != nil {
				log.Panic(err)
			}
		}(tl)
	}
	route.RoutingClient = client.NewRoutingServiceClient()

	wg.Wait()

	// log: 运行时才修改
	logrus.SetLevel(LOG_LEVEL[c.Log.Level])
}

func Step() {
	clock.GlobalTime = float64(clock.Step) * clock.STEP_INTERVAL

	// Prepare

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		person.Manager.Prepare() // person
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		aoi.Manager.Prepare() // aoi
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		junction.Manager.Prepare() // junction
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		lane.Manager.Prepare() // lane
	}()
	wg.Wait()

	// Update

	wg.Add(1)
	go func() {
		defer wg.Done()
		person.Manager.Update(clock.STEP_INTERVAL) // person
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		aoi.Manager.Update(clock.STEP_INTERVAL) // aoi
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		junction.Manager.Update(clock.STEP_INTERVAL) // junction
	}()

	wg.Wait()

	// 等待导航请求处理完成
	route.RoutingClient.Wait()
}

func Run() {
	for clock.Step < clock.END_STEP {
		Step()
		clock.Step++
	}
	route.RoutingClient.Close()
	log.Info("engine complete")
}

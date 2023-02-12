package config

type MongoPath struct {
	DB  string `yaml:"db"`
	Col string `yaml:"col"`
}

func (m MongoPath) GetDb() string {
	return m.DB
}

func (m MongoPath) GetCol() string {
	return m.Col
}

type Log struct {
	Level string `yaml:"level"`
}

type Mongo struct {
	URI          string     `yaml:"uri"`
	Map          MongoPath  `yaml:"map"`
	Person       MongoPath  `yaml:"person"`
	TrafficLight *MongoPath `yaml:"traffic_light,omitempty"`
}

type ControlStep struct {
	Start    int32   `yaml:"start"`
	Total    int32   `yaml:"total"`
	Interval float64 `yaml:"interval"`
}
type Control struct {
	Step ControlStep `yaml:"step"`
}

type Config struct {
	Log     Log     `yaml:"log"`
	Mongo   Mongo   `yaml:"mongo"`
	Control Control `yaml:"control"`
}

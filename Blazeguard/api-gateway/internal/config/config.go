package config

import "os"
type Config struct{
	Port string
	KafkaBroker string
}
func Load() Config{
	port:=os.Getenv("PORT")
	if port==""{
		port="8080"
	}
	broker:=os.Getenv("KAFKA_BROKER")
	if broker==""{
		broker="localhost:9092"
	}
	return Config{
		Port: port,
		KafkaBroker: broker,
	}
}

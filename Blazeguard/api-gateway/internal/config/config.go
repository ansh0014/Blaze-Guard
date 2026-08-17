package config

import "os"
type Config struct{
	Port string
	KafkaBroker string
	MLModelURL string
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
	mlURL:=os.Getenv("ML_MODEL_URL")
	if mlURL==""{
		mlURL="http://localhost:9000"
	}
	return Config{
		Port: port,
		KafkaBroker: broker,
		MLModelURL: mlURL,
	}
}

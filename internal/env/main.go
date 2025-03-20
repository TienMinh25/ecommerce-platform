package env

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type RedisConfig struct {
	// redis
	RedisAddress  string `envconfig:"REDIS_ADDRESS"`
	RedisPassword string `envconfig:"REDIS_PASSWORD"`
}

type MinioConfig struct {
	// minio
	MinioEndpointURL         string `envconfig:"MINIO_ENDPOINT_URL"`
	MinioAccessKey           string `envconfig:"MINIO_ACCESS_KEY"`
	MinioSecretKey           string `envconfig:"MINIO_SECRET_KEY"`
	MinioBucketAvatars       string `envconfig:"MINIO_BUCKET_AVATARS"`
	MinioRegion              string `envconfig:"MINIO_REGION"`
	MinioBucketAvatarsPolicy string `envconfig:"MINIO_BUCKET_AVATARS_POLICY"`
}

type JeagerConfig struct {
	// jeager (used for distributed tracing)
	JeagerExporterHttpEndpoint string `envconfig:"JAEGER_EXPORTER_HTTP_ENDPOINT"`
	JeagerExporterGrpcEndpoint string `envconfig:"JAEGER_EXPORTER_GRPC_ENDPOINT"`
}

type KafkaConfig struct {
	// kafka
	KafkaBrokers               string `envconfig:"KAFKA_BROKERS" default:"localhost:29092,localhost:29093,localhost:29094"`
	KafkaGroupID               string `envconfig:"KAFKA_GROUP_ID" default:"my-consumer-group"`
	KafkaRequestTimeout        int    `envconfig:"KAFKA_REQUEST_TIMEOUT" default:"3000"`
	KafkaRetryAttempts         int    `envconfig:"KAFKA_RETRY_ATTEMPTS" default:"3"`
	KafkaRetyDelay             int    `envconfig:"KAFKA_RETRY_DELAY" default:"2000"`
	KafkaConsumerFetchMinBytes int    `envconfig:"KAFKA_CONSUMER_FETCH_MIN_BYTES" default:"5"`
	KafkaConsumerFetchMaxBytes int    `envconfig:"KAFKA_CONSUMER_FETCH_MAX_BYTES" default:"1000000"`
	KafkaConsumerMaxWait       int    `envconfig:"KAFKA_CONSUMER_MAX_WAIT" default:"300"`
	KafkaProducerMaxWait       int    `envconfig:"KAFKA_PRODUCER_MAX_WAIT" default:"300"`
}

type MailConfig struct {
	MailHogFrom     string `envconfig:"MAILHOG_FROM"`
	MailHogSmtpHost string `envconfig:"MAILHOG_SMTP_SERVER"`
	MailHogPort     string `envconfig:"MAILHOG_PORT"`
	MailHogUserName string `envconfig:"MAILHOG_USERNAME"`
	MailHogPassword string `envconfig:"MAILHOG_PASSWORD"`
}

type MongoDBConfig struct {
	MongoDbUri      string `envconfig:"MONGODB_URI"`
	MongoDbDatabase string `envconfig:"MONGODB_DATABASE"`
	MongoDbUsername string `envconfig:"MONGODB_USERNAME"`
	MongoDbPassword string `envconfig:"MONGODB_PASSWORD"`
}

type ServiceWorkerPoolConfig struct {
	CapacityWorkerPool int `envconfig:"NUM_WORKER_POOL_DEFAULT" default:"5"`
	MessageSize        int `envconfig:"NUM_MESSAGE_IN_QUEUE_DEFAULT" default:"1000"`
}

type EnvManager struct {
	Redis             *RedisConfig
	Jeager            *JeagerConfig
	Mail              *MailConfig
	Mongo             *MongoDBConfig
	Kafka             *KafkaConfig
	ServiceWorkerPool *ServiceWorkerPoolConfig

	ApiSecretKey string `envconfig:"API_SECRET_KEY"`

	OneSignalAppId      string `envconfig:"ONESIGNAL_APP_ID"`
	OneSignalRestApiKey string `envconfig:"ONESIGNAL_REST_API_KEY"`

	JetStreamName          string `envconfig:"JETSTREAM_STREAM_NAME"`
	JetStreamDurable       string `envconfig:"JETSTREAM_DURABLE"`
	JetStreamOrdersTopic   string `envconfig:"JETSTREAM_ORDERS_TOPIC"`
	JetStreamProductsTopic string `envconfig:"JETSTREAM_PRODUCTS_TOPIC"`
	JetStreamPartnersTopic string `envconfig:"JETSTREAM_PARTNERS_TOPIC"`
}

func NewEnvManager() *EnvManager {
	var config EnvManager

	if err := envconfig.Process("", config); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	return &config
}

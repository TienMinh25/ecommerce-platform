package env

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"path/filepath"
	"runtime"
)

type PostgreSQLConfig struct {
	PostgresHost     string `envconfig:"POSTGRES_HOST"`
	PostgresPort     int    `envconfig:"POSTGRES_PORT"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	PostgresDSN      string `envconfig:"POSTGRES_DSN"`
}

type RedisConfig struct {
	// redis
	RedisAddress  string `envconfig:"REDIS_ADDRESS"`
	RedisPassword string `envconfig:"REDIS_PASSWORD"`
	RedisDB       int    `envconfig:"REDIS_DB"`
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

type ServerConfig struct {
	ServerAddresss string `envconfig:"SERVER_ADDRESS"`
	ConsumeGroup   string `envconfig:"API_GATEWAY_CONSUME_GROUP"`
}

type NotificationServerConfig struct {
	ServerAddresss string `envconfig:"NOTIFICATION_ADDRESS"`
	ConsumeGroup   string `envconfig:"NOTIFICATION_CONSUME_GROUP"`
}

type EnvManager struct {
	ServerConfig             *ServerConfig
	PostgreSQL               *PostgreSQLConfig
	Redis                    *RedisConfig
	Jeager                   *JeagerConfig
	Mail                     *MailConfig
	Mongo                    *MongoDBConfig
	MinioConfig              *MinioConfig
	Kafka                    *KafkaConfig
	ServiceWorkerPool        *ServiceWorkerPoolConfig
	NotificationServerConfig *NotificationServerConfig

	OTPVerifyEmailTimeout int `envconfig:"OTP_VERIFY_EMAIL_TIMEOUT"`

	PrivateKeyPath string `envconfig:"PRIVATE_KEY_PATH"`
	PublicKeyPath  string `envconfig:"PUBLIC_KEY_PATH"`

	ExpireAccessToken  int `envconfig:"EXPIRE_ACCESS_TOKEN"`
	ExpireRefreshToken int `envconfig:"EXPIRE_REFRESH_TOKEN"`

	TopicVerifyOTP string `envconfig:"TOPIC_VERIFY_OTP"`
}

func NewEnvManager() *EnvManager {
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		return nil
	}

	configPath := filepath.Join(filepath.Dir(filename), "../../configs/.env.prod")
	err := godotenv.Load(configPath)

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	log.Println("✅ Loaded env file successfully")

	var config EnvManager

	if err = envconfig.Process("", &config); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	log.Println("✅ Loaded environment variables successfully")

	return &config
}

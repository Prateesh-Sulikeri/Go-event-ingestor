package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST string
	DB_PORT string
	DB_USER string
	DB_PASSWORD string
	DB_NAME string

	REDIS_ADDR string

	JWT_SECRET string
	JWT_ISSUER string
	JWT_EXP_HOURS int

	RATE_LIMIT_BUCKET_SIZE int
	RATE_LIMIT_REFILL_RATE int
}

func Load() Config {
	_ = godotenv.Load()

	exp, err := strconv.Atoi(os.Getenv("JWT_EXP_HOURS"))
	if err != nil {
		exp = 1
	}

	bucket, err := strconv.Atoi(os.Getenv("RATE_LIMIT_BUCKET_SIZE"))
	if err != nil {
		log.Fatal("invalid RATE_LIMIT_BUCKET_SIZE")
	}

	refill, err := strconv.Atoi(os.Getenv("RATE_LIMIT_REFILL_RATE"))
	if err != nil {
		log.Fatal("invalid RATE_LIMIT_REFILL_RATE")
	}

	return Config{
		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: os.Getenv("DB_PORT"),
		DB_USER: os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME: os.Getenv("DB_NAME"),

		REDIS_ADDR: os.Getenv("REDIS_ADDR"),

		JWT_SECRET: os.Getenv("JWT_SECRET"),
		JWT_ISSUER: os.Getenv("JWT_ISSUER"),
		JWT_EXP_HOURS: exp,

		RATE_LIMIT_BUCKET_SIZE: bucket,
		RATE_LIMIT_REFILL_RATE: refill,
	}
}

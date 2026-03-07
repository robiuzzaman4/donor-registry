package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Version       string
	ServiceName   string
	Port          int
	DbUrl         string
	JwtSecret     string
	AdminName     string
	AdminPhone    string
	AdminPassword string
}

var cnf *Config

func loadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load env variables")
		os.Exit(1)
	}

	version := os.Getenv("VERSION")
	if version == "" {
		fmt.Println("Missing env variable 'VERSION'")
		os.Exit(1)
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		fmt.Println("Missing env variable 'SERVICE_NAME'")
		os.Exit(1)
	}
	portStr := os.Getenv("PORT")
	if portStr == "" {
		fmt.Println("Missing env variable 'PORT'")
		os.Exit(1)
	}

	port, err := strconv.ParseInt(portStr, 10, 64)
	if portStr == "" {
		fmt.Println("PORT must be a number")
		os.Exit(1)
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		fmt.Println("Missing env variable 'DB_URL'")
		os.Exit(1)
	}

	jwtSceret := os.Getenv("JWT_SECRET")
	if jwtSceret == "" {
		fmt.Println("Missing env variable 'JWT_SECRET'")
		os.Exit(1)
	}

	adminName := os.Getenv("ADMIN_NAME")
	if adminName == "" {
		fmt.Println("Missing env variable 'ADMIN_NAME'")
		os.Exit(1)
	}

	adminPhone := os.Getenv("ADMIN_PHONE")
	if adminPhone == "" {
		fmt.Println("Missing env variable 'ADMIN_PHONE'")
		os.Exit(1)
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		fmt.Println("Missing env variable 'ADMIN_PASSWORD'")
		os.Exit(1)
	}

	cnf = &Config{
		Version:       version,
		ServiceName:   serviceName,
		Port:          int(port),
		DbUrl:         dbUrl,
		JwtSecret:     jwtSceret,
		AdminName:     adminName,
		AdminPhone:    adminPhone,
		AdminPassword: adminPassword,
	}
	return cnf
}

func GetConfig() *Config {
	if cnf == nil {
		loadConfig()
	}
	return cnf
}

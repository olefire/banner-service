package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Port       string
	PgDSN      string
	PublicKey  string
	PrivateKey string
	TokenTTL   time.Duration
	BannerTTL  time.Duration
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(fmt.Errorf("error config file: %w", err))
		viper.SetDefault("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432")
		viper.SetDefault("PORT", ":8000")
		viper.SetDefault("TOKEN_TTL", time.Hour)
		viper.SetDefault("BANNER_TTL", 5*time.Minute)
		viper.SetDefault("ACCESS_TOKEN_PRIVATE_KEY", "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlCUEFJQkFBSkJBTzVIKytVM0xrWC91SlRvRHhWN01CUURXSTdGU0l0VXNjbGFFKzlaUUg5Q2VpOGIxcUVmCnJxR0hSVDVWUis4c3UxVWtCUVpZTER3MnN3RTVWbjg5c0ZVQ0F3RUFBUUpCQUw4ZjRBMUlDSWEvQ2ZmdWR3TGMKNzRCdCtwOXg0TEZaZXMwdHdtV3Vha3hub3NaV0w4eVpSTUJpRmI4a25VL0hwb3piTnNxMmN1ZU9wKzVWdGRXNApiTlVDSVFENm9JdWxqcHdrZTFGY1VPaldnaXRQSjNnbFBma3NHVFBhdFYwYnJJVVI5d0loQVBOanJ1enB4ckhsCkUxRmJxeGtUNFZ5bWhCOU1HazU0Wk1jWnVjSmZOcjBUQWlFQWhML3UxOVZPdlVBWVd6Wjc3Y3JxMTdWSFBTcXoKUlhsZjd2TnJpdEg1ZGdjQ0lRRHR5QmFPdUxuNDlIOFIvZ2ZEZ1V1cjg3YWl5UHZ1YStxeEpXMzQrb0tFNXdJZwpQbG1KYXZsbW9jUG4rTkVRdGhLcTZuZFVYRGpXTTlTbktQQTVlUDZSUEs0PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ==")
		viper.SetDefault("ACCESS_TOKEN_PUBLIC_KEY", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTzVIKytVM0xrWC91SlRvRHhWN01CUURXSTdGU0l0VQpzY2xhRSs5WlFIOUNlaThiMXFFZnJxR0hSVDVWUis4c3UxVWtCUVpZTER3MnN3RTVWbjg5c0ZVQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==")
	}

	return &Config{
		PgDSN:      viper.GetString("DATABASE_DSN"),
		Port:       viper.GetString("PORT"),
		PrivateKey: viper.GetString("ACCESS_TOKEN_PRIVATE_KEY"),
		PublicKey:  viper.GetString("ACCESS_TOKEN_PUBLIC_KEY"),
		TokenTTL:   viper.GetDuration("TOKEN_TTL"),
		BannerTTL:  viper.GetDuration("BANNER_TTL"),
	}
}

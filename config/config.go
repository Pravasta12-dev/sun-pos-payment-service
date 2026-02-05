package config

import "github.com/spf13/viper"

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`
}

type PsqlDB struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	DBName    string `json:"db_name"`
	DBMaxOpen int    `json:"db_max_open"`
	DBMaxIdle int    `json:"db_max_idle"`
}

type Midtrans struct {
	Env     string `json:"env"`
	BaseURL string `json:"base_url"`
}

type Payment struct {
	QrisAcquirer     string `json:"qris_acquirer"`
	QrisExpireMinute int    `json:"qris_expire_minute"`
}

type Security struct {
	EncryptionSecret string `json:"encryption_secret"`
}

type Websocket struct {
	Enabled bool `json:"enabled"`
}

type Config struct {
	App       App       `json:"app"`
	PsqlDB    PsqlDB    `json:"psql_db"`
	Midtrans  Midtrans  `json:"midtrans"`
	Payment   Payment   `json:"payment"`
	Security  Security  `json:"security"`
	Websocket Websocket `json:"websocket"`
}

func NewConfig() *Config {
	var baseUrl string
	baseProductionUrl := viper.GetString("MIDTRANS_BASE_URL")
	baseSandboxUrl := viper.GetString("MIDTRANS_BASE_SANDBOX_URL")
	midtransEnv := viper.GetString("MIDTRANS_ENV")

	if midtransEnv == "production" {
		baseUrl = baseProductionUrl
	} else {
		baseUrl = baseSandboxUrl
	}

	return &Config{
		App: App{
			AppPort: viper.GetString("APP_PORT"),
			AppEnv:  viper.GetString("APP_ENV"),
		},
		PsqlDB: PsqlDB{
			Host:      viper.GetString("DATABASE_HOST"),
			Port:      viper.GetInt("DATABASE_PORT"),
			User:      viper.GetString("DATABASE_USER"),
			Password:  viper.GetString("DATABASE_PASSWORD"),
			DBName:    viper.GetString("DATABASE_NAME"),
			DBMaxOpen: viper.GetInt("DATABASE_MAX_OPEN_CONNECTIONS"),
			DBMaxIdle: viper.GetInt("DATABASE_MAX_IDLE_CONNECTIONS"),
		},
		Midtrans: Midtrans{
			Env:     midtransEnv,
			BaseURL: baseUrl,
		},
		Payment: Payment{
			QrisAcquirer:     viper.GetString("PAYMENT_QRIS_ACQUIRER"),
			QrisExpireMinute: viper.GetInt("PAYMENT_QRIS_EXPIRE_MINUTE"),
		},
		Security: Security{
			EncryptionSecret: viper.GetString("ENCRYPTION_SECRET_KEY"),
		},
		Websocket: Websocket{
			Enabled: viper.GetBool("WS_ENABLED"),
		},
	}
}

package config

type Session interface {
	LoadVariables() (map[string]string, error)
}

func GetSession(c Config, rc rawConfig) {

}

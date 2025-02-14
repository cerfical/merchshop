package postgres

type Config struct {
	Host string
	Port string
	Name string

	User     string
	Password string

	Migrations string
}

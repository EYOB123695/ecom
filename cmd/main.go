package main  


func main() {
	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: "",
		},
	}

	_ = cfg
}

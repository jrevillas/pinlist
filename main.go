package main

func main() {
	config := Load()

	DB := &Connection{}

	// Init database
	DB.initDatabase(config.ConnectionString)

	Start(DB, config)
}

package initialize

func Run() error {
	LoadConfig()

	SetupLogger()
	SetupMongoDB()
	SetupRedis()
	server := InitializeServer()
	
	// LoadData()

	return server.Run()
}

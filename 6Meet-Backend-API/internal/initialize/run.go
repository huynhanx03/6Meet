package initialize

func Run() error {
	LoadConfig()

	InitLogger()
	InitMongoDB()
	InitRedis()
	InitDependencies()
	InitRouter()

	// LoadData()

	return nil
}

package initialize

func Run() error {
	LoadConfig()
	
	InitLogger()
	// InitMongoDB()
	// InitDependencies()
	// InitRouter()

	LoadData()
	
	return nil
}

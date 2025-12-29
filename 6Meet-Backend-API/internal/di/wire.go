package di

// SetupDependencies initializes all dependencies and assigns them to GlobalContainer
func SetupDependencies() {
	// Initialize user dependencies
	userRepo, userCache, userSearch, userService, userHandler, crawlService, userProducer, userCrawlConsumer, userCdcConsumer := InitUserDependencies()

	GlobalContainer = &Container{
		UserRepo:            userRepo,
		UserCache:           userCache,
		UserSearch:          userSearch,
		UserService:         userService,
		UserHandler:         userHandler,
		CrawlService:        crawlService,
		UserCrawlerConsumer: userCrawlConsumer,
		UserCDCConsumer:     userCdcConsumer,
		UserProducer:        userProducer,
	}
}

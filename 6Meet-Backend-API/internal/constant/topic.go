package constant

const (
	// TopicCrawlWikiPages is the topic for crawl wiki pages
	TopicCrawlWikiPages = "crawl-wiki-pages"

	// TopicUserDebezium is the topic name for Debezium user changes
	TopicUserDebezium = "mongo.users.sixmeet.users"

	// ConsumerGroupUserSaver is the group ID for user saver consumer
	ConsumerGroupUserSaver = "user_saver_group"

	// ConsumerGroupUserCDC is the group ID for user CDC consumer
	ConsumerGroupUserCDC = "user_cdc_group"
)

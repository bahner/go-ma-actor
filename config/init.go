package config

// This should be done after flag.Parse() in main.
func Init() {
	InitLogging()
	InitNodeIdentity()
	InitIdentity()

}

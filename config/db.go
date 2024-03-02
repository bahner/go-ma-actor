package config

const (
	defaultDBFilename = "ma.db"
)

var (
	defaultDbFile = dataHome + "/" + defaultDBFilename
)

func DefaultDBFile() string {
	return defaultDbFile
}

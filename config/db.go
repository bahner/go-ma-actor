package config

const defaultDBFilename = "ma.db"

var DefaultDbFile = NormalisePath(dataHome + defaultDBFilename)

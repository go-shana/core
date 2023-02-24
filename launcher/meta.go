package launcher

// metaData stores build information.
var metaData = map[string]string{}

// SetMeta sets the meta data.
func SetMeta(key, value string) {
	metaData[key] = value
}

func init() {
	SetMeta("version", "tip")
}

package base

type Config struct {
	ScanEventCountLimit int
	MaxAreasByUser      int
}

func (b *Base) loadConfig() {
	config := Config{
		ScanEventCountLimit: 100,
		MaxAreasByUser:      10,
	}

	b.Config = &config
}

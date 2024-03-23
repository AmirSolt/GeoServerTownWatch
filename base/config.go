package base

type Config struct {
	ScanEventCountLimit int
}

func (b *Base) loadConfig() {
	config := Config{
		ScanEventCountLimit: 100,
	}

	b.Config = &config
}

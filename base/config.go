package base

type Config struct {
	ScanEventCountLimit int
}

func (b *Base) loadConfig() {
	config := Config{
		ScanEventCountLimit: 50,
	}

	b.Config = &config
}

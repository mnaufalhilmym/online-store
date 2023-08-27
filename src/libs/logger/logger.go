package logger

type logger struct {
	prefix string
}

func New(prefix string) logger {
	return logger{
		prefix: prefix,
	}
}

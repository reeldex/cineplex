package logger

import "go.uber.org/zap"

func MustNew(srvName string) (*zap.Logger, func() error) {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return l.With(zap.String("service_name", srvName)), l.Sync
}

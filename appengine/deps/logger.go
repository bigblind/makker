package deps

import (
	"google.golang.org/appengine/log"
	"context"
	"github.com/bigblind/makker/di"
	"github.com/bigblind/makker/logging"
)

func loggerFactory(ctx context.Context)  {
}


type gaeLogger struct{}

func (gl gaeLogger) Debugf(ctx context.Context, format string, args ... interface{}) {
	log.Debugf(ctx, format, args...)
}

func (gl gaeLogger) Infof(ctx context.Context, format string, args ... interface{}) {
	log.Infof(ctx, format, args...)
}

func (gl gaeLogger) Warnf(ctx context.Context, format string, args ... interface{}) {
	log.Warningf(ctx, format, args...)
}

func (gl gaeLogger) Error(ctx context.Context, format string, args ... interface{}) {
	log.Errorf(ctx, format, args...)
}

func init()  {
	di.Graph.Provide(func() *logging.StructuredLogger {
		return logging.NewStructuredLogger(gaeLogger{})
	})
}





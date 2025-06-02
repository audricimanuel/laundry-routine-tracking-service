package logging

import (
	"context"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/constants"
	"github.com/sirupsen/logrus"
)

func WithContext(ctx context.Context) *logrus.Entry {
	log := logrus.WithContext(ctx)
	if ctx.Value(constants.TRACE_ID) != nil {
		log = log.WithFields(logrus.Fields{
			constants.TRACE_ID: ctx.Value(constants.TRACE_ID).(string),
		})
	}
	return log
}

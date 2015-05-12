package model

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/tkusd/server/util"
)

var sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)

type dbLogger struct{}

func (l *dbLogger) Print(v ...interface{}) {
	if len(v) == 0 {
		return
	}

	level := v[0]
	source := v[1]
	now := time.Now()

	if level == "sql" {
		duration := v[2].(time.Duration)

		var formatedValues []interface{}

		for _, value := range v[4].([]interface{}) {
			indirectValue := reflect.Indirect(reflect.ValueOf(value))

			if indirectValue.IsValid() {
				value = indirectValue.Interface()

				if t, ok := value.(time.Time); ok {
					formatedValues = append(formatedValues, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
				} else if b, ok := value.([]byte); ok {
					formatedValues = append(formatedValues, fmt.Sprintf("'%v'", string(b)))
				} else if r, ok := value.(driver.Valuer); ok {
					if value, err := r.Value(); err == nil && value != nil {
						formatedValues = append(formatedValues, fmt.Sprintf("'%v'", value))
					} else {
						formatedValues = append(formatedValues, "NULL")
					}
				} else {
					formatedValues = append(formatedValues, fmt.Sprintf("'%v'", value))
				}
			} else {
				formatedValues = append(formatedValues, fmt.Sprintf("'%v'", value))
			}
		}

		msg := fmt.Sprintf(sqlRegexp.ReplaceAllString(v[3].(string), "%v"), formatedValues...)

		util.Log().WithFields(logrus.Fields{
			"level":    level,
			"source":   source,
			"duration": duration,
			"time":     now,
		}).Info(msg)
	} else {
		util.Log().WithFields(logrus.Fields{
			"level":  level,
			"source": source,
			"time":   now,
		}).Info(v[2:]...)
	}
}

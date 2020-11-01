package importer

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Stats struct {
	duplicates   int
	success      int
	tooMany      int
	nonIndexable int
	fail         int
}

func (stats Stats) String() string {
	return fmt.Sprintf("success:%d duplicates:%d tooMany:%d nonIndexable:%d  fail:%d",
		stats.success,
		stats.duplicates,
		stats.tooMany,
		stats.nonIndexable,
		stats.fail)
}

func (stats Stats) print() {
	log.WithFields(log.Fields{
		"success":      stats.success,
		"fail":         stats.fail,
		"nonIndexable": stats.nonIndexable,
		"tooMany":      stats.tooMany,
		"duplicates":   stats.duplicates,
	}).Info("print stats")
}

package output

import (
	"github.com/yalhyane/another-redis-memory-analyzer/utils"
	"io"
)

type ReportOutput interface {
	Output(r utils.DBReports)
}

type ReportWriter io.Writer

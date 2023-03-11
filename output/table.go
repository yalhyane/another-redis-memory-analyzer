package output

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/yalhyane/another-redis-memory-analyzer/utils"
	"log"
	"os"
	"strconv"
)

type TableOutput struct {
	MinSize uint64
	Out     ReportWriter
}

func NewTableOutput(minSize uint64, out ReportWriter) *TableOutput {
	return &TableOutput{minSize, out}
}

func DefaultTableOutput() *TableOutput {
	return &TableOutput{MinSize: 0, Out: os.Stdout}
}

func (o *TableOutput) Output(r utils.DBReports) {
	var (
		size string
	)
	var redisTotalSize uint64
	var redisTotalKeys uint64
	dbsTable := tablewriter.NewWriter(o.Out)
	dbsTable.SetHeader([]string{"DB", "Count", "Size"})
	for db, reports := range r {
		table := tablewriter.NewWriter(o.Out)
		table.SetHeader([]string{"Key", "Count", "Size"})
		var dbTotalSize uint64
		var dbTotalCount uint64
		for _, value := range reports {
			dbTotalSize += value.Size
			dbTotalCount += value.Count
			// if less than 1Mb ignore...
			if value.Size < uint64(o.MinSize) {
				continue
			}
			size = humanize.Bytes(value.Size)
			table.Append([]string{value.Key, strconv.Itoa(int(value.Count)), size})
		}
		redisTotalSize += dbTotalSize
		redisTotalKeys += dbTotalCount

		// db total size...
		size := humanize.Bytes(dbTotalSize)
		table.Append([]string{"Total", strconv.Itoa(int(dbTotalCount)), size})

		log.Println("DB", db, "total size is:", size)
		table.Render()
		dbsTable.Append([]string{fmt.Sprintf("DB %d", db), strconv.Itoa(int(dbTotalCount)), size})
	}

	size = humanize.Bytes(redisTotalSize)
	log.Println("Redis total size is:", size)
	dbsTable.Append([]string{"All", strconv.Itoa(int(redisTotalKeys)), size})
	dbsTable.Render()
}

package output

import (
	"encoding/json"
	"github.com/dustin/go-humanize"
	"github.com/yalhyane/another-redis-memory-analyzer/utils"
	"log"
)

type JsonOutput struct {
	MinSize uint64
	Out     ReportWriter
}

type JsonKeyStructure struct {
	Key   string `json:"key"`
	Count uint64 `json:"count"`
	Size  string `json:"size"`
}

type JsonDBStructure struct {
	DB    uint64 `json:"db"`
	Count uint64 `json:"count"`
	Size  string `json:"size"`
}

type DBsDataStructure map[uint64][]JsonKeyStructure
type SummaryStructure []JsonDBStructure
type JsonStructure struct {
	DBs     DBsDataStructure `json:"DBs"`
	Summary SummaryStructure `json:"summary"`
	Total   JsonKeyStructure `json:"total"`
}

func (o *JsonOutput) Output(r utils.DBReports) {
	var (
		size string
	)
	var redisTotalSize uint64
	var redisTotalKeys uint64
	var dbsData = DBsDataStructure{}
	summaryData := make(SummaryStructure, 0, len(r))
	for db, reports := range r {
		var dbTotalSize uint64
		var dbTotalCount uint64
		dbsData[db] = make([]JsonKeyStructure, 0, len(reports))
		for _, value := range reports {
			dbTotalSize += value.Size
			dbTotalCount += value.Count
			// if less than 1Mb ignore...
			if value.Size < uint64(o.MinSize) {
				continue
			}
			size = humanize.Bytes(value.Size)
			dbsData[db] = append(dbsData[db], JsonKeyStructure{
				value.Key,
				value.Count,
				size,
			})
		}
		redisTotalSize += dbTotalSize
		redisTotalKeys += dbTotalCount

		// db total size...
		size := humanize.Bytes(dbTotalSize)
		summaryData = append(summaryData, JsonDBStructure{
			db,
			dbTotalCount,
			size,
		})

		log.Println("DB", db, "total size is:", size)
	}
	size = humanize.Bytes(redisTotalSize)
	log.Println("Redis total size is:", size)
	data := JsonStructure{
		DBs:     dbsData,
		Summary: summaryData,
		Total: JsonKeyStructure{
			"*",
			redisTotalKeys,
			size,
		},
	}
	d, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("Could not output json:", err)
	}

	_, err = o.Out.Write(d)
	if err != nil {
		log.Fatal("Could not write json:", err)
	}
	// add br at the end
	_, _ = o.Out.Write([]byte("\n"))

}

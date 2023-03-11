package utils

type Report struct {
	Key   string `json:"key"`
	Count uint64 `json:"count"`
	Size  uint64 `json:"size"`
}

type DBReports map[uint64][]Report
type KeyReports map[string]Report
type SizeReports []Report

func (sr SizeReports) Len() int {
	return len(sr)
}

func (sr SizeReports) Less(i, j int) bool {
	return sr[i].Size > sr[j].Size
}

func (sr SizeReports) Swap(i, j int) {
	sr[i], sr[j] = sr[j], sr[i]
}

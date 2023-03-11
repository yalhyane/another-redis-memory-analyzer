package utils

import (
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
	"sync"
)

func StartProgressBar(wg *sync.WaitGroup, width int) *mpb.Progress {
	return mpb.New(mpb.WithWaitGroup(wg), mpb.WithWidth(width))
}

func InitProgressBar(p *mpb.Progress, barTitle string, total int64) *mpb.Bar {
	bar := p.AddBar(total,
		mpb.PrependDecorators(
			// simple name decorator
			decor.Name(barTitle, decor.WC{W: len(barTitle) + 1, C: decor.DidentRight}),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
			// replace ETA decorator with "done" message, OnComplete event
			//decor.OnComplete(
			//	// ETA decorator with ewma age of 60
			//	decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
			//),
		),
	)

	return bar
}

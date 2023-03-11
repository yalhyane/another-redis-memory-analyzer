package cmd

import (
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/go-redis/redis"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v5"
	"github.com/yalhyane/another-redis-memory-analyzer/output"
	"github.com/yalhyane/another-redis-memory-analyzer/utils"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

var (
	host            string
	port            int
	password        string
	askPassword     bool
	delimiter       string
	delimiterLevel  int
	onlyDb          int
	scanLength      int
	humanMinSize    string
	minSize         uint64
	outFormat       []string
	goroutinesPerDb int
	disableBar      bool
)
var tableOutFormatStr = "table"
var jsonOutFormatStr = "json"
var supportedOutFormat = map[string]bool{
	tableOutFormatStr: true,
	jsonOutFormatStr:  true,
}
var supportedOutFormatStr string

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "analyze memory of redis instance",
	Long: `
arma - Another Redis Memory Analyzer - A Command Line Tool for Analyzing Redis Memory

The arma command line tool is designed to help Redis developers understand how DB and keys patterns are using memory in Redis instance.
This tool uses the Redis MEMORY USAGE command to retrieve the number of bytes that a key and its value are taking in RAM, and then group them by key pattern and by db.

Examples:
  To analyze memory usage of a Redis instance running on the default port on the local machine:
    arma analyze
  
  To analyze memory usage of a Redis instance running on a remote server with a custom port and password:
    arma analyze --host example.com --port 6380 --ask-password

  To analyze memory usage of a Redis instance and pattern keys by custom delimiter and a minimum level:
    arma analyze --delimiter "-" --level 2

  To analyze memory usage of a Redis instance on a specific database with larger scan range and only show big keys:
    arma analyze --db 3 --scan 10000 --min-size 1MB

Note: This tool is intended to provide insights into memory usage of group of keys in Redis, but should not be used as a substitute for proper memory management practices. It is important to regularly monitor your Redis instance and take appropriate actions to optimize memory usage, such as configuring eviction policies or increasing memory capacity.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		for _, f := range outFormat {
			f = strings.ToLower(f)
			if _, ok := supportedOutFormat[f]; !ok {
				return errors.New("invalid value for output format parameter, supported values: " + supportedOutFormatStr)
			}
		}
		var err error
		// ask for password
		if askPassword {
			prompt := promptui.Prompt{
				Label: "Password",
				Mask:  '*',
			}

			password, err = prompt.Run()
			if err != nil {
				return err
			}
		}
		minSize, err = humanize.ParseBytes(humanMinSize)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		analyze()
		return nil
	},
}

func init() {
	initOutputFormatStr()
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "Redis server host")
	analyzeCmd.Flags().IntVarP(&port, "port", "P", 6379, "Redis server port")
	analyzeCmd.Flags().StringVarP(&password, "password", "p", "", "Redis server password, it's recommended to use -a flag which will ask for your password which will prevent it from appearing in tty history")
	analyzeCmd.Flags().BoolVarP(&askPassword, "ask-password", "a", false, "Let tty asks for password (recommended)")
	analyzeCmd.Flags().BoolVarP(&disableBar, "no-progress-bar", "b", false, "No progress bar")
	analyzeCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ":", "The delimiter of keys")
	analyzeCmd.Flags().IntVarP(&delimiterLevel, "level", "l", 1, "The delimiter level of keys")
	analyzeCmd.Flags().IntVarP(&onlyDb, "db", "D", -1, "Redis DB to analyse ( -1 will analyze all dbs )")
	analyzeCmd.Flags().IntVarP(&scanLength, "scan", "S", 500, "Redis scan range length, number of keys to scan at once")
	analyzeCmd.Flags().StringVarP(&humanMinSize, "min-size", "s", "1KB", "Minimum size of group of keys to show, if group of keys is less than this size it won't be shown. Human readable size (KB, MB, GB...)")
	analyzeCmd.Flags().StringSliceVarP(&outFormat, "format", "o", []string{"table"}, "Output format, supported values: "+supportedOutFormatStr)
	analyzeCmd.Flags().IntVarP(&goroutinesPerDb, "goroutines", "g", 50, "Number of concurrent goroutines to analyze a redis DB")
}

func initOutputFormatStr() {
	keys := make([]string, len(supportedOutFormat))
	for k, _ := range supportedOutFormat {
		keys = append(keys, k)
	}
	supportedOutFormatStr = strings.Join(keys, ", ")
}

func analyze() {
	redisClient := utils.Connection(host, port, password, 0)
	defer utils.Close(redisClient)
	reports := Start(redisClient, delimiter, delimiterLevel, onlyDb, scanLength)

	writer := os.Stdout
	for _, f := range outFormat {
		var o output.ReportOutput
		switch f {
		case tableOutFormatStr:
			o = &output.TableOutput{
				MinSize: minSize,
				Out:     writer,
			}
		case jsonOutFormatStr:
			o = &output.JsonOutput{
				MinSize: minSize,
				Out:     writer,
			}
		}
		o.Output(reports)
	}
}

func Start(client *redis.Client, delimiter string, delimiterLevel int, onlyDb int, scanLength int) utils.DBReports {
	var dbReports = utils.DBReports{}
	// var reportsMapMutex = sync.RWMutex{}

	log.Println("Starting memory analysis")
	log.Println("Reading databases...")
	databases := utils.GetDatabases(client)

	sortedDbs := make([]int, 0, len(databases))
	for db := range databases {
		sortedDbs = append(sortedDbs, int(db))
	}
	sort.Ints(sortedDbs)

	var wg sync.WaitGroup
	var progressBar *mpb.Progress
	if !disableBar {
		progressBar = utils.StartProgressBar(&wg, 65)
	}

	for _, dbInt := range sortedDbs {

		if onlyDb > -1 && onlyDb != dbInt {
			// log.Printf("Ignore db %d: \n", dbInt)
			continue
		}

		//		log.Printf("Start db %d: \n", dbInt)

		wg.Add(1)

		go func(db uint64, keyCount int64) {
			defer wg.Done()
			var bar *mpb.Bar
			if !disableBar {
				barTitle := fmt.Sprintf("DB %d (%d keys)", db, keyCount)
				bar = utils.InitProgressBar(progressBar, barTitle, keyCount)
			}

			mr := utils.KeyReports{}

			// new connection to DB....
			rClient := utils.ConnectionToDb(int(db))

			var cursor uint64
			chunkSize := int64(scanLength)
			progress := 0
			var mrMapMutex = sync.RWMutex{}

			// allow 50 goroutines at once
			ch := make(chan int, goroutinesPerDb)

			var wwg sync.WaitGroup
			goroIndex := 0
			for {
				var keys []string
				var err error

				keys, cursor, err = rClient.Scan(cursor, "*", chunkSize).Result()
				if err != nil {
					panic(err)
				}

				// if cursor is not empty and no keys found
				// increase chunkSize
				if len(keys) == 0 && cursor != 0 && chunkSize < keyCount {
					chunkSize = chunkSize + chunkSize/2
				}

				keysLength := len(keys)
				progress += keysLength

				// block if ch reach 50...
				ch <- 1

				wwg.Add(1)
				goroIndex++
				go func(keys []string, ch chan int, index int, bar *mpb.Bar) {

					defer wwg.Done()

					groupKey := ""
					for _, key := range keys {
						if !disableBar {
							bar.IncrBy(1)
						}
						tmp := strings.Split(key, delimiter)
						if len(tmp) > 1 && delimiterLevel < len(tmp) {
							groupKey = strings.Join(tmp[0:delimiterLevel], delimiter) + delimiter + "*"
						} else {
							groupKey = key
						}

						// use MEMORY USAGE
						length, err := rClient.MemoryUsage(key, 0).Result()

						if err != nil {
							log.Fatal(err)
						}

						r := utils.Report{}
						mrMapMutex.RLock()
						if _, ok := mr[groupKey]; ok {
							r = mr[groupKey]
						} else {
							r = utils.Report{Key: groupKey}
						}
						mrMapMutex.RUnlock()
						r.Size += uint64(length)
						r.Count++
						mrMapMutex.Lock()
						mr[groupKey] = r
						mrMapMutex.Unlock()
					}

					// free space for another goroutine
					_ = <-ch

				}(keys, ch, goroIndex, bar)
				// stop scan
				if cursor == 0 {
					break
				}
			}

			wwg.Wait()

			sr := utils.SizeReports{}
			for _, report := range mr {
				sr = append(sr, report)
			}
			sort.Sort(sr)
			// reportsMapMutex.Lock()
			dbReports[db] = sr
			// reportsMapMutex.Unlock()

		}(uint64(dbInt), databases[uint64(dbInt)])
	}

	// wait for our bar to complete and flush
	if disableBar {
		wg.Wait()
	} else {
		progressBar.Wait()
	}

	return dbReports
}

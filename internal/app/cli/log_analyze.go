package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/loganalyzer"
)

func newLogAnalyzeCmd() *cobra.Command {
	var (
		top   int
		group string
	)
	cmd := &cobra.Command{
		Use:   "log-analyze",
		Short: "Analyze log files",
		Run: func(cmd *cobra.Command, args []string) {
			filePath := os.Getenv("ACCESS_LOG_FILE")
			if filePath == "" {
				filePath = "./logs/access.jsonl"
			}
			f, err := os.OpenFile(filePath, os.O_RDONLY, 0)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open log file: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()

			stats, kvs, err := loganalyzer.AnalyzeFile(filePath, loganalyzer.Options{
				GroupBy: group,
				TopN:    top,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "analyze error: %v\n", err)
				os.Exit(1)
			}
			printReport(stats, kvs)
		},
	}

	cmd.Flags().IntVarP(&top, "top", "t", 10, "Top N routes/paths")
	cmd.Flags().StringVarP(&group, "group", "g", "route", "Group by: route|path")

	return cmd
}

func printReport(s loganalyzer.Stats, top []loganalyzer.KV) {
	fmt.Printf("Requests: %d\n", s.Count)
	if !s.From.IsZero() {
		loc := time.Local
		if tz := os.Getenv("TZ"); tz != "" {
			if l, err := time.LoadLocation(tz); err == nil {
				loc = l
			}
		}
		fmt.Printf("Period : %s 〜 %s\n",
			s.From.In(loc).Format(time.RFC3339),
			s.To.In(loc).Format(time.RFC3339))
	}
	fmt.Printf("Status : ")
	for _, code := range []int{200, 201, 204, 400, 401, 403, 404, 500, 503} {
		if c, ok := s.ByStatus[code]; ok {
			fmt.Printf("%d:%d ", code, c)
		}
	}
	fmt.Printf("\nClass  : 2xx:%d 3xx:%d 4xx:%d 5xx:%d other:%d\n",
		s.ByClass["2xx"], s.ByClass["3xx"], s.ByClass["4xx"], s.ByClass["5xx"], s.ByClass["other"])

	fmt.Printf("Latency(ms) min:%d avg:%.1f p50:%d p95:%d p99:%d max:%d\n",
		s.LatMin, s.LatAvg, s.LatP50, s.LatP95, s.LatP99, s.LatMax)

	fmt.Printf("\nTop endpoints:\n")
	for i, kv := range top {
		fmt.Printf("%2d. %6d  %s\n", i+1, kv.Count, kv.Key)
	}
}

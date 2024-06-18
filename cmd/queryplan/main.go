package queryplan

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"pbench/log"
	"pbench/presto/plan_node"

	"github.com/spf13/cobra"
)

var (
	csvFile         string
	queryPlanColumn int
	hasHeader       bool
	output          string
)

// type lineHandler func()

func Run(_ *cobra.Command, args []string) {
	if csvFile == "" {
		log.Fatal().Str("CSV", csvFile).Msg("please specify the CSV file to read")
	}

	err := readCSVFile()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to process the CSV file")
	}
}

func readCSVFile() error {
	f, err := os.Open(csvFile)
	if err != nil {
		return err
	}

	defer f.Close()

	r := csv.NewReader(f)
	var rowNum = 1
	if hasHeader {
		if _, err := r.Read(); err != nil {
			log.Fatal().Err(err).Msg("failed to consume the header line")
		}
		rowNum++
	}

	for ; ; rowNum++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := parseQueryPlan(record[queryPlanColumn]); err != nil {
			log.Error().Err(err).Msgf("failed to parse the query plan at row:%d", rowNum)
		}
	}
	return nil
}

func parseQueryPlan(plan string) error {
	planTree := make(plan_node.PlanTree)
	if err := json.Unmarshal([]byte(plan), &planTree); err != nil {
		return err
	}

	joins, err := planTree.ParseJoins()

	if err != nil {
		return err
	}

	fmt.Printf("Joins: %v\n", joins)
	return nil
}

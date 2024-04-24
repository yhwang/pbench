package load

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"pbench/log"
	"pbench/presto/query_json"
	"pbench/stage"
	"pbench/utils"
	"reflect"
	"time"
)

var (
	Name          string
	Comment       string
	OutputPath    string
	MySQLCfgPath  string
	InfluxCfgPath string

	mysqlDb      *sql.DB
	runRecorders []stage.RunRecorder
	pseudoStage  *stage.Stage
	queryResults []*stage.QueryResult
)

func Run(_ *cobra.Command, args []string) {
	utils.PrepareOutputDirectory(OutputPath)
	runRecorders = make([]stage.RunRecorder, 0, 3)
	queryResults = make([]*stage.QueryResult, 0, 8)

	mysqlDb = utils.InitMySQLConnFromCfg(MySQLCfgPath)
	registerRunRecorder(stage.NewFileBasedRunRecorder())
	registerRunRecorder(stage.NewMySQLRunRecorderWithDb(mysqlDb))
	registerRunRecorder(stage.NewInfluxRunRecorder(InfluxCfgPath))

	pseudoStage = &stage.Stage{
		Id:       "load_json",
		ColdRuns: 1,
		States: &stage.SharedStageStates{
			RunName:      Name,
			Comment:      Comment,
			RunStartTime: time.Now(),
			OutputPath:   OutputPath,
		},
	}

	for _, path := range args {
		if err := processPath(path); err != nil {
			log.Error().Str("path", path).Err(err).Msg("failed to process path")
		}
	}

	for _, r := range runRecorders {
		r.RecordRun(utils.GetCtxWithTimeout(time.Second*5), pseudoStage, queryResults)
	}

	log.Info().Int("file_loaded", len(queryResults)).Send()
}

func processPath(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return processFile(path)
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fullPath := filepath.Join(path, entry.Name())
		if err = processFile(fullPath); err != nil {
			return err
		}
	}
	return nil
}

func processFile(path string) (err error) {
	bytes, ioErr := os.ReadFile(path)
	if ioErr != nil {
		return ioErr
	}
	queryInfo := new(query_json.QueryInfo)
	if unmarshalErr := json.Unmarshal(bytes, queryInfo); unmarshalErr != nil {
		return unmarshalErr
	}

	queryResult := &stage.QueryResult{
		StageId: pseudoStage.Id,
		Query: &stage.Query{
			Text:             queryInfo.Query,
			ColdRun:          true,
			ExpectedRowCount: -1,
		},
		QueryId:   queryInfo.QueryId,
		InfoUrl:   queryInfo.Self,
		RowCount:  int(queryInfo.QueryStats.OutputPositions),
		StartTime: *queryInfo.QueryStats.CreateTime,
		EndTime:   queryInfo.QueryStats.EndTime,
	}
	if queryInfo.ErrorCode != nil {
		queryResult.QueryError = fmt.Errorf(*queryInfo.ErrorCode.Name)
	}
	if queryResult.StartTime.Before(pseudoStage.States.RunStartTime) {
		pseudoStage.States.RunStartTime = queryResult.StartTime
	}
	if queryResult.EndTime != nil {
		dur := queryResult.EndTime.Sub(queryResult.StartTime)
		queryResult.Duration = &dur
		if queryResult.EndTime.After(pseudoStage.States.RunFinishTime) {
			pseudoStage.States.RunFinishTime = *queryResult.EndTime
		}
	}

	for _, r := range runRecorders {
		r.RecordQuery(utils.GetCtxWithTimeout(time.Second*5), pseudoStage, queryResult)
	}

	queryResults = append(queryResults, queryResult)

	if mysqlDb != nil {
		// OutputStage is in a tree structure, and we need to flatten it for its ORM to be correctly parsed.
		e := queryInfo.PrepareForInsert()
		if e == nil {
			e = utils.SqlInsertObject(mysqlDb, queryInfo,
				"presto_query_creation_info",
				"presto_query_operator_stats",
				"presto_query_plans",
				"presto_query_stage_stats",
				"presto_query_statistics",
			)
		}
		if e != nil {
			log.Error().Err(e).Msg("failed to insert")
			return e
		}
	}
	return nil
}

func registerRunRecorder(r stage.RunRecorder) {
	if r == nil || reflect.ValueOf(r).IsNil() {
		return
	}
	runRecorders = append(runRecorders, r)
}
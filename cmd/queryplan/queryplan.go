// Copyright 2024.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package queryplan

import (
	"pbench/cmd"

	"github.com/spf13/cobra"
)

// queryplanCmd represents the run command
var queryplanCmd = &cobra.Command{
	Use:                   `queryplan [flags] [list of root-level benchmark stage JSON files]`,
	Short:                 "Parse query plan",
	Long:                  `Read a CSV file, parse the "query plan" column, and write the JOIN information into a JSON file`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(0),
	Run:                   Run,
}

func init() {
	cmd.RootCmd.AddCommand(queryplanCmd)
	queryplanCmd.Flags().StringVarP(&csvFile, "csv", "f", "", `CSV file to read`)
	queryplanCmd.Flags().IntVarP(&queryPlanColumn, "column", "c", 0, `The column index for the Query Plans in the CSV file(index starts with 0)`)
	queryplanCmd.Flags().StringVarP(&output, "output", "o", "queryplan.json", "Output JSON file")
	queryplanCmd.Flags().BoolVarP(&hasHeader, "has-header", "s", true, "contain the header line or not")
}

package stage

type Query struct {
	Text             string
	File             *string
	Index            int
	BatchSize        int
	ColdRun          bool
	SequenceNo       int
	ExpectedRowCount int
}

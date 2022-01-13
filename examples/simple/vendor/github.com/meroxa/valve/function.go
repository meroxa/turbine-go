package valve

type Function interface {
	Process(r []Record) ([]Record, []RecordWithError)
}

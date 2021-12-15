package valve

type Function interface {
	Process(r []Record) ([]Record, []RecordWithError)
}

func Process(s []Record, fn Function) ([]Record, []RecordWithError) {
	return fn.Process(s)
}
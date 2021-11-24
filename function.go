package valve

type Function interface {
	Process(r []Record) ([]Record, error)
}

func Process(s []Record, fn Function) ([]Record, []RecordWithError) {
	return nil, nil
}
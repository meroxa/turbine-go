package app

type Function interface {
	Process(r []Record) []Record
}

package valve

type Resource struct {
	Name    string
	ID      int
}

type ResourceConfig struct {
	Field string
	Value string
}

type ResourceConfigs []ResourceConfig

func (r Resource) Records(collection string, cfg ResourceConfigs) ([]Record, error) {
	if local {
		// TODO: retrieve records from fixtures
		return ReadFixtures(r.Name, collection)
	}
	return []Record{}, nil
}

func (r Resource) Write(rr []Record, collection string, cfg ResourceConfigs) error {
	return nil
}

func Resources(name string) (Resource, error) {
	r, err := C.GetResource(name)
	if err != nil {
		return Resource{}, err
	}

	return Resource{
		Name: r.Name,
		ID:   r.ID,
	}, nil
}

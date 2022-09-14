package common

type Counter struct {
	Count int `json:"count"`
}

func (r *Counter) Record(count int) Counter {
	r.Count = count
	return *r
}

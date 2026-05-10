package reducer

type Reducer struct {
	dirName string
}

func NewReducer(dirName string) *Reducer {
	return &Reducer{dirName: dirName}
}

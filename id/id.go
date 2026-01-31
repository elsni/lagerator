package id

type IdMaker struct {
	LastId uint32 `json:"lastid"`
}

var IdSource = NewIdMaker()

// NewIdMaker returns an IdMaker with LastId initialized to zero.
func NewIdMaker() IdMaker {
	return IdMaker{LastId: 0}
}

// GetNewId increments and returns the next id.
func (id *IdMaker) GetNewId() uint32 {
	id.LastId += 1
	return id.LastId
}

// SetLastId sets the current last id.
func (id *IdMaker) SetLastId(lid uint32) {
	id.LastId = lid
}

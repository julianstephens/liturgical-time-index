package plan

type RbRef struct {
	Chapter    string `validate:"required,alphanum"                 json:"chapter"`
	StartVerse int    `validate:"required,gte=1"                    json:"start_verse"`
	EndVerse   int    `validate:"required,gte=1,gtfield=StartVerse" json:"end_verse,omitempty"`
}

func NewRbRef(rbStr string) (*RbRef, error) {
	ref, err := parseRbRef(rbStr)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

func parseRbRef(rbStr string) (*RbRef, error) {
	// TODO: Implement parsing logic to convert rbStr into RbRef struct

	return &RbRef{
		Chapter:    "ExampleChapter",
		StartVerse: 1,
		EndVerse:   5,
	}, nil
}

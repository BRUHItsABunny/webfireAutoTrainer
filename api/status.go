package api

// Make x ENUM(Initialized,Progressed)
//
//go:generate go-enum -f=$GOFILE --nocase --flag --names
type LessonStatus int

func (x LessonStatus) ToParameterString() string {
	switch x {
	case LessonStatusInitialized:
		return "I"
	default:
		return "P"
	}
}

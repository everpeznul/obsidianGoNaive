package use_case

type Noter interface {
	FindFather() (Noter, error)
	FindAncestor() (Noter, error)
	FindFounder() (Noter, error)
}

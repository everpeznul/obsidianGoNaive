package use_case

type Noter interface {
	FindFather() (string, error)
	FindAncestor() (string, error)
	FindFounder() (string, error)
}

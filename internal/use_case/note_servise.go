package use_case

type Noter interface {
	FindFather() string
	FindAncestor() string
	FindFounder() string
}

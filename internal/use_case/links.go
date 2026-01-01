package use_case

type Linker struct {
}

func (l *Linker) Format(n Noter) []string {

	founder, _ := n.FindFounder()
	ancestor, _ := n.FindAncestor()
	father, _ := n.FindFather()

	return []string{founder, ancestor, father}
}

package use_case

type Linker struct {
}

func (l *Linker) Format(n Noter) []string {

	founder := n.FindFounder()
	ancestor := n.FindAncestor()
	father := n.FindFather()

	return []string{founder, ancestor, father}
}

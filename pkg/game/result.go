package game

type result int

const (
	unknown result = iota
	notIncluded
	wrongPos
	correct
)

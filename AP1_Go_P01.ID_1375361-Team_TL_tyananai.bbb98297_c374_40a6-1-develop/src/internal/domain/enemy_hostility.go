package domain

type HostilityLevel int

const (
	HostilityLow HostilityLevel = iota
	HostilityAverage
	HostilityHigh
)

func (h HostilityLevel) Radius() int {
	switch h {
	case HostilityLow:
		return 4
	case HostilityAverage:
		return 6
	case HostilityHigh:
		return 8
	default:
		return 6
	}
}

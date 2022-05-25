package stat

import (
	"testing"
)

func TestStat(t *testing.T) {
	s := NewStat("zeromicro/go-zero")
	// s.Repository()
	s.LatestMonthStars()
	// fmt.Println(s.ContributorCount())
}

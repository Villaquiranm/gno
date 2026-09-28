package connect4

import (
	"strconv"
	"strings"
	"testing"

	"gno.land/p/nt/testutils/v0"
	"gno.land/p/nt/uassert/v0"
)

func TestLeaderboardTopTen(t *testing.T) {
	reset()
	for i := 1; i <= 12; i++ {
		statsOf(testutils.TestAddress("p" + strconv.Itoa(i))).Wins = i
	}
	out := leaderboard(func(s *Stats) int64 { return int64(s.Wins) }, func(v int64) string { return strconv.FormatInt(v, 10) })
	uassert.True(t, strings.HasPrefix(strings.Split(out, "\n")[2], "| 1 | "+testutils.TestAddress("p12").String()+" | 12 |"))
	uassert.True(t, strings.Contains(out, "| 10 | "+testutils.TestAddress("p3").String()+" | 3 |"))
	uassert.False(t, strings.Contains(out, testutils.TestAddress("p2").String()))
}

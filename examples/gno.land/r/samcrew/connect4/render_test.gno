package connect4

import (
	"strconv"
	"strings"
	"testing"
	"time"

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

func TestGamePageShowsActionsForEveryState(cur realm, t *testing.T) {
	reset()
	contains := func(out string, want ...string) {
		for _, w := range want {
			uassert.True(t, strings.Contains(out, w), "missing "+w+" in:\n"+out)
		}
	}

	testing.SetRealm(testing.NewUserRealm(alice))
	pay(stake)
	id := Offer(cross(cur), "", 10, commit())
	page := "game/" + strconv.Itoa(id)
	contains(Render(page), "func=Accept", ".send=10000000ugnot", "func=Cancel", "Open offer")

	g := mustGet(id)
	g.ExpiresAt = time.Now().Add(-time.Second)
	out := Render(page)
	contains(out, "expired", "func=Cancel")
	uassert.False(t, strings.Contains(out, "func=Accept"))

	g, p := accepted(cur)
	page = "game/" + strconv.Itoa(g.ID)
	contains(Render(page), "Waiting for", "func=Reveal", "func=Resign")

	testing.SetRealm(testing.NewUserRealm(alice))
	pay(0)
	Reveal(cross(cur), g.ID, p)
	contains(Render(page), "func=Play", "column=7", "func=Resign", "Turn:")

	testing.SetRealm(testing.NewUserRealm(alice))
	pay(0)
	Resign(cross(cur), g.ID)
	contains(Render(page), "Winner", bob.String(), gnot(2*stake-fee))

	for status, want := range map[Status]string{
		Draw:      "both stakes refunded",
		Void:      "both stakes refunded",
		Cancelled: "refunded to the creator",
	} {
		g.Status = status
		contains(Render(page), want)
	}
}

func TestLobbyLinksEveryGameToItsPage(cur realm, t *testing.T) {
	reset()
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(stake)
	id := Offer(cross(cur), "", 10, commit())
	g, _ := accepted(cur)
	// Render as if deployed elsewhere: links must follow the realm's real path.
	testing.SetRealm(testing.NewCodeRealm("gno.land/r/nym-mikecito001/connect4"))
	out := Render("")
	base := "(/r/nym-mikecito001/connect4:"
	uassert.True(t, strings.Contains(out, base+"game/"+strconv.Itoa(id)+")"), out)
	uassert.True(t, strings.Contains(out, base+"game/"+strconv.Itoa(g.ID)+")"), out)
	uassert.True(t, strings.Contains(out, base+"new)"), out)
	uassert.False(t, strings.Contains(out, "samcrew"), "links must not hardcode the realm path")
}

func TestNewOfferPage(t *testing.T) {
	reset()
	c := sha256Hex("some four random words")

	out := Render("new")
	uassert.True(t, strings.Contains(out, `<gno-form path="new">`))
	uassert.False(t, strings.Contains(out, "func=Offer"))

	out = Render("new?stake=12.5&validFor=10&opponent=" + bob.String() + "&commitment=" + c)
	for _, want := range []string{"func=Offer", ".send=12500000ugnot", "validFor=10", "opponent=" + bob.String(), "commitment=" + c, "12.5 GNOT"} {
		uassert.True(t, strings.Contains(out, want), "missing "+want+" in:\n"+out)
	}

	out = Render("new?stake=3&commitment=" + c) // validFor defaults, opponent optional
	uassert.True(t, strings.Contains(out, ".send=3000000ugnot"))
	uassert.True(t, strings.Contains(out, "validFor=15"))

	for _, q := range []string{
		"stake=%3Cscript%3E&validFor=10&commitment=" + c, // not a number
		"stake=0.5&validFor=10&commitment=" + c,          // below minimum
		"stake=1.0000001&validFor=10&commitment=" + c,    // more than 6 decimals
		"stake=5&validFor=61&commitment=" + c,
		"stake=5&validFor=10&opponent=%3Cscript%3E&commitment=" + c,
		"stake=5&validFor=10&commitment=%3Cscript%3E",
		"stake=5&validFor=10",
	} {
		out := Render("new?" + q)
		uassert.False(t, strings.Contains(out, "func=Offer"), q)
		uassert.True(t, strings.Contains(out, "Fix the following"), q)
		uassert.False(t, strings.Contains(out, "<script>"), "form input echoed: "+q)
	}
}

package connect4

import (
	"chain"
	"chain/banker"
	"testing"

	"gno.land/p/nt/avl/v0"
	"gno.land/p/nt/testutils/v0"
	"gno.land/p/nt/uassert/v0"
)

const stake = int64(10_000_000)

var (
	alice   = testutils.TestAddress("alice")
	bob     = testutils.TestAddress("bob")
	carol   = testutils.TestAddress("carol")
	admin   = address("g1cvr48r7l7lkmvp77cr6zg2zhu26jgfwr0y8pew")
	pkgAddr = chain.PackageAddress("gno.land/r/samcrew/connect4")
)

func reset() {
	games, active, stats = avl.Tree{}, avl.Tree{}, avl.Tree{}
	nextID, fee, feesCollected = 0, 100_000, 0
}

// pay sets the ugnot sent with the next call and credits the realm with it,
// mirroring what the chain does. testing.SetRealm must be called in the same
// frame that makes the cross call, so tests call it inline rather than here.
func pay(amount int64) {
	if amount == 0 {
		testing.SetOriginSend(nil)
		return
	}
	coins := chain.Coins{{"ugnot", amount}}
	testing.SetOriginSend(coins)
	testing.IssueCoins(pkgAddr, coins)
}

func balance(addr address) int64 {
	return banker.NewReadonlyBanker().GetCoin(addr, "ugnot")
}

// started returns a Playing game between alice (creator) and bob.
func started(cur realm) *Game {
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(stake)
	id := Offer(cross(cur), "", 10)
	testing.SetRealm(testing.NewUserRealm(bob))
	pay(stake)
	Accept(cross(cur), id)
	return mustGet(id)
}

// move plays column for whoever is on turn.
func move(cur realm, g *Game, column int) {
	testing.SetRealm(testing.NewUserRealm(g.player(g.Turn)))
	pay(0)
	Play(cross(cur), g.ID, column)
}

// expireClock advances block time past the 90s move clock (SkipHeights is 5s/block).
func expireClock() { testing.SkipHeights(19) }

func TestWinPaysPotMinusFee(cur realm, t *testing.T) {
	reset()
	g := started(cur)
	first := g.player(g.Turn)
	before := balance(first)
	for _, c := range []int{1, 2, 1, 2, 1, 2, 1} { // first mover stacks column 1
		move(cur, g, c)
	}
	uassert.Equal(t, "won", g.Status.String())
	uassert.Equal(t, first, g.Winner)
	uassert.Equal(t, 2*stake-100_000, balance(first)-before)
	uassert.Equal(t, int64(100_000), feesCollected)
	uassert.Equal(t, 1, statsOf(first).Wins)
	uassert.Equal(t, 2*stake-100_000, statsOf(first).GnotWon)
	uassert.Equal(t, 1, statsOf(g.player(3-g.Turn)).Games)
	uassert.Equal(t, 0, active.Size())
}

func TestOfferRejects(cur realm, t *testing.T) {
	reset()
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(minStake - 1)
	uassert.AbortsContains(t, cur, "stake below minimum", func() { Offer(cross(cur), "", 10) })
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(stake)
	uassert.AbortsContains(t, cur, "validFor", func() { Offer(cross(cur), "", 0) })
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(stake)
	uassert.AbortsContains(t, cur, "validFor", func() { Offer(cross(cur), "", 61) })
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(stake)
	uassert.AbortsContains(t, cur, "yourself", func() { Offer(cross(cur), alice, 10) })
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(0)
	uassert.AbortsContains(t, cur, "send exactly one", func() { Offer(cross(cur), "", 10) })
	testing.SetRealm(testing.NewUserRealm(alice))
	testing.SetOriginSend(chain.Coins{{"ugnot", stake}, {"foo", 1}})
	uassert.AbortsContains(t, cur, "send exactly one", func() { Offer(cross(cur), "", 10) })
}

func TestAcceptRejects(cur realm, t *testing.T) {
	reset()
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(stake)
	id := Offer(cross(cur), bob, 1)
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(stake)
	uassert.AbortsContains(t, cur, "own offer", func() { Accept(cross(cur), id) })
	testing.SetRealm(testing.NewUserRealm(carol))
	pay(stake)
	uassert.AbortsContains(t, cur, "reserved", func() { Accept(cross(cur), id) })
	testing.SetRealm(testing.NewUserRealm(bob))
	pay(stake - 1)
	uassert.AbortsContains(t, cur, "must send exactly", func() { Accept(cross(cur), id) })
	testing.SkipHeights(13) // 65s > 1 minute
	testing.SetRealm(testing.NewUserRealm(bob))
	pay(stake)
	uassert.AbortsContains(t, cur, "expired", func() { Accept(cross(cur), id) })
	testing.SetRealm(testing.NewUserRealm(bob))
	pay(stake)
	uassert.AbortsContains(t, cur, "not found", func() { Accept(cross(cur), 99) })
}

func TestPlayRejects(cur realm, t *testing.T) {
	reset()
	g := started(cur)
	testing.SetRealm(testing.NewUserRealm(g.player(3 - g.Turn)))
	pay(0)
	uassert.AbortsContains(t, cur, "not your turn", func() { Play(cross(cur), g.ID, 1) })
	testing.SetRealm(testing.NewUserRealm(carol))
	pay(0)
	uassert.AbortsContains(t, cur, "not your turn", func() { Play(cross(cur), g.ID, 1) })
	testing.SetRealm(testing.NewUserRealm(g.player(g.Turn)))
	pay(0)
	uassert.AbortsContains(t, cur, "column must be", func() { Play(cross(cur), g.ID, 0) })
	testing.SetRealm(testing.NewUserRealm(g.player(g.Turn)))
	pay(0)
	uassert.AbortsContains(t, cur, "column must be", func() { Play(cross(cur), g.ID, 8) })
	expireClock()
	testing.SetRealm(testing.NewUserRealm(g.player(g.Turn)))
	pay(0)
	uassert.AbortsContains(t, cur, "clock ran out", func() { Play(cross(cur), g.ID, 1) })
}

func TestFullColumnRejected(cur realm, t *testing.T) {
	reset()
	g := started(cur)
	for i := 0; i < rows; i++ {
		move(cur, g, 4) // alternating players never make 4 in a column
	}
	testing.SetRealm(testing.NewUserRealm(g.player(g.Turn)))
	pay(0)
	uassert.AbortsContains(t, cur, "column is full", func() { Play(cross(cur), g.ID, 4) })
}

func TestTimeoutBeforeFirstMoveIsVoid(cur realm, t *testing.T) {
	reset()
	g := started(cur)
	testing.SetRealm(testing.NewUserRealm(carol))
	pay(0)
	uassert.AbortsContains(t, cur, "has not run out", func() { ClaimTimeout(cross(cur), g.ID) })
	a, b := balance(alice), balance(bob)
	expireClock()
	testing.SetRealm(testing.NewUserRealm(carol))
	pay(0) // anyone may settle
	ClaimTimeout(cross(cur), g.ID)
	uassert.Equal(t, "void", g.Status.String())
	uassert.Equal(t, stake, balance(alice)-a)
	uassert.Equal(t, stake, balance(bob)-b)
	uassert.Equal(t, int64(0), feesCollected)
	uassert.Equal(t, 0, stats.Size())
}

func TestTimeoutAfterMoveForfeits(cur realm, t *testing.T) {
	reset()
	g := started(cur)
	mover := g.player(g.Turn)
	move(cur, g, 1)
	expireClock()
	testing.SetRealm(testing.NewUserRealm(mover))
	pay(0)
	ClaimTimeout(cross(cur), g.ID)
	uassert.Equal(t, "won", g.Status.String())
	uassert.Equal(t, mover, g.Winner)
}

func TestResign(cur realm, t *testing.T) {
	reset()
	g := started(cur)
	testing.SetRealm(testing.NewUserRealm(carol))
	pay(0)
	uassert.AbortsContains(t, cur, "not a player", func() { Resign(cross(cur), g.ID) })
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(0)
	Resign(cross(cur), g.ID)
	uassert.Equal(t, bob, g.Winner)
}

func TestCancel(cur realm, t *testing.T) {
	reset()
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(stake)
	id := Offer(cross(cur), "", 1)
	testing.SetRealm(testing.NewUserRealm(bob))
	pay(0)
	uassert.AbortsContains(t, cur, "only the creator", func() { Cancel(cross(cur), id) })
	testing.SkipHeights(13) // expired offers can still be cancelled
	before := balance(alice)
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(0)
	Cancel(cross(cur), id)
	uassert.Equal(t, "cancelled", mustGet(id).Status.String())
	uassert.Equal(t, stake, balance(alice)-before)

	g := started(cur)
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(0)
	uassert.AbortsContains(t, cur, "not open", func() { Cancel(cross(cur), g.ID) })
}

func TestDrawRefundsBoth(cur realm, t *testing.T) {
	reset()
	g := started(cur)
	// Fill every cell but the top of column 7 with a pattern that has no
	// four in a row; player 1 owns the last cell.
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			g.Board[c][r] = byte(1 + (c/2+r)%2)
		}
	}
	g.Board[6][5] = 0
	g.Moves, g.Turn = cols*rows-1, 1
	a, b := balance(alice), balance(bob)
	move(cur, g, 7)
	uassert.Equal(t, "draw", g.Status.String())
	uassert.Equal(t, stake, balance(alice)-a)
	uassert.Equal(t, stake, balance(bob)-b)
	uassert.Equal(t, 1, statsOf(alice).Draws)
}

func TestFeeSnapshotAndWithdraw(cur realm, t *testing.T) {
	reset()
	testing.SetRealm(testing.NewUserRealm(carol))
	pay(0)
	uassert.AbortsContains(t, cur, "not owner", func() { SetFee(cross(cur), 0) })
	g := started(cur) // offered at the default 0.1 GNOT fee
	testing.SetRealm(testing.NewUserRealm(admin))
	pay(0)
	uassert.AbortsContains(t, cur, "negative", func() { SetFee(cross(cur), -1) })
	SetFee(cross(cur), 5_000_000)
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(0)
	Resign(cross(cur), g.ID)
	uassert.Equal(t, int64(100_000), feesCollected)

	before := balance(carol)
	testing.SetRealm(testing.NewUserRealm(admin))
	pay(0)
	WithdrawFees(cross(cur), carol)
	uassert.Equal(t, int64(100_000), balance(carol)-before)
	uassert.Equal(t, int64(0), feesCollected)

	// A raised fee applies to new offers, which must stake more than it.
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(5_000_000)
	uassert.AbortsContains(t, cur, "exceed the house fee", func() { Offer(cross(cur), "", 10) })
}

func TestFinishedGameRejects(cur realm, t *testing.T) {
	reset()
	g := started(cur)
	testing.SetRealm(testing.NewUserRealm(alice))
	pay(0)
	Resign(cross(cur), g.ID)
	testing.SetRealm(testing.NewUserRealm(bob))
	pay(0)
	uassert.AbortsContains(t, cur, "not in progress", func() { Play(cross(cur), g.ID, 1) })
	uassert.AbortsContains(t, cur, "not in progress", func() { Resign(cross(cur), g.ID) })
	expireClock()
	testing.SetRealm(testing.NewUserRealm(bob))
	pay(0)
	uassert.AbortsContains(t, cur, "not in progress", func() { ClaimTimeout(cross(cur), g.ID) })
	testing.SetRealm(testing.NewUserRealm(carol))
	pay(stake)
	uassert.AbortsContains(t, cur, "not open", func() { Accept(cross(cur), g.ID) })
	testing.SetRealm(testing.NewUserRealm(carol))
	pay(stake)
	uassert.AbortsContains(t, cur, "invalid opponent", func() { Offer(cross(cur), "notanaddress", 10) })
}

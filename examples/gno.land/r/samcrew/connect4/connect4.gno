// Package connect4 is a staked Connect 4 game. A creator posts an offer with
// a ugnot stake; anyone (or a named opponent) accepts with the same stake;
// the winner takes the pot minus a flat house fee.
//
// Every move has a 90s clock measured in block time. Nothing on-chain fires
// by itself, so an expired clock is settled by anyone calling ClaimTimeout:
// if the first mover never moved the game is void and both stakes are
// refunded, otherwise the player who ran out of time forfeits.
package connect4

import (
	"chain"
	"chain/banker"
	"chain/runtime"
	"chain/runtime/unsafe"
	"crypto/sha256"
	"strconv"
	"strings"
	"time"

	"gno.land/p/nt/avl/v0"
	"gno.land/p/nt/ownable/v0"
)

const (
	moveTimeout = 90 * time.Second
	maxOfferTTL = 60        // minutes
	minStake    = 1_000_000 // ugnot
	denom       = "ugnot"
)

type Status int

const (
	Open Status = iota
	Playing
	Won
	Draw
	Void
	Cancelled
)

func (s Status) String() string {
	return [...]string{"open", "playing", "won", "draw", "void", "cancelled"}[s]
}

type Game struct {
	ID        int
	Creator   address
	Opponent  address // optional; if set, only this address may accept
	Acceptor  address
	Stake     int64 // ugnot, per player
	Fee       int64 // ugnot house fee, fixed at offer time
	ExpiresAt time.Time

	Board    board
	Turn     byte // 1 = creator, 2 = acceptor
	Moves    int
	LastMove time.Time
	Status   Status
	Winner   address
}

type Stats struct {
	Games, Wins, Draws int
	GnotWon            int64 // total ugnot received from wins
}

var (
	games         avl.Tree // key(id) -> *Game
	active        avl.Tree // key(id) -> *Game, only Open and Playing
	stats         avl.Tree // address -> *Stats
	nextID        int
	fee           int64 = 100_000 // ugnot per decided game
	feesCollected int64
	owner         = ownable.NewWithAddress("g1cvr48r7l7lkmvp77cr6zg2zhu26jgfwr0y8pew")
)

// Offer posts a game open for validFor minutes. The stake is the ugnot sent
// with the call. A non-empty opponent makes the offer private.
func Offer(cur realm, opponent address, validFor int) int {
	stake := received(cur)
	caller := cur.Previous().Address()
	if stake < minStake {
		panic("stake below minimum of " + strconv.Itoa(minStake) + denom)
	}
	if stake <= fee {
		panic("stake must exceed the house fee of " + strconv.FormatInt(fee, 10) + denom)
	}
	if validFor < 1 || validFor > maxOfferTTL {
		panic("validFor must be between 1 and 60 minutes")
	}
	if opponent != "" && !opponent.IsValid() {
		panic("invalid opponent address")
	}
	if opponent == caller {
		panic("cannot challenge yourself")
	}
	nextID++
	g := &Game{
		ID:        nextID,
		Creator:   caller,
		Opponent:  opponent,
		Stake:     stake,
		Fee:       fee,
		ExpiresAt: time.Now().Add(time.Duration(validFor) * time.Minute),
	}
	games.Set(key(g.ID), g)
	active.Set(key(g.ID), g)
	return g.ID
}

// Accept joins an open offer by sending exactly its stake. The first mover
// is picked from the accepting block and has 90s to play.
func Accept(cur realm, id int) {
	stake := received(cur)
	caller := cur.Previous().Address()
	g := mustGet(id)
	if g.Status != Open {
		panic("game is not open")
	}
	if !time.Now().Before(g.ExpiresAt) {
		panic("offer expired")
	}
	if caller == g.Creator {
		panic("cannot accept your own offer")
	}
	if g.Opponent != "" && caller != g.Opponent {
		panic("offer is reserved for another player")
	}
	if stake != g.Stake {
		panic("must send exactly " + strconv.FormatInt(g.Stake, 10) + denom)
	}
	g.Acceptor = caller
	g.Status = Playing
	g.LastMove = time.Now()
	g.Turn = firstMover(g)
}

// Play drops a piece in column (1-7).
func Play(cur realm, id, column int) {
	g := mustGet(id)
	if g.Status != Playing {
		panic("game is not in progress")
	}
	if cur.Previous().Address() != g.player(g.Turn) {
		panic("not your turn")
	}
	if time.Now().After(g.deadline()) {
		panic("move clock ran out")
	}
	if column < 1 || column > cols {
		panic("column must be between 1 and 7")
	}
	r := g.Board.drop(column-1, g.Turn)
	if r < 0 {
		panic("column is full")
	}
	g.Moves++
	switch {
	case g.Board.wins(column-1, r):
		settleWin(cur, g, g.player(g.Turn))
	case g.Moves == cols*rows:
		settleDraw(cur, g)
	default:
		g.Turn = 3 - g.Turn
		g.LastMove = time.Now()
	}
}

// ClaimTimeout settles a game whose move clock ran out. Anyone may call it.
func ClaimTimeout(cur realm, id int) {
	g := mustGet(id)
	if g.Status != Playing {
		panic("game is not in progress")
	}
	if !time.Now().After(g.deadline()) {
		panic("move clock has not run out")
	}
	if g.Moves == 0 {
		settleVoid(cur, g)
		return
	}
	settleWin(cur, g, g.player(3-g.Turn))
}

// Resign concedes a game in progress.
func Resign(cur realm, id int) {
	g := mustGet(id)
	if g.Status != Playing {
		panic("game is not in progress")
	}
	switch cur.Previous().Address() {
	case g.Creator:
		settleWin(cur, g, g.Acceptor)
	case g.Acceptor:
		settleWin(cur, g, g.Creator)
	default:
		panic("not a player in this game")
	}
}

// Cancel withdraws an open offer, expired or not, and refunds the stake.
func Cancel(cur realm, id int) {
	g := mustGet(id)
	if g.Status != Open {
		panic("game is not open")
	}
	if cur.Previous().Address() != g.Creator {
		panic("only the creator can cancel")
	}
	finish(g, Cancelled)
	send(cur, g.Creator, g.Stake)
}

// SetFee changes the flat house fee, in ugnot, for new offers. Games already
// offered keep the fee they were created with.
func SetFee(cur realm, amount int64) {
	owner.AssertOwnedBy(cur.Previous().Address())
	if amount < 0 {
		panic("fee must not be negative")
	}
	fee = amount
}

// WithdrawFees sends all collected fees to to.
func WithdrawFees(cur realm, to address) {
	owner.AssertOwnedBy(cur.Previous().Address())
	amount := feesCollected
	feesCollected = 0
	send(cur, to, amount)
}

func TransferOwnership(cur realm, newOwner address) {
	if err := owner.TransferOwnership(0, cur, newOwner); err != nil {
		panic(err)
	}
}

// received returns the ugnot sent with a direct user call. OriginSend only
// describes a real payment to this realm when the caller is an EOA using
// maketx call; see docs/resources/effective-gno.md.
func received(cur realm) int64 {
	if !cur.Previous().IsUserCall() {
		panic("only user-call (maketx call) accepted")
	}
	sent := unsafe.OriginSend()
	if len(sent) != 1 || sent[0].Denom != denom {
		panic("send exactly one " + denom + " coin")
	}
	return sent[0].Amount
}

func firstMover(g *Game) byte {
	seed := strings.Join([]string{
		strconv.FormatInt(runtime.ChainHeight(), 10),
		strconv.FormatInt(time.Now().UnixNano(), 10),
		strconv.Itoa(g.ID),
		g.Creator.String(),
		g.Acceptor.String(),
	}, "|")
	h := sha256.Sum256([]byte(seed))
	return 1 + h[0]&1
}

func settleWin(cur realm, g *Game, winner address) {
	loser := g.Creator
	if winner == g.Creator {
		loser = g.Acceptor
	}
	pot := 2 * g.Stake
	finish(g, Won)
	g.Winner = winner
	feesCollected += g.Fee
	w := statsOf(winner)
	w.Games++
	w.Wins++
	w.GnotWon += pot - g.Fee
	statsOf(loser).Games++
	send(cur, winner, pot-g.Fee)
}

func settleDraw(cur realm, g *Game) {
	finish(g, Draw)
	for _, p := range []address{g.Creator, g.Acceptor} {
		s := statsOf(p)
		s.Games++
		s.Draws++
		send(cur, p, g.Stake)
	}
}

func settleVoid(cur realm, g *Game) {
	finish(g, Void)
	send(cur, g.Creator, g.Stake)
	send(cur, g.Acceptor, g.Stake)
}

func finish(g *Game, s Status) {
	g.Status = s
	active.Remove(key(g.ID))
}

func send(cur realm, to address, amount int64) {
	if amount == 0 {
		return
	}
	banker.NewBanker(banker.BankerTypeRealmSend, cur).
		SendCoins(cur.Address(), to, chain.Coins{{denom, amount}})
}

func (g *Game) player(turn byte) address {
	if turn == 1 {
		return g.Creator
	}
	return g.Acceptor
}

func (g *Game) deadline() time.Time { return g.LastMove.Add(moveTimeout) }

func statsOf(addr address) *Stats {
	if v := stats.Get(addr.String()); v != nil {
		return v.(*Stats)
	}
	s := &Stats{}
	stats.Set(addr.String(), s)
	return s
}

// key zero-pads ids so avl iteration order matches creation order.
func key(id int) string {
	s := strconv.Itoa(id)
	return strings.Repeat("0", 10-len(s)) + s
}

func get(id int) *Game {
	if v := games.Get(key(id)); v != nil {
		return v.(*Game)
	}
	return nil
}

func mustGet(id int) *Game {
	g := get(id)
	if g == nil {
		panic("game not found")
	}
	return g
}

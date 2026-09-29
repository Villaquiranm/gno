// PKGPATH: gno.land/r/demo/main

package main

import (
	"chain"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"gno.land/p/nt/testutils/v0"
	"gno.land/r/samcrew/connect4"
)

func commit(passphrase string) string {
	h := sha256.Sum256([]byte(passphrase))
	return hex.EncodeToString(h[:])
}

func main(cur realm) {
	alice := testutils.TestAddress("alice")
	bob := testutils.TestAddress("bob")
	carol := testutils.TestAddress("carol")
	c4 := chain.PackageAddress("gno.land/r/samcrew/connect4")
	stake := chain.Coins{{"ugnot", 5_000_000}}

	testing.SetRealm(testing.NewUserRealm(alice))
	testing.SetOriginSend(stake)
	testing.IssueCoins(c4, stake)
	g1 := connect4.Offer(cross(cur), "", 15, commit("alice four random words"))

	testing.SetRealm(testing.NewUserRealm(carol))
	testing.SetOriginSend(stake)
	testing.IssueCoins(c4, stake)
	connect4.Offer(cross(cur), bob, 5, commit("carol four random words"))

	testing.SetRealm(testing.NewUserRealm(bob))
	testing.SetOriginSend(stake)
	testing.IssueCoins(c4, stake)
	connect4.Accept(cross(cur), g1)

	// Render runs as the realm itself, as it does under gnoweb.
	testing.SetRealm(testing.NewCodeRealm("gno.land/r/samcrew/connect4"))
	println(connect4.Render(""))
	println("-----")
	println(connect4.Render("game/1"))
	println("-----")
	println(connect4.Render("game/2"))
	println("-----")
	println(connect4.Render("new?stake=12.5&validFor=10&opponent=&commitment=5fb8d3d0c4f3b0fe7ab6ab1f7eda42c6e2de0c3c1d4a6c5e3d0f2f5d0e4f5a6b"))
	println("-----")
	println(connect4.Render("game/x"))
	println("-----")
	println(connect4.Render("game/99999999999"))
}

// Output:
// # Connect 4
//
// Both players stake the same amount; the winner takes the pot minus a 0.1 GNOT fee. 90s per move: run out and you forfeit. If the first mover never moves, both stakes are refunded.
//
// **[Post an offer](/r/samcrew/connect4:new)** · open a game below to play it
//
// ## Open offers
//
// | Game | Stake | Creator | For | Expires in | |
// |---|---|---|---|---|---|
// | [#2](/r/samcrew/connect4:game/2) | 5 GNOT | g1vdshymmvta047h6lta047h6lta047h6l2asz94 | g1vfhkyh6lta047h6lta047h6lta047h6l03vdhu | 5m0s | [Accept](/r/samcrew/connect4$help&func=Accept&.send=5000000ugnot&id=2) |
//
// ## Live games
//
// | Game | Pot | Turn | Time left |
// |---|---|---|---|
// | [#1](/r/samcrew/connect4:game/1) | 10 GNOT | g1v9kxjcm9ta047h6lta047h6lta047h6lzd40gh (reveal) | 1m30s |
//
// ## Most games won
//
// No finished games yet.
//
// ## Most gnot won
//
// No finished games yet.
//
//
// -----
// # Game #1 — playing
//
// [← Lobby](/r/samcrew/connect4)
//
// 🔴 g1v9kxjcm9ta047h6lta047h6lta047h6lzd40gh
// 🟡 g1vfhkyh6lta047h6lta047h6lta047h6l03vdhu
// Stake: 5 GNOT each · fee 0.1 GNOT
//
// 1️⃣2️⃣3️⃣4️⃣5️⃣6️⃣7️⃣
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
//
// Waiting for the creator to reveal the passphrase.
//
// Time left: 1m30s
//
// **[Reveal](/r/samcrew/connect4$help&func=Reveal&id=1&passphrase=)** (creator)
//
// [Resign](/r/samcrew/connect4$help&func=Resign&id=1)
//
// -----
// # Game #2 — open
//
// [← Lobby](/r/samcrew/connect4)
//
// 🔴 g1vdshymmvta047h6lta047h6lta047h6l2asz94
// 🟡 waiting for opponent
// Stake: 5 GNOT each · fee 0.1 GNOT
//
// 1️⃣2️⃣3️⃣4️⃣5️⃣6️⃣7️⃣
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
//
// Open offer for g1vfhkyh6lta047h6lta047h6lta047h6l03vdhu, expires in 5m0s.
//
// **[Accept for 5 GNOT](/r/samcrew/connect4$help&func=Accept&.send=5000000ugnot&id=2)**
//
// [Cancel offer](/r/samcrew/connect4$help&func=Cancel&id=2) (creator)
//
// -----
// # Post an offer
//
// [← Lobby](/r/samcrew/connect4)
//
// Pick a passphrase of 4+ random words, fresh for every offer, and keep it: once someone accepts you have 90s to reveal it, or you forfeit. Paste only its sha256 below, never the passphrase itself:
//
//     printf %s 'your pass phrase' | shasum -a 256
//
// <gno-form path="new">
// <gno-input name="stake" type="text" placeholder="Stake in GNOT, e.g. 10 (min 1)" />
// <gno-input name="validFor" type="number" placeholder="Minutes the offer stays open, 1-60 (default 15)" />
// <gno-input name="opponent" type="text" placeholder="Opponent address (optional, leave empty for anyone)" />
// <gno-input name="commitment" type="text" placeholder="sha256 of your passphrase (64 hex chars)" />
// </gno-form>
//
// ## Ready
//
// Stake 12.5 GNOT, open to anyone for 10 min.
//
// **[Post offer for 12.5 GNOT](/r/samcrew/connect4$help&func=Offer&.send=12500000ugnot&commitment=5fb8d3d0c4f3b0fe7ab6ab1f7eda42c6e2de0c3c1d4a6c5e3d0f2f5d0e4f5a6b&opponent=&validFor=10)**
//
// -----
// # Not found
//
// -----
// # Not found

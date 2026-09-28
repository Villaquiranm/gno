// PKGPATH: gno.land/r/demo/main

package main

import (
	"chain"
	"testing"

	"gno.land/p/nt/testutils/v0"
	"gno.land/r/samcrew/connect4"
)

func main(cur realm) {
	alice := testutils.TestAddress("alice")
	bob := testutils.TestAddress("bob")
	carol := testutils.TestAddress("carol")
	c4 := chain.PackageAddress("gno.land/r/samcrew/connect4")
	stake := chain.Coins{{"ugnot", 5_000_000}}

	testing.SetRealm(testing.NewUserRealm(alice))
	testing.SetOriginSend(stake)
	testing.IssueCoins(c4, stake)
	g1 := connect4.Offer(cross(cur), "", 15)

	testing.SetRealm(testing.NewUserRealm(carol))
	testing.SetOriginSend(stake)
	testing.IssueCoins(c4, stake)
	connect4.Offer(cross(cur), bob, 5)

	testing.SetRealm(testing.NewUserRealm(bob))
	testing.SetOriginSend(stake)
	testing.IssueCoins(c4, stake)
	connect4.Accept(cross(cur), g1)

	println(connect4.Render(""))
	println("-----")
	println(connect4.Render("game/1"))
	println("-----")
	println(connect4.Render("game/x"))
}

// Output:
// # Connect 4
//
// Both players stake the same amount; the winner takes the pot minus a 0.1 GNOT fee. 90s per move: run out and you forfeit. If the first mover never moves, both stakes are refunded.
//
// [Post an offer]($help&func=Offer&opponent=&validFor=15)
//
// ## Open offers
//
// | Stake | Creator | For | Expires in | |
// |---|---|---|---|---|
// | 5 GNOT | g1vdshymmvta047h6lta047h6lta047h6l2asz94 | g1vfhkyh6lta047h6lta047h6lta047h6l03vdhu | 5m0s | [Accept]($help&func=Accept&.send=5000000ugnot&id=2) |
//
// ## Live games
//
// | Game | Pot | Turn | Time left |
// |---|---|---|---|
// | [#1](/r/samcrew/connect4:game/1) | 10 GNOT | g1v9kxjcm9ta047h6lta047h6lta047h6lzd40gh | 1m30s |
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
// 🔴 g1v9kxjcm9ta047h6lta047h6lta047h6lzd40gh
// 🟡 g1vfhkyh6lta047h6lta047h6lta047h6l03vdhu
// Stake: 5 GNOT each
//
// 1️⃣2️⃣3️⃣4️⃣5️⃣6️⃣7️⃣
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
// ⚪⚪⚪⚪⚪⚪⚪
//
// Turn: 🔴 g1v9kxjcm9ta047h6lta047h6lta047h6lzd40gh
//
// Time left: 1m30s
//
// [1]($help&func=Play&column=1&id=1) [2]($help&func=Play&column=2&id=1) [3]($help&func=Play&column=3&id=1) [4]($help&func=Play&column=4&id=1) [5]($help&func=Play&column=5&id=1) [6]($help&func=Play&column=6&id=1) [7]($help&func=Play&column=7&id=1)
//
// [Resign]($help&func=Resign&id=1)
//
// -----
// # Not found

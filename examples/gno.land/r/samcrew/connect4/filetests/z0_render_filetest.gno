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

	println(connect4.Render(""))
	println("-----")
	println(connect4.Render("game/1"))
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
// The first mover is picked by a secret passphrase you choose when posting an offer. Send its sha256 as `commitment`:
//
//     printf %s 'your pass phrase' | shasum -a 256
//
// Once your offer is accepted you have 90s to reveal the passphrase, or you forfeit. Make it unguessable (use 4+ random words): the commitment is public, so a weak passphrase can be brute-forced.
//
// [Post an offer]($help&func=Offer&.send=1000000ugnot&commitment=&opponent=&validFor=15)
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
// Waiting for creator to reveal
//
// Time left: 1m30s
//
// [Reveal]($help&func=Reveal&id=1&passphrase=) (creator only)
//
// [Resign]($help&func=Resign&id=1)
//
// -----
// # Not found
//
// -----
// # Not found

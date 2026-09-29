# Staked Connect 4 realm (`r/samcrew/connect4`)

## Context

We want a real-money Connect 4 on gno.land: equal ugnot stakes, winner takes
the pot minus a small house fee, fast turns, a lobby of offers open for the
next N minutes, and a guarantee both players are present once a game starts.

Two platform facts shape the design:

- `time.Now()` is the block timestamp. It is deterministic and only advances
  per block; nothing executes on its own, so deadlines are enforced lazily by
  whichever transaction checks them.
- Presence cannot be observed on-chain. The only proof a player is there is a
  transaction they signed.

## Decision

- **Lobby of offers.** `Offer` escrows the creator's stake with an expiry
  (1-60 min) and an optional reserved opponent. `Accept` must send exactly the
  same stake. Both use `cur.Previous().IsUserCall()` + `unsafe.OriginSend()`.
- **Presence via the move clock.** Each move has 90s of block time. `Accept`
  proves the acceptor is present, `Reveal` the creator, and the first move
  the first mover. If
  the clock runs out with zero moves the game is void and both stakes are
  refunded (the absent player may be the creator who posted long ago); after
  at least one move, the player on turn forfeits. Anyone can call
  `ClaimTimeout`, so bystanders can settle stuck games. A late `Play` panics
  rather than settling, keeping one settlement path.
- **First mover by commit-reveal.** `Offer` takes `commitment`, the lowercase
  hex `sha256(passphrase)` (e.g. `printf %s 'pass phrase' | shasum -a 256`);
  commitments need not be unique. After `Accept` the game
  waits (`Turn = 0`, `Play` rejected) for the creator to `Reveal(id,
  passphrase)` within 90s; the first mover is the low bit of
  `sha256(passphrase|id|creator|acceptor)`. The acceptor cannot know the
  outcome before accepting and cannot re-roll after. The creator knows it
  before revealing, so not revealing in time forfeits: `ClaimTimeout` then
  pays the acceptor as a win. The passphrase must be unguessable, since the
  commitment is public, and fresh per offer, since a revealed passphrase
  makes any later game using it predictable; the lobby says both. Reuse only
  hurts the reusing creator, so it is not rejected (rejecting it also let
  anyone block an offer by front-running its commitment).
- **Fee** is a flat 0.1 GNOT (100,000 ugnot) per decisive game (one with a
  winner), adjustable by the owner and snapshotted per game at `Offer`, so a
  live game's terms never change. `Offer` requires stake > fee so a winner
  always profits. Draws and void games pay no fee. Owner (`p/nt/ownable/v0`)
  can `SetFee`, `WithdrawFees`, `TransferOwnership`.
- **Cancel.** The creator may cancel an open offer at any time; anyone may
  once it has expired, so stale offers can be cleared. The refund always goes
  to the creator.
- **Non-payable calls reject coins.** Every function other than `Offer` and
  `Accept` panics if coins are sent, so nothing gets stuck in the realm.
- **Payouts are pushed** with a RealmSend banker inside the settling tx, after
  state is updated. Native coin sends run no callee code.
- **Leaderboards** (wins, total ugnot won) update only on Won/Draw, so void
  or cancelled games cannot farm stats. Render keeps a capped top-10 while
  scanning stats (Gno has no `sort.Slice`).
- An `active` tree holds only Open/Playing games so the lobby render cost does
  not grow with history.

## Alternatives considered

- First mover from a hash computed in `Accept` (block height, time, ids).
  Rejected: a tx can carry several messages and is atomic, so an acceptor can
  send `[Accept, Play]` and let the whole tx revert whenever the pick goes
  against them, then retry next block. `IsUserCall()` does not prevent this;
  it only rules out `maketx run`, not multi-message txs. Any randomness
  settled inside the accepting tx has the same problem.
- Commit-reveal by both players, or a helper app that manages secrets: fairer
  against a weak passphrase, but an extra tx or off-chain component per game.

- Matchmaking queue by stake tier: faster pairing, but no browsable lobby and
  it pairs people with idle players.
- Explicit "ready" handshake after accept: an extra tx per game that the first
  move already provides.
- Forfeit on first-move timeout: punishes an acceptor's opponent who simply
  posted an offer and left before it was taken; void is fairer.
- Chess-style time banks / per-game timeouts: more state and UI for little
  gain at this stage.

## Consequences

- Clocks are only as precise as block production: a move can land a few
  seconds past 90s if no block was produced in between. Symmetric for both
  players.
- A weak creator passphrase can be brute-forced offline from the public
  commitment, letting an acceptor take only offers where they move first.
  This cannot be enforced on-chain; the lobby warns creators.
- Each game needs an extra `Reveal` tx from the creator within 90s of accept.
- Self-play with an alt account inflates wins but costs the fee.
- Leaderboard rendering is O(players x 10); add a sorted index if the player
  count reaches tens of thousands.

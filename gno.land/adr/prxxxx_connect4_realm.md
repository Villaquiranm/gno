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
  proves the acceptor is present; the first move proves the first mover is. If
  the clock runs out with zero moves the game is void and both stakes are
  refunded (the absent player may be the creator who posted long ago); after
  at least one move, the player on turn forfeits. Anyone can call
  `ClaimTimeout`, so bystanders can settle stuck games. A late `Play` panics
  rather than settling, keeping one settlement path.
- **First mover** is the low bit of `sha256(height|time|id|creator|acceptor)`
  computed in `Accept`. The acceptor does not pick their inclusion block and
  `IsUserCall()` stops a `maketx run` script from aborting on a bad flip.
- **Fee** is a flat 0.1 GNOT (100,000 ugnot) per decided game, adjustable by
  the owner and snapshotted per game at `Offer`, so a live game's terms never
  change. `Offer` requires stake > fee so a winner always profits. Only decided games pay a fee. Owner
  (`p/nt/ownable/v0`) can `SetFee`, `WithdrawFees`, `TransferOwnership`.
- **Payouts are pushed** with a RealmSend banker inside the settling tx, after
  state is updated. Native coin sends run no callee code.
- **Leaderboards** (wins, total ugnot won) update only on Won/Draw, so void
  or cancelled games cannot farm stats. Render keeps a capped top-10 while
  scanning stats (Gno has no `sort.Slice`).
- An `active` tree holds only Open/Playing games so the lobby render cost does
  not grow with history.

## Alternatives considered

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
- A validator who is also the acceptor could bias the first-mover bit.
- Self-play with an alt account inflates wins but costs the fee.
- Leaderboard rendering is O(players x 10); add a sorted index if the player
  count reaches tens of thousands.

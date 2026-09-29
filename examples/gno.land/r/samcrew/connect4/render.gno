package connect4

import (
	"strconv"
	"strings"
	"time"

	"gno.land/p/moul/txlink/v0"
)

const gamePath = "/r/samcrew/connect4:game/"

// Render serves "" (lobby and leaderboards) and "game/<id>".
// Every string written here comes from the realm's own state (addresses,
// numbers); path input is only ever parsed as an integer.
func Render(path string) string {
	if path == "" {
		return renderLobby()
	}
	if s, ok := strings.CutPrefix(path, "game/"); ok {
		if id, err := strconv.Atoi(s); err == nil {
			if g := get(id); g != nil {
				return renderGame(g)
			}
		}
	}
	return "# Not found\n"
}

func renderLobby() string {
	now := time.Now()
	var open, live strings.Builder
	active.Iterate("", "", func(_ string, v any) bool {
		g := v.(*Game)
		switch {
		case g.Status == Open && now.Before(g.ExpiresAt):
			who := "anyone"
			if g.Opponent != "" {
				who = g.Opponent.String()
			}
			accept := txlink.NewLink("Accept").
				AddArgs("id", strconv.Itoa(g.ID)).
				SetSend(strconv.FormatInt(g.Stake, 10) + denom).URL()
			open.WriteString("| " + gnot(g.Stake) + " | " + g.Creator.String() + " | " + who +
				" | " + remaining(g.ExpiresAt.Sub(now)) + " | [Accept](" + accept + ") |\n")
		case g.Status == Playing:
			turn := g.player(g.Turn).String()
			if g.Turn == 0 {
				turn = g.Creator.String() + " (reveal)"
			}
			live.WriteString("| [#" + strconv.Itoa(g.ID) + "](" + gamePath + strconv.Itoa(g.ID) + ") | " +
				gnot(2*g.Stake) + " | " + turn + " | " + remaining(g.deadline().Sub(now)) + " |\n")
		}
		return false
	})

	var b strings.Builder
	b.WriteString("# Connect 4\n\n")
	b.WriteString("Both players stake the same amount; the winner takes the pot minus a " + gnot(fee) + " fee. " +
		"90s per move: run out and you forfeit. If the first mover never moves, both stakes are refunded.\n\n")
	b.WriteString("The first mover is picked by a secret passphrase you choose when posting an offer. " +
		"Send its sha256 as `commitment`:\n\n" +
		"    printf %s 'your pass phrase' | shasum -a 256\n\n" +
		"Once your offer is accepted you have 90s to reveal the passphrase, or you forfeit. " +
		"Make it unguessable (use 4+ random words) and use a fresh one for every offer: the commitment is public, so a weak or already-revealed passphrase lets others predict who moves first.\n\n")
	b.WriteString("[Post an offer](" + txlink.NewLink("Offer").
		AddArgs("opponent", "", "validFor", "15", "commitment", "").
		SetSend(strconv.Itoa(minStake)+denom).URL() + ")\n\n")
	b.WriteString("## Open offers\n\n")
	if open.Len() == 0 {
		b.WriteString("No open offers.\n\n")
	} else {
		b.WriteString("| Stake | Creator | For | Expires in | |\n|---|---|---|---|---|\n" + open.String() + "\n")
	}
	b.WriteString("## Live games\n\n")
	if live.Len() == 0 {
		b.WriteString("No games in progress.\n\n")
	} else {
		b.WriteString("| Game | Pot | Turn | Time left |\n|---|---|---|---|\n" + live.String() + "\n")
	}
	b.WriteString("## Most games won\n\n")
	b.WriteString(leaderboard(func(s *Stats) int64 { return int64(s.Wins) }, func(v int64) string { return strconv.FormatInt(v, 10) }))
	b.WriteString("## Most gnot won\n\n")
	b.WriteString(leaderboard(func(s *Stats) int64 { return s.GnotWon }, gnot))
	return b.String()
}

func renderGame(g *Game) string {
	var b strings.Builder
	id := strconv.Itoa(g.ID)
	b.WriteString("# Game #" + id + " — " + g.Status.String() + "\n\n")
	b.WriteString("🔴 " + g.Creator.String() + "  \n🟡 " + orNone(g.Acceptor) + "  \n")
	b.WriteString("Stake: " + gnot(g.Stake) + " each\n\n")

	b.WriteString("1️⃣2️⃣3️⃣4️⃣5️⃣6️⃣7️⃣  \n")
	for r := rows - 1; r >= 0; r-- {
		for c := 0; c < cols; c++ {
			b.WriteString([...]string{"⚪", "🔴", "🟡"}[g.Board[c][r]])
		}
		b.WriteString("  \n")
	}
	b.WriteString("\n")

	switch {
	case g.Status == Playing && g.Turn == 0:
		left := g.deadline().Sub(time.Now())
		b.WriteString("Waiting for creator to reveal\n\n")
		if left > 0 {
			b.WriteString("Time left: " + remaining(left) + "\n\n")
			b.WriteString("[Reveal](" + txlink.NewLink("Reveal").AddArgs("id", id, "passphrase", "").URL() + ") (creator only)\n\n")
		} else {
			b.WriteString("Reveal clock ran out: [Claim timeout](" +
				txlink.NewLink("ClaimTimeout").AddArgs("id", id).URL() + ")\n\n")
		}
		b.WriteString("[Resign](" + txlink.NewLink("Resign").AddArgs("id", id).URL() + ")\n")
	case g.Status == Playing:
		left := g.deadline().Sub(time.Now())
		b.WriteString("Turn: " + [...]string{"", "🔴", "🟡"}[g.Turn] + " " + g.player(g.Turn).String() + "\n\n")
		if left > 0 {
			b.WriteString("Time left: " + remaining(left) + "\n\n")
			for c := 1; c <= cols; c++ {
				b.WriteString("[" + strconv.Itoa(c) + "](" +
					txlink.NewLink("Play").AddArgs("id", id, "column", strconv.Itoa(c)).URL() + ") ")
			}
			b.WriteString("\n\n")
		} else {
			b.WriteString("Move clock ran out: [Claim timeout](" +
				txlink.NewLink("ClaimTimeout").AddArgs("id", id).URL() + ")\n\n")
		}
		b.WriteString("[Resign](" + txlink.NewLink("Resign").AddArgs("id", id).URL() + ")\n")
	case g.Status == Won:
		b.WriteString("Winner: " + g.Winner.String() + "\n")
	}
	return b.String()
}

// leaderboard renders the top 10 players by score, skipping zero scores.
func leaderboard(score func(*Stats) int64, format func(int64) string) string {
	type row struct {
		addr  string
		score int64
	}
	var top []row // kept sorted, highest first, at most 10
	stats.Iterate("", "", func(addr string, v any) bool {
		sc := score(v.(*Stats))
		if sc == 0 {
			return false
		}
		i := len(top)
		for i > 0 && top[i-1].score < sc {
			i--
		}
		if i < 10 {
			top = append(top, row{})
			copy(top[i+1:], top[i:])
			top[i] = row{addr, sc}
			if len(top) > 10 {
				top = top[:10]
			}
		}
		return false
	})
	if len(top) == 0 {
		return "No finished games yet.\n\n"
	}
	var b strings.Builder
	b.WriteString("| # | Player | |\n|---|---|---|\n")
	for i, r := range top {
		b.WriteString("| " + strconv.Itoa(i+1) + " | " + r.addr + " | " + format(r.score) + " |\n")
	}
	b.WriteString("\n")
	return b.String()
}

// gnot formats ugnot as GNOT with up to 6 decimals.
func gnot(ugnot int64) string {
	whole := strconv.FormatInt(ugnot/1_000_000, 10)
	frac := ugnot % 1_000_000
	if frac == 0 {
		return whole + " GNOT"
	}
	f := strconv.FormatInt(frac, 10)
	f = strings.TrimRight(strings.Repeat("0", 6-len(f))+f, "0")
	return whole + "." + f + " GNOT"
}

func remaining(d time.Duration) string { return d.Truncate(time.Second).String() }

func orNone(a address) string {
	if a == "" {
		return "waiting for opponent"
	}
	return a.String()
}

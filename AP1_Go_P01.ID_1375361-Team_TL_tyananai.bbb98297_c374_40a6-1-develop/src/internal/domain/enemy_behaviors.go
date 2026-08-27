package domain

import (
	"fmt"
	"math/rand/v2"

	"rogue/internal/domain/entity"
)

type EnemyBehavior interface {
	TakeTurn(g *Game, e *Enemy)
	OnAttacked(g *Game, e *Enemy) (hit bool, msg string)
	OnAttackPlayer(g *Game, e *Enemy) (forceHit bool, dmgBonus int, msg string)
	OnHitPlayer(g *Game, e *Enemy) string
}

type ZombieBehavior struct{}

func (b *ZombieBehavior) TakeTurn(g *Game, e *Enemy) {
	g.defaultMove(e, 1, false)
}
func (b *ZombieBehavior) OnAttacked(g *Game, e *Enemy) (bool, string) { return true, "" }
func (b *ZombieBehavior) OnAttackPlayer(g *Game, e *Enemy) (bool, int, string) {
	return false, 0, ""
}
func (b *ZombieBehavior) OnHitPlayer(g *Game, e *Enemy) string { return "" }

type VampireBehavior struct {
	FirstHitPending bool
}

func (b *VampireBehavior) TakeTurn(g *Game, e *Enemy) {
	g.defaultMove(e, 1, false)
}

func (b *VampireBehavior) OnAttacked(g *Game, e *Enemy) (bool, string) {
	if b.FirstHitPending {
		b.FirstHitPending = false
		return false, "Vampire dodges the first strike!"
	}
	return true, ""
}

func (b *VampireBehavior) OnAttackPlayer(g *Game, e *Enemy) (bool, int, string) {
	return false, 0, ""
}

func (b *VampireBehavior) OnHitPlayer(g *Game, e *Enemy) string {
	drain := 2 + rand.IntN(3)
	g.Player.MaxHealth = max(10, g.Player.MaxHealth-drain)
	if g.Player.Health > g.Player.MaxHealth {
		g.Player.Health = g.Player.MaxHealth
	}
	return fmt.Sprintf("Vampire drains %d MaxHP!", drain)
}

type GhostBehavior struct {
	Invisible bool
	HadCombat bool
}

func (b *GhostBehavior) TakeTurn(g *Game, e *Enemy) {
	if !b.HadCombat {
		if rand.IntN(100) < 25 {
			b.Invisible = !b.Invisible
		}
	}
	e.IsVisible = !b.Invisible

	playerPos := g.Player.Pos
	dist := manhattan(e.Pos, playerPos)

	if !e.IsAwake && dist <= e.AggroRadius() && g.HasLineOfSight(e.Pos, playerPos) {
		e.IsAwake = true
	}

	if e.IsAwake {
		if g.chasePlayer(e, 1, true) {
			return
		}
		if dist > e.AggroRadius()*2 {
			e.IsAwake = false
		}
	}

	room, ok := g.getRoomAt(e.Pos)
	if ok {
		for i := 0; i < 10; i++ {
			ry := room.StartPos.Y + rand.IntN(room.Height)
			rx := room.StartPos.X + rand.IntN(room.Width)
			np := entity.Point{Y: ry, X: rx}
			if g.walkableForEnemy(np, e) {
				e.Pos = np
				break
			}
		}
	}
}

func (b *GhostBehavior) OnAttacked(g *Game, e *Enemy) (bool, string) {
	b.HadCombat = true
	b.Invisible = false
	e.IsVisible = true
	return true, ""
}

func (b *GhostBehavior) OnAttackPlayer(g *Game, e *Enemy) (bool, int, string) {
	return false, 0, ""
}

func (b *GhostBehavior) OnHitPlayer(g *Game, e *Enemy) string {
	b.HadCombat = true
	b.Invisible = false
	e.IsVisible = true
	return ""
}

type OgreBehavior struct {
	NeedRest bool
	SureHit  bool
}

func (b *OgreBehavior) TakeTurn(g *Game, e *Enemy) {
	if b.NeedRest {
		b.NeedRest = false
		b.SureHit = true
		return
	}
	g.defaultMove(e, 2, true)
}
func (b *OgreBehavior) OnAttacked(g *Game, e *Enemy) (bool, string) { return true, "" }
func (b *OgreBehavior) OnAttackPlayer(g *Game, e *Enemy) (bool, int, string) {
	if b.SureHit {
		b.SureHit = false
		return true, 4, "heavy strike"
	}
	return false, 0, ""
}

func (b *OgreBehavior) OnHitPlayer(g *Game, e *Enemy) string {
	b.NeedRest = true
	return ""
}

type SnakeBehavior struct {
	LastDiag entity.Point
}

func (b *SnakeBehavior) TakeTurn(g *Game, e *Enemy) {
	playerPos := g.Player.Pos
	dist := manhattan(e.Pos, playerPos)

	if !e.IsAwake && dist <= e.AggroRadius() && g.HasLineOfSight(e.Pos, playerPos) {
		e.IsAwake = true
	}

	if e.IsAwake {
		if g.chasePlayer(e, 1, false) {
			return
		}
		if dist > e.AggroRadius()*2 {
			e.IsAwake = false
		}
	}
	if b.LastDiag == (entity.Point{}) {
		b.LastDiag = entity.Point{Y: 1, X: 1}
	}
	np := entity.Point{Y: e.Pos.Y + b.LastDiag.Y, X: e.Pos.X + b.LastDiag.X}
	if g.walkableForEnemy(np, e) {
		e.Pos = np
	} else {
		dirs := []entity.Point{
			{Y: -1, X: -1},
			{Y: -1, X: 1},
			{Y: 1, X: -1},
			{Y: 1, X: 1},
		}
		b.LastDiag = dirs[rand.IntN(4)]
	}
}
func (b *SnakeBehavior) OnAttacked(g *Game, e *Enemy) (bool, string) { return true, "" }
func (b *SnakeBehavior) OnAttackPlayer(g *Game, e *Enemy) (bool, int, string) {
	return false, 0, ""
}

func (b *SnakeBehavior) OnHitPlayer(g *Game, e *Enemy) string {
	if rand.IntN(100) < 15 {
		g.Player.IsPlayerSleep = true
		return "Snake puts you to sleep!"
	}
	return ""
}

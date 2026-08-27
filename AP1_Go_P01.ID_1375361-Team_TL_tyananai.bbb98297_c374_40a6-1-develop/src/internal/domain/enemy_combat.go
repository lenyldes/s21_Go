package domain

import (
	"fmt"
	"math/rand/v2"

	"rogue/internal/domain/entity"
)

const (
	baseHitPercent = 50
	dexHitBonus    = 5
)

func rollHit(attackerDex, defenderDex int) bool {
	p := baseHitPercent + (attackerDex-defenderDex)*dexHitBonus
	if p < 5 {
		p = 5
	}
	if p > 95 {
		p = 95
	}
	return rand.IntN(100) < p
}

func playerDamage(p *Player) int {
	weaponPower := 0
	if p.CurrentWeapon != nil {
		weaponPower = p.CurrentWeapon.Power
	}
	base := p.Strength + weaponPower
	variance := rand.IntN(p.Strength/2 + 2)
	return max(1, base+variance)
}

func enemyDamage(e *Enemy) int {
	base := e.Strength
	variance := rand.IntN(e.Strength/3 + 2)
	return max(1, base+variance)
}

func (g *Game) enemyAtPoint(p entity.Point) *Enemy {
	for i := range g.CurrentLevelMap.Enemies {
		e := g.CurrentLevelMap.Enemies[i]
		if e.Alive && e.Pos == p {
			return e
		}
	}
	return nil
}

func (g *Game) applyPlayerDamage(dmg int) {
	g.Player.Health -= dmg
	if g.Player.Health <= 0 {
		g.Player.Health = 0
		g.PlayerDead = true
		g.Log("YOU DIED!")
	}
}

func (g *Game) MeleeExchange(e *Enemy) {
	if g.PlayerDead || e == nil || !e.Alive {
		return
	}

	e.IsAwake = true

	hit, msg := e.Behavior.OnAttacked(g, e)
	if msg != "" {
		g.Log(msg)
	}

	if hit {
		if rollHit(g.Player.Dexterity, e.Dexterity) {
			dmg := playerDamage(g.Player)
			e.Health -= dmg
			g.Log(fmt.Sprintf("You hit %s for %d dmg.", e.Name(), dmg))
		} else {
			g.Log(fmt.Sprintf("You miss %s.", e.Name()))
		}
	}

	if e.Health <= 0 {
		e.Alive = false
		hostilityMult := int(e.Hostility) + 1
		baseLoot := hostilityMult * (e.Strength + e.Dexterity + (e.MaxHealth / 2))

		variance := 0
		if baseLoot > 0 {
			variance = rand.IntN((baseLoot / 5) + 1)
		}
		loot := baseLoot + variance

		g.Player.Backpack.Gold += loot

		g.Log(fmt.Sprintf("%s dies! +%d Gold.", e.Name(), loot))
		return
	}
}

func (g *Game) resolveEnemyHitOnPlayer(e *Enemy) {
	forceHit, dmgBonus, msg := e.Behavior.OnAttackPlayer(g, e)

	hit := forceHit || rollHit(e.Dexterity, g.Player.Dexterity)

	if !hit {
		g.Log(fmt.Sprintf("%s misses.", e.Name()))
		return
	}

	hitMsg := e.Behavior.OnHitPlayer(g, e)
	if hitMsg != "" {
		g.Log(hitMsg)
	}

	dmg := enemyDamage(e) + dmgBonus
	if msg != "" {
		g.Log(fmt.Sprintf("%s (%s) hits for %d dmg!", e.Name(), msg, dmg))
	} else {
		g.Log(fmt.Sprintf("%s hits for %d dmg.", e.Name(), dmg))
	}

	g.applyPlayerDamage(dmg)
}

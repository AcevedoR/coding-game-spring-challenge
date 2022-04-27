package main

import (
	"fmt"
	"math"
	"sort"
)

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/
var maxX = 17630
var maxY = 9000
var isBlue bool
var enemyBaseX, enemyBaseY int
var patrolState Patrol

func main() {
	// baseX: The corner of the map representing your base
	var baseX, baseY int
	fmt.Scan(&baseX, &baseY)

	isBlue = baseX < (maxY / 2)
	if isBlue {
		enemyBaseX = maxX
		enemyBaseY = maxY
	} else {
		enemyBaseX = 0
		enemyBaseY = 0
	}

	// heroesPerPlayer: Always 3
	var heroesPerPlayer int
	fmt.Scan(&heroesPerPlayer)

	attacker := 0
	defender1 := 1
	defender2 := 2
	// var attackerId, defender1Id, defender2Id int
	var turn = 1

	for {
		var health, mana, enemyHealth, enemyMana int
		fmt.Scan(&health, &mana)
		fmt.Scan(&enemyHealth, &enemyMana)

		// entityCount: Amount of heros and monsters you can see
		var entityCount int
		fmt.Scan(&entityCount)

		heroes := make([]*Hero, heroesPerPlayer, heroesPerPlayer)
		monsters := make([]Monster, 0, 20)
		// dangerousMonsters := make([]Monster, 0, 5)
		// strategy := DEFAULT
		alreadyInitHeroes := 0

		for i := 0; i < entityCount; i++ {
			// id: Unique identifier
			// type: 0=monster, 1=your hero, 2=opponent hero
			// x: Position of this entity
			// shieldLife: Ignore for this league; Count down until shield spell fades
			// isControlled: Ignore for this league; Equals 1 when this entity is under a control spell
			// health: Remaining health of this monster
			// vx: Trajectory of this monster
			// nearBase: 0=monster with no target yet, 1=monster targeting a base
			// threatFor: Given this monster's trajectory, is it a threat to 1=your base, 2=your opponent's base, 0=neither
			var id, ttype, x, y, shieldLife, isControlled, health, vx, vy, nearBase, threatFor int
			fmt.Scan(&id, &ttype, &x, &y, &shieldLife, &isControlled, &health, &vx, &vy, &nearBase, &threatFor)
			actor := Actor{id, ttype, x, y, shieldLife, isControlled, health, vx, vy, nearBase, threatFor}
			if actor.ttype == 1 {
				hero := &Hero{actor: actor}
				if alreadyInitHeroes == attacker {
					// attackerId = hero.actor.id
					hero.index = attacker
					heroes[attacker] = hero
				} else if alreadyInitHeroes == defender1 {
					hero.index = defender1
					// defender1Id = hero.actor.id
					heroes[defender1] = hero
				} else if alreadyInitHeroes == defender2 {
					hero.index = defender2
					// defender2Id = hero.actor.id
					heroes[defender2] = hero
				}
				assignDuty(hero, isBlue)
				alreadyInitHeroes++
			}
			if actor.ttype == 0 {
				distance := calculateDistance(baseX, baseY, actor.x, actor.y)

				monsters = append(monsters, Monster{actor: actor, distanceFromBase: distance})
			}
		}

		sort.Slice(monsters, func(i, j int) bool {
			return monsters[i].distanceFromBase < monsters[j].distanceFromBase
		})
		if len(monsters) > 0 {
			for id, hero := range heroes {
				if id == attacker {
					// maybeTarget := monsters[0]
					// if heroes[defender2].target != nil && heroes[defender2].target.actor.id == maybeTarget.actor.id {
					// 	if len(monsters) > 1 {
					// 		maybeTarget = monsters[1]
					// 	}
					// } else if hero.target == nil || maybeTarget.distanceFromBase < hero.target.distanceFromBase {
					// 	hero.target = &maybeTarget
					// }
					if mana > 30 && calculateDistance(hero.actor.x, hero.actor.y, enemyBaseX, enemyBaseY) < 8000 && findCloseMonsterDescending(hero.actor, monsters, 1280) != nil {
						hero.wind(isBlue)
					} else {
						hero.patrol()
					}
				} else if id == defender1 {
					if monsters[0].distanceFromBase < 5000 {
						hero.target = &monsters[0]
						hero.attack(monsters, "P")
						continue
					}
					oo, i := targetClosestMob(*hero, filterMonstersToArea(monsters, toInt(!isBlue), 1))
					hero.target = oo
					if oo != nil {
						hero.attack(monsters, fmt.Sprintf("R%d", i))
					} else {
						hero.moveToDuty()
					}
				} else if id == defender2 {
					if monsters[0].distanceFromBase < 5000 {
						if countMonstersInBaseDistance(monsters, 3000) > 1 && mana > 20 && (calculateDistanceBetweenActors(hero.actor, monsters[1].actor) < 1280 || calculateDistanceBetweenActors(hero.actor, monsters[0].actor) < 1280) {
							hero.wind(isBlue)
						} else {
							hero.target = &monsters[0]
							hero.attack(monsters, "P")
						}
						continue
					}
					oo, i := targetClosestMob(*hero, filterMonstersToArea(monsters, toInt(!isBlue), -1))
					hero.target = oo
					if hero.target != nil {
						hero.attack(monsters, fmt.Sprintf("R%d", i))
					} else {
						hero.moveToDuty()
					}
				}
			}
		} else {
			for _, hero := range heroes {
				hero.moveToDuty()
			}
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")

		// In the first league: MOVE <x> <y> | WAIT; In later leagues: | SPELL <spellParams>;

		// fmt.Fprintf(os.Stderr, "current strategy: %s", strategy)
		// fmt.Fprintf(os.Stderr, "number of heroes: %d", len(heroes))
		// fmt.Fprintf(os.Stderr, "number of monsters: %d", len(monsters))
		// fmt.Fprintf(os.Stderr, "%+v\n", heroes)

		// saviourCount := 0
		// heroCount := 0
		// for _, hero := range heroes {
		// 	if heroCount == -1 {
		// 	} else {
		// 		if hero.target != nil {
		// 			if strategy == PANIC {
		// 				monsterToMC := dangerousMonsters[0]
		// 				if saviourCount == 0 && mana >= 10 && len(dangerousMonsters) > 1 && calculateDistance(hero.actor.x, hero.actor.y, monsterToMC.actor.x, monsterToMC.actor.y) < 2200 {
		// 					fmt.Println("SPELL CONTROL ", monsterToMC.actor.id, 9200, 4400)
		// 					saviourCount++
		// 					continue
		// 				}
		// 			}
		// 			point := getOptimalPoint(*hero.target, monsters)
		// 			fmt.Println("MOVE ", point.x, " ", point.y)
		// 		} else {
		// 			if calculateDistance(hero.actor.x, hero.actor.y, hero.dutyY, hero.dutyY) < 300 {
		// 				if mana > 20 && turn > 5 {
		// 					if tryToShieldSomeone(heroes) {
		// 						continue
		// 					}
		// 				}
		// 			}
		// 			fmt.Println("MOVE ", hero.dutyX, hero.dutyY)
		// 		}
		// 	}
		// 	heroCount++
		// }
		turn++
	}
}

type Actor struct {
	id, ttype, x, y, shieldLife, isControlled, health, vx, vy, nearBase, threatFor int
}
type Monster struct {
	actor            Actor
	distanceFromBase int
}
type Hero struct {
	index          int
	actor          Actor
	target         *Monster
	fallbackTarget *Monster
	dutyX          int
	dutyY          int
}
type Patrol struct {
	coeff       int
	isAscending bool
}

func getAlreadyTargerMonsters(heroes map[int]*Hero) []int {
	list := make([]int, 0, 10)
	for _, value := range heroes {
		if value.target != nil {
			list = append(list, value.target.actor.id)
		}
	}
	return list
}
func contains(list []int, x int) bool {
	for i := 0; i < len(list); i++ {
		if x == list[i] {
			return true
		}
	}
	return false
}
func assignDuty(hero *Hero, isBlue bool) {
	if isBlue {
		if hero.index == 0 {
			hero.dutyX = enemyBaseX - 5000
			hero.dutyY = enemyBaseY - 5000
		} else if hero.index == 1 {
			hero.dutyX = 6000
			hero.dutyY = 2500
		} else {
			hero.dutyX = 3500
			hero.dutyY = 6000
		}
	} else {
		if hero.index == 0 {
			hero.dutyX = enemyBaseX + 5000
			hero.dutyY = enemyBaseY + 5000
		} else if hero.index == 1 {
			hero.dutyX = 14300
			hero.dutyY = 4400
		} else {
			hero.dutyX = 11800
			hero.dutyY = 7000
		}
	}
}
func calculateDistanceBetweenActors(actor1 Actor, actor2 Actor) int {
	return calculateDistance(actor1.x, actor1.y, actor2.x, actor2.y)
}
func calculateDistance(p1X int, p1Y int, p2X int, p2Y int) int {
	first := PowInt(p2X-p1X, 2)
	second := PowInt(p2Y-p1Y, 2)
	return int(math.Round(math.Sqrt(float64(first) + float64(second))))
}
func PowInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
func tryToShieldSomeone(heroes []*Hero) bool {
	for id, ally := range heroes {
		if ally.actor.shieldLife < 1 {
			fmt.Println("SPELL SHIELD ", id)
			return true
		}
	}
	return false
}

const (
	DEFAULT = "DEFAULT"
	PANIC   = "PANIC"
)

func isInArea(x int, y int, ax int, ay int, bx int, by int) bool {
	// a = upper left
	if x > ax && x < bx && y > ay && y < by {
		return true
	}
	return false
}
func isInDefensiveArea(x int, y int, isBlue bool) bool {
	var ax, ay, bx, by int
	if !isBlue {
		ax = maxX - 7500
		ay = maxY - 7500
		bx = maxX - 5000
		by = maxY
	} else {
		ax = 5000
		ay = 0
		bx = 7500
		by = 7500
	}
	if isInArea(x, y, ax, ay, bx, by) {
		return true
	}
	if !isBlue {
		ax = maxX - 7500
		ay = maxY - 7500
		bx = maxX
		by = maxY - 5000
	} else {
		ax = 0
		ay = 5000
		bx = 7500
		by = 7500
	}
	if isInArea(x, y, ax, ay, bx, by) {
		return true
	}
	return false
}

type Point struct {
	x int
	y int
}

func getOptimalPoint(target Monster, monsters []Monster) Point {
	for _, m := range monsters {
		distance := calculateDistance(target.actor.x, target.actor.y, m.actor.x, m.actor.y)
		if target.actor.id != m.actor.id && distance < 800 {
			return Point{x: (target.actor.x + m.actor.x) / 2, y: (target.actor.y + m.actor.y) / 2}
		}
	}
	return Point{target.actor.x, target.actor.y}
}
func Abs(x int) int {
	if x < 0 {
		return x * -1
	}
	return x
}
func filterMonstersToArea(monsters []Monster, vectorx int, vectory int) []Monster {
	result := make([]Monster, 0, len(monsters))
	for _, m := range monsters {
		cX := false
		cY := false
		if vectorx == 1 {
			// right
			cX = m.actor.x > maxX/2
		} else if vectorx == -1 {
			cX = m.actor.x < maxX/2
		}
		if vectory == 1 {
			// up
			cY = m.actor.y < maxY/2
		} else if vectory == -1 {
			cY = m.actor.y > maxY/2
		}
		if cX && cY {
			result = append(result, m)
		}
	}
	return result
}
func targetClosestMob(hero Hero, monsters []Monster) (*Monster, int) {
	var target *Monster = nil
	var targetDistance int
	var targetRank int
	count := 0
	for i, m := range monsters {
		if target == nil {
			target = &m
			targetRank = i
			targetDistance = calculateDistanceBetweenActors(hero.actor, m.actor)
		} else {
			if m.distanceFromBase > 2000 && calculateDistanceBetweenActors(hero.actor, m.actor) < targetDistance {
				target = &m
				targetRank = i
			}
		}
		if count > -1 { // TODO fix disabled
			break
		}
		count++
	}
	return target, targetRank
}
func toInt(b bool) int {
	if b {
		return 1
	} else {
		return -1
	}
}
func (hero Hero) attack(monsters []Monster, message string) {
	point := getOptimalPoint(*hero.target, monsters)
	fmt.Println("MOVE ", point.x, " ", point.y, hero.index, ": ", message)
}
func (hero Hero) moveToDuty() {
	fmt.Println("MOVE ", hero.dutyX, hero.dutyY, hero.index, ": D")
}
func (hero Hero) patrol() {
	if patrolState.isAscending && patrolState.coeff > 4 {
		patrolState.isAscending = false
	} else if !patrolState.isAscending && patrolState.coeff < -4 {
		patrolState.isAscending = true
	}
	if patrolState.isAscending {
		patrolState.coeff++
	} else {
		patrolState.coeff--
	}
	var xcoeff int
	if isBlue {
		xcoeff = Abs(patrolState.coeff)
	} else {
		xcoeff = -Abs(patrolState.coeff)
	}
	fmt.Println("MOVE ", hero.dutyX+xcoeff*800, hero.dutyY+patrolState.coeff*700, hero.index, "p", patrolState.coeff)
}
func (hero Hero) wind(isBlue bool) {
	fmt.Println("SPELL WIND ", enemyBaseX, enemyBaseY, hero.index)
}
func countMonstersInBaseDistance(monsters []Monster, dist int) int {
	var count int
	for _, m := range monsters {
		if m.distanceFromBase < dist {
			count++
		}
	}
	return count
}
func findCloseMonsterDescending(actor Actor, monsters []Monster, distance int) *Monster {
	for i := len(monsters) - 1; i >= 0; i-- {
		m := monsters[i]
		if calculateDistanceBetweenActors(actor, m.actor) < distance {
			return &m
		}
	}
	return nil
}

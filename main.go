package main

import (
	"fmt"
	"os"
)

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
    // baseX: The corner of the map representing your base
    var baseX, baseY int
    fmt.Scan(&baseX, &baseY)
    
    // heroesPerPlayer: Always 3
    var heroesPerPlayer int
    fmt.Scan(&heroesPerPlayer)



    for {
        for i := 0; i < 2; i++ {
            // health: Each player's base health
            // mana: Ignore in the first league; Spend ten mana to cast a spell
            var health, mana int
            fmt.Scan(&health, &mana)
        }
        // entityCount: Amount of heros and monsters you can see
        var entityCount int
        fmt.Scan(&entityCount)
        
		heroes := make(map[int]*Hero)
		monsters := make([]Monster, 0, 20)

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
				if val, ok := heroes[actor.id]; ok {
					val.actor = actor
				} else {
					heroes[actor.id] = &Hero{actor: actor}
				}
			}
			if actor.ttype == 0  {
				threat := 0
				if actor.threatFor == 1 {
					threat += 10
					if actor.nearBase == 1 {
						threat += 5
					}
				}
				monsters = append(monsters, Monster{actor, threat})
			}
        }

        for i := 0; i < len(monsters); i++ {
			monster := monsters[i]
			for _, hero := range heroes {
				if hero.target == nil || monster.threat > hero.target.threat {
					
					// fmt.Fprintf(os.Stderr, "hero: %d targeting: %d", hero.actor.id, hero.target.actor.id)
					hero.target = &monster
				}
			}
		}
            
		// fmt.Fprintln(os.Stderr, "Debug messages...")
		
		// In the first league: MOVE <x> <y> | WAIT; In later leagues: | SPELL <spellParams>;

		fmt.Fprintf(os.Stderr, "number of heroes: %d", len(heroes))
		fmt.Fprintf(os.Stderr, "number of monsters: %d", len(monsters))
		fmt.Fprintf(os.Stderr, "%+v\n", heroes)

		for _, hero := range heroes {
			if (hero.target != nil) {
				fmt.Println("MOVE ", hero.target.actor.x, " ", hero.target.actor.y)
			} else {
				fmt.Println("WAIT")
			}
		}
	
    }
}
type Actor struct {
	id, ttype, x, y, shieldLife, isControlled, health, vx, vy, nearBase, threatFor int
}
type Monster struct {
	actor Actor
	threat int
}
type Hero struct {
	actor Actor
	target *Monster
}
func getAlreadyTargerMonsters(heroes map[int]*Hero) []int{
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
package main

import "fmt"

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
        
		heroes := make([]Actor, 0, 20)
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
				heroes = append(heroes, actor)
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
		var target1, target2, target3 Monster
        for i := 0; i < len(monsters); i++ {
			m := monsters[i]
			if m.threat > target1.threat {
				target3 = target2
				target2 = target1
				target1 = m
			} else if m.threat > target2.threat {
				target3 = target2
				target2 = m
			} else if m.threat > target3.threat {
				target3 = m
			}
		}
		targets := []Monster{target1, target2, target3}

            
		// fmt.Fprintln(os.Stderr, "Debug messages...")
		
		// In the first league: MOVE <x> <y> | WAIT; In later leagues: | SPELL <spellParams>;

		activeHeroes := 0
		for i := 0; i < len(targets); i++ {
			t := targets[i]
			if t.threat != 0 {
				fmt.Println("MOVE ", t.actor.x, " ", t.actor.y)
				activeHeroes++ 
			}
		}
		inactiveHeroes := heroesPerPlayer - activeHeroes 
		for i := 0; i < inactiveHeroes; i++ {
			fmt.Println("WAIT")
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
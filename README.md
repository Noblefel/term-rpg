A simple terminal-based RPG where you can battle against various enemies.

<img src="https://github.com/Noblefel/term-rpg/blob/main/sample.PNG">

# How to play

#### clone the repository
```bash
git clone https://github.com/Noblefel/term-rpg
``` 

#### install dependencies 
```sh
go mod tidy
```

#### run the app 
```sh
go run main.go
```

# About the game

### Actions
| Name | Effect |
| -------- | ------- |  
| ⚔️ Attack | Deals dmg based on attack attribute + random value |
| 🛡️ Guard | Boost dmg reduction by 20% for 2 turns |
| 🔥 Fury | Sacrifice hp for +5 attack point for 2 turns (player only) |
| 🏃 Flee | Escape from the battle (player only) |

### Player Perks
| Name | Effect |
| -------- | ------- |  
| 💰 Greed | Gain 15% more loot |
| 🛡️ Resiliency | +1 defense point and 10% dmg reduction |
| ⚔️ Havoc | +25% attack bonus, but -15 HP cap|
| ⌛ Temporal | +1 extra turn for bonus effects |

### Enemies
| Name | Feat/Special |
| -------- | ------- |  
| Acolyte 🧙| Has damage reduction |
| Assassin 🗡️| Good in offense with high attack attribute |
| Changeling 🎭 | Will mimic player's attributes |
| Evil Genie 🧞 | Can straight up **curse** (debuffs) your attributes |
| Golem 🗿  | High defense, massive damage, but more likely to skip their own turn |
| Snakes 🐍 | Though low hp, they could deal high damage |
| Spike Turtle 🐢 | Reflect some of the original damage back to the attacker |
| Thug 🥊 | A good all-rounder with decent attributes |
| Vampire 🧛 | Heals every attack and drains 5% current hp as extra damage |
| Wraith 👻 | Absorbs fixed number of hp, ignoring any defense and effects |

### Other
- **Rest**: recover (5 + 10% of player's hp cap + 0-8) of hp
- **Train**: 30% chance to buff random attributes
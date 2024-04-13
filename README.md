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
| 🛡️ Defend | Boost dmg reduction by 20% for 5 seconds (cannot stack) |
| 🔥 Fury | Sacrifice hp for +4 attack point for 5 seconds (player only) |
| 🏃 Flee | Escape from the battle (player only) |

### Player Perks
| Name | Effect |
| -------- | ------- |  
| 💰 Greed | Gain 15% more loot |
| 🛡️ Resiliency | +1 defense point and 10% dmg reduction |
| ⚔️ Havoc | +25% Attack, but -15 HP cap|
| ⏰ Temporal | +8 seconds to actions bonus modifier |

### Enemies
| Name | Feat |
| -------- | ------- |  
| Acolyte 🧙| Has damage reduction |
| Assassin 🗡️| Deals high damage |
| Golem 🗿  | High defense. 30% chance of dealing massive damage, otherwise 0 (miss) |
| Snakes 🐍 | Though low hp, they could deal high damage. Drops no loot |
| Thug 🥊 | A good all-rounder with decent attributes |

### Other
- **Rest**: recover (5 + 10% of player's hp cap + 0-8) of hp
- **Train**: 20% chance to buff random attributes
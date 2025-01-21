A simple turn-based RPG game in terminal.

https://github.com/user-attachments/assets/74143399-d9da-4b58-b220-ffb93c6d9c41

### Actions

| Name      | Effect                                             |
| --------- | -------------------------------------------------- |
| ⚔️ Attack | deal dmg based on strength                         |
| 🔥 Skills | use variety of skills, cost energy                 |
| 🧰 Items  | to be added                                        |
| 🏃 Flee   | try to escape from the battle, may or may not fail |
| ⌛ Skip   | do nothing                                         |

### Player perks

| Name          | Effect                                       |
| ------------- | -------------------------------------------- |
| 🛡️ Resiliency | increase survivability                       |
| ⚔️ Havoc      | extra damage, but low starting gold & max hp |
| 🐻 Berserk    | more powerful the lower your hp is           |
| 🐇 Ingenious  | skill cooldown is reduced by 1               |
| 🍹 Poisoner   | give poison effect at the start of battle    |

### Player skills

| Name          | Desc                                                            |
| ------------- | --------------------------------------------------------------- |
| charge        | attack with 130% strength                                       |
| heal          | recover hp by atleast 8% of your hpcap                          |
| frenzy        | sacrifice hp to attack with 250% strength (no inherited effect) |
| vision        | see enemy attributes (no cost)                                  |
| drain         | take 20% of enemy current hp                                    |
| absorb        | take 8% of enemy hp cap and ignore defense                      |
| trick         | make the enemy target themselves                                |
| poison        | attack 60% strength and poison enemy for 3 turns                |
| stun          | attack 30% strength and stun enemy for 2 turns                  |
| fireball      | deal moderate amount of damage (fixed number/rng)               |
| meteor strike | deal huge amount of damage (fixed number)                       |

### Enemies

| Name          | Summary                                                | Special              |
| ------------- | ------------------------------------------------------ | -------------------- |
| Knight 🛡️     | decent all-rounder with good defense                   | defense buff         |
| Wizard 🧙     | easy to take out, but has great damage                 | various spells, heal |
| Changeling 🎭 | will morph itself to be like **you**                   | attribute copy       |
| Vampire 🧛    | powerful enemy with both good attack and survivability | lifesteal            |
| Demon 👹      | powerful enemy that cares little about your defense    | ignore defense       |

### Other

- **Battle**: go to battle, the enemy difficulty is scaled
- **Deep forest**: explore the deep forest, may get something
- **Dungeon**: explore the dungeon, TO BE ADDED
- **Shop**: TO BE ADDED
- **Rest**: recover your healthpoint
- **Train**: 40% chance to buff random attributes

```bash
git clone https://github.com/Noblefel/term-rpg
```

```sh
go mod tidy
```

```sh
go run .
```

A simple turn-based RPG game in terminal.

https://github.com/user-attachments/assets/74143399-d9da-4b58-b220-ffb93c6d9c41

### Actions

| Name      | Effect                                             |
| --------- | -------------------------------------------------- |
| ⚔️ Attack | deal dmg based on strength + random value          |
| 🔥 Skills | to be added                                        |
| 🧰 Items  | to be added                                        |
| 🏃 Flee   | try to escape from the battle, may or may not fail |
| ⌛ Skip   | do nothing                                         |

### Player perks

| Name          | Effect                                       |
| ------------- | -------------------------------------------- |
| 🛡️ Resiliency | increase survivability                       |
| ⚔️ Havoc      | extra damage, but low starting gold & max hp |
| 🐻 Berserk    | more powerful the lower your hp is           |

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

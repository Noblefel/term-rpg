### Player perks

| Name         | Effect                                                        |
| ------------ | ------------------------------------------------------------- |
| 🛡️ Resilient | increase overall defense                                      |
| ⚔️ Havoc     | +20% damage, but low starting hp & energy cap                 |
| 🐻 Berserk   | more powerful the lower your hp is                            |
| 🐇 Ingenious | +2 energy cap, skill cooldown is reduced by 2                 |
| 🍹 Poisoner  | give severe poisoning effect on every battle                  |
| ⚰️ Deadman   | give weaken effect on every battle                            |
| 🏃 Survivor  | almost always succeed when fleeing                            |
| 🌙 Coven     | (tba) unlock special skills, but -10% damage & extra cooldown |
| 🦅 Falconry  | (tba) get a companion that assist you in battle               |

For **deadman**, effect will be on the enemy for 3 turns. Weaken means they do 13% less damage and reducing their defense by half. Also you will get buffed when facing Undeads.

For **coven** (to be added), each attack will give 0-3 stack of "hex" which can be used to cast special skills. You will get penalty when facing against Clerics.

### Enemies

| Name            | Summary                                                    | Special                            |
| --------------- | ---------------------------------------------------------- | ---------------------------------- |
| Knight 🛡️       | decent all-rounder with good defense                       | defense buff, strengthen           |
| Wizard 🧙       | lacks in traditional defense, makes it up with spells      | heal, immunity, barrier, confuse   |
| Changeling 🎭   | will morph itself to be like you, atleast tried            | attribute copy x 0.75              |
| Vampire 🧛      | powerful enemy with both good attack and survivability     | lifesteal, poison                  |
| Demon 👹        | powerful enemy that cares little about your defense        | ignore defense, burning            |
| Shardling ⛰️    | has a tough defense and will reflect damage                | reflect damage                     |
| Evil Genie 🔮   | can straight up **curse** you permanently                  | debuff, illusion, force-field      |
| Celestial ☄️    | a being from beyond this world with tons of hp             | prevent skills, heal aura, burning |
| Shapeshift 🎭   | a straight-up better version of changeling                 | attribute copy x 1.5               |
| Undead 🧟       | the mere sight of this enemy weakens you                   | weaken, poison                     |
| Scorpion 🦂     | (tba) powerful enemy with high damage and is very venomous | severe poison                      |
| Inferno 🔥      | (tba) this creature was born from the strongest of fires   | severe burning                     |
| Vine Monster 🌲 | (tba) attacks have high chance to ensnare you              | stun, heal aura                    |
| Werewolf 🐺     | (tba) starts off weak until it gets **awakened**           | extreme buff                       |
| Cleric ☀️       | (tba) may convert "purify" stack into **MASSIVE** dmg      | **final exorcism**                 |
| Artificer 🛠️    | (tba) advanced equipments gave them **powerful** abilities | swap hp, lay traps, force-field    |

### Player skills (can equip 5 at a time)

| Name          | Desc                                                                  |
| ------------- | --------------------------------------------------------------------- |
| charge        | attack with 130% strength                                             |
| frenzy        | sacrifice hp to attack with 250% strength                             |
| great blow    | sacrifice the next turn to attack with 210% strength                  |
| poison        | attack 85% strength and poison enemy for 3 turns                      |
| stun          | attack 60% strength and stun enemy for 2 turns                        |
| swift strike  | attack 85% strength (doesnt consume turn)                             |
| knives throw  | attack 15 fixed damage (doesnt consume turn, no cd)                   |
| fireball      | deal fixed amount of damage and give burning effect                   |
| meteor strike | deal huge amount of damage (fixed number rng)                         |
| strengthen    | attack 100% strength to increase damage by 10% for 3 turns            |
| barrier       | reduce incoming damage by 40% for 2 turns                             |
| force-field   | reduce incoming damage by 15% for 5 turns                             |
| heal spell    | recover hp by atleast 12% of hpcap                                    |
| heal aura     | recover hp by atleast 7% of hpcap for 3 turns                         |
| heal potion   | recover hp by 34 (fixed number)                                       |
| drain         | take 22% of enemy current hp                                          |
| absorb        | take 10% of enemy hp cap and ignore defense                           |
| trick         | make the enemy self-target                                            |
| vision        | see enemy attributes (no cost, doesnt consume turn)                   |
| hex chant     | (coven) sacrifice hp to give 5 stack of hex (doesnt consume turn)     |
| hex cleanse   | (coven) 5 hex to remove all effects & cd on you (doesnt consume turn) |
| hex barrier   | (coven) 8 hex to grant 65% damage reduction for 3 turns               |
| hex curse     | (coven) 6 hex to give poison, weaken, and severe burning for 3 turns  |
| blood ritual  | (coven) 1-12 hex to deal from small to **MASSIVE** damage             |

### Other

- **Battle**: go to battle, the enemy difficulty is scaled
- **Deep forest**: explore the deep forest, may get something
- **Dungeon**: explore the dungeon, TO BE ADDED
- **Shop**: TO BE ADDED
- **Guest house**: recover your healthpoint
- **Training grounds**: 40% chance to buff random attributes
- **Switch perk**: allow you to change perk.

tips: for quick gameplay pick poisoner then spam stun and poison skill.

```bash
git clone https://github.com/Noblefel/term-rpg
```

```sh
go mod tidy
```

```sh
go run .
```

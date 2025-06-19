### Player perks

| Name         | Effect                                                        |
| ------------ | ------------------------------------------------------------- |
| 🛡️ Resilient | increase overall defense                                      |
| ⚔️ Havoc     | +20% damage, but low starting hp & energy cap                 |
| 🐻 Berserk   | more powerful the lower your hp is                            |
| 🐇 Ingenious | +2 energy cap, skill cooldown is reduced by 2                 |
| 🍹 Poisoner  | inflict poisoning effect at the start of battle               |
| ⚰️ Deadman   | inflict weaken effect at the start of battle                  |
| 🏃 Survivor  | almost always succeed when fleeing                            |
| 🎃 Insanity  | it can either go really well or really bad                    |
| 🌙 Coven     | (tba) unlock special skills, but -10% damage & extra cooldown |
| 🦅 Falconry  | (tba) get a companion that assist you in battle               |

For **resilient**, you get 10% damage reduction plus small extra hp cap and defense, also incoming damage cannot exceed 16% of your hp cap.

For **berserk**, bonuses starts to apply when your hp gets below 40%.

For **poisoner**, inflicted effect will stay for 5 turns.

For **deadman**, inflicted effect will stay for 3 turns. Weaken means they do 13% less damage and reducing their defense by half. Also you will get buffed when facing Undeads.

For **survivor**, you also get small extra agility. Agility is used to trigger crit, dodge, and to flee the battle.

For **insanity**, you may or may not: deal extra/less damage, skip on your own turn, get random buff/debuff to your attributes, randomized skill cooldown and 1% chance to instantly slain the enemy each turn.

For **coven** (to be added), each attack will inflict 0-3 stack of "hex" which can be used to cast special skills. You will get penalty when facing against Clerics.

### Enemies

| Name            | Summary                                                  | Special                            |
| --------------- | -------------------------------------------------------- | ---------------------------------- |
| Knight 🛡️       | decent all-rounder with good defense                     | defense buff, strengthen           |
| Wizard 🧙       | lacks in traditional defense, makes it up with spells    | heal, immunity, barrier, confuse   |
| Changeling 🎭   | will morph itself to be like you, atleast tried          | attribute copy x 0.75              |
| Vampire 🧛      | powerful enemy that recovers hp every attack             | lifesteal, poison                  |
| Demon 👹        | powerful enemy that cares little about your defense      | ignore defense, burning            |
| Shardling ⛰️    | has a tough defense and will reflect damage              | reflect                            |
| Evil Genie 🔮   | can straight up **curse** you permanently                | debuff, illusion, force-field      |
| Celestial ☄️    | a being from beyond this world with tons of hp           | prevent skills, heal aura, burning |
| Shapeshift 🎭   | a better, more deadlier version of changeling            | attribute copy x 1.25              |
| Undead 🧟       | the mere sight of this enemy weakens you                 | weaken, poison                     |
| Scorpion 🦂     | powerful enemy with high damage and venom                | severe poison, ignore 30% defense  |
| Goblin 🗡️       | weak defense but makes it up with agility                | confuse                            |
| Infernal 🔥     | (tba) this creature was born from the strongest of fires | severe burning                     |
| Vine Monster 🌲 | (tba) attacks have high chance to ensnare you            | stun, reflect (small), heal aura   |
| Werewolf 🐺     | (tba) starts off weak until it gets **awakened**         | extreme buff, bleeding             |
| Cleric ☀️       | (tba) may convert "purify" stack into **MASSIVE** dmg    | **final exorcism**                 |
| Artificer 🛠️    | (tba) advanced equipments gave them **unique** abilities | swap hp, lay traps, force-field    |

### Player skills (can equip 5 at a time)

| Name          | Desc                                                                  |
| ------------- | --------------------------------------------------------------------- |
| charge        | attack with 130% strength                                             |
| frenzy        | sacrifice hp to attack with 250% strength                             |
| great blow    | sacrifice the next turn to attack with 210% strength                  |
| poison        | attack 85% strength and poison enemy for 3 turns                      |
| stun          | attack 60% strength and stun enemy for 2 turns                        |
| swift strike  | attack 85% strength (doesnt consume turn)                             |
| knives throw  | attack 40 fixed damage (doesnt consume turn, no cd)                   |
| fireball      | deal moderate amount of damage and inflict burning                    |
| meteor strike | deal huge amount of damage (fixed number rng)                         |
| strengthen    | attack 100% strength to increase damage by 10% for 3 turns            |
| barrier       | reduce incoming damage by 40% for 2 turns                             |
| force-field   | reduce incoming damage by 15% for 5 turns                             |
| heal spell    | recover hp by atleast 15% of hpcap                                    |
| heal aura     | recover hp by atleast 7% of hpcap for 3 turns                         |
| heal potion   | recover hp by 40 (fixed number)                                       |
| drain         | take 22% of enemy current hp                                          |
| absorb        | take 10% of enemy hp cap, ignore defense and effects                  |
| trick         | make the enemy self-target                                            |
| vision        | see enemy attributes (no cost, no cd, doesnt consume turn)            |
| hex chant     | (coven) sacrifice hp to inflict 5 stacks of hex (doesnt consume turn) |
| hex cleanse   | (coven) 5 hex to remove all effects & cd on you (doesnt consume turn) |
| hex barrier   | (coven) 8 hex to grant 65% damage reduction for 3 turns               |
| hex curse     | (coven) 6 hex to inflict poison, weaken, and burning for 3 turns      |
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

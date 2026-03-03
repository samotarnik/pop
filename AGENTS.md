# Prince of Prekmurje

`Prince of Prekmurje` is a desktop 2D arcade game built in Go programming language with the 
[Ebiten](https://github.com/hajimehoshi/ebiten) game engine. It's sort of a parody of the original 
[Prince of Persia](https://en.wikipedia.org/wiki/Prince_of_Persia_(1989_video_game)) from 1989.

The Prince of Prekmurje is a farmer-like character moving through the fields of Premurje, avoiding obstacles.

## Gameplay

- The Prince can run, jump, fight, duck.
- Carries a pitchfork to fight, but it can only be used when the Prince is standing still.
- Beer cans are life units, start out with 3 life, max out at 5 life.
- Lose all life and you die, respawn at the start of the level.
- Pick up beer cans to gain life.
- Jump over cats, touch them and the Prince loses all life on the spot.
- Duck flying hay bales, get hit and you lose 1 life.
- Fight other enemy farmers using pitchforks, every successful stab is worth 1 life, get hit and you lose 1 life. 
- Typical enemy farmer has 1 life, final boss has 3 life.

## Other features in v1

- 3-level mini campaign
- Lives system with respawn
- Intro, inter-level, and outro story cards
- 8-bit background music

## Controls

- `J/L` or `Left/Right`: move
- `I/Up`: jump
- `K/Down`: duck
- `Space`: use pitchfork
- `R`: restart current level
- `Esc`: pause
- `M`: mute toggle placeholder
- `Enter`: continue in menu/story/game-over/victory scenes

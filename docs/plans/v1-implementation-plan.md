# Prince of Prekmurje v1 Implementation Plan (From Empty Repo)

## Summary
Build a full v1 playable desktop game in Go + Ebiten with:
- All AGENTS.md mechanics and scenes.
- Placeholder art/audio.
- Simple struct-based runtime architecture (no ECS).
- JSON-authored levels.
- macOS-first build/release.
- Core unit + integration test coverage.
- Embedded assets via `go:embed`.
- English story cards.
- Separate lives system: 3 lives total per run.

## Scope and Success Criteria
1. Deliver a playable 3-level campaign with intro, inter-level, outro, pause, game-over, and victory flows.
2. Implement controls exactly as specified (`J/L`, arrows, `I/Up`, `K/Down`, `Space`, `R`, `Esc`, `M`, `Enter`).
3. Enforce combat/hazard rules exactly:
- Cat contact: instant loss of all current health.
- Hay bale hit: minus 1 health.
- Enemy stab received: minus 1 health.
- Player stab landed: enemy loses 1 health.
4. Health starts at 3, max 5. Death triggers respawn at level start and consumes 1 life. Game over when lives reach 0.
5. Final boss has 3 initial health; normal enemy has 1.
6. Runs on macOS locally with `go run` and packaged binary.

## Architecture
1. Project layout:
- `cmd/pop/main.go` bootstraps window, game loop.
- `internal/game` core loop, scene manager, shared context.
- `internal/scene` menu/story/level/pause/gameover/victory scenes.
- `internal/entity` player, enemies, projectiles/hazards, pickups.
- `internal/physics` collision, hitboxes, world bounds.
- `internal/level` JSON loading/validation/spawn building.
- `internal/audio` BGM + SFX + mute toggle placeholder wiring.
- `assets` placeholder spritesheets, audio, level JSON, story JSON.
- `test` integration scenario tests.
2. Runtime model:
- Fixed 60 TPS update.
- Scene-based state machine.
- Per-scene `Update()`, `Draw()`, `OnEnter()`, `OnExit()`.
- Entities as structs implementing small interfaces (`Updateable`, `Drawable`, `Collider`, `Damageable`).

## Public APIs / Interfaces / Types
1. `type Scene interface { OnEnter(*Context); OnExit(*Context); Update(*Context) error; Draw(screen *ebiten.Image) }`
2. `type SceneID string` constants for `Intro`, `StoryCard`, `Level`, `Pause`, `GameOver`, `Victory`.
3. `type GameState struct { CurrentLevel int; Health int; MaxHealth int; Lives int; Score int; Muted bool }`
4. `type Player struct` with explicit state enum: `Idle`, `Run`, `Jump`, `Duck`, `Attack`, `Hit`, `Dead`.
5. `type Enemy struct` with AI mode enum: `Patrol`, `Engage`, `Attack`, `Dead`.
6. `type LevelData struct` decoded from JSON with:
- `id`, `name`, `width`, `height`, `tileSize`
- `tiles` (2D int array)
- `spawn` (`player`, `enemies`, `hazards`, `pickups`, `goal`)
- `bgm` and optional `storyCardId`
7. `type CombatEvent struct { SourceID string; TargetID string; Damage int; Heal int; InstantKill bool }`
8. Config constants package for tuning values (movement/combat timing/damage) to avoid hardcoded literals.

## Gameplay Implementation Decisions
1. Display and timing:
- Internal resolution `640x480`, scaled to window `1280x960`.
- `SetTPS(60)` with deterministic update logic.
2. Player mechanics:
- Run speed, jump velocity, gravity, duck hitbox reduction, and attack active window in config constants.
- Attack only allowed while stationary and grounded.
3. Damage/heal:
- Health clamped `[0,5]`.
- Non-lethal damage applies short invulnerability window.
- Cat damage bypasses invulnerability and sets health to 0.
4. Enemy behavior:
- Normal farmer: approach + attack when in range; 1 HP.
- Boss: same base model with tuned range/timing; 3 HP.
5. Level completion:
- Reach goal trigger with player alive.
- Preserve health/lives between levels.
6. Scene flow:
- Start -> Intro -> Level1 -> Story -> Level2 -> Story -> Level3(Boss) -> Outro -> Victory.
- Any zero-lives state -> Game Over -> Enter returns to Intro with fresh run.

## Data and Content
1. JSON schemas:
- `levels/level1.json`, `level2.json`, `level3.json`.
- `story/cards.json` keyed by card ID.
2. Asset strategy:
- Minimal placeholder sprites by entity type and state.
- Single looped 8-bit BGM track plus simple SFX stubs.
3. Asset packaging:
- Use `go:embed` for all assets.
- Optional dev override is out of scope for v1.

## Testing and Acceptance
1. Unit tests:
- Health/lives transitions and clamp behavior.
- Damage rules per hazard/combat type.
- Attack gating (must be standing still).
- Enemy HP and boss HP behavior.
- JSON parsing/validation failures and success paths.
2. Integration tests:
- Scene transition progression across 3-level campaign.
- Respawn at level start after death.
- Game over when lives hit 0.
- Victory after boss defeat.
3. Manual acceptance checklist:
- All controls respond as specified.
- Pause/restart/mute placeholder actions work.
- Story cards advance on Enter.
- No soft-lock on death, pause, or scene transitions.

## Delivery Phases
1. Phase 1: Bootstrap
- Initialize Go module, Ebiten entrypoint, scene manager, config, embed pipeline.
2. Phase 2: Core Gameplay Systems
- Player controller, physics/collision, hazards, pickups, combat core.
3. Phase 3: Level and Enemy Systems
- JSON loader, level spawning, enemy AI, boss tuning.
4. Phase 4: Campaign Flow
- Story card scenes, inter-level transitions, game-over/victory.
5. Phase 5: Audio/UI/Polish
- BGM/SFX placeholders, HUD (health/lives/level), pause/menu polish.
6. Phase 6: QA and Packaging
- Tests, bug fixes, macOS build artifact and run instructions.

## Assumptions and Defaults
1. Placeholder assets are acceptable for first complete v1 implementation.
2. Story text is English in v1.
3. macOS is the only required release platform in this iteration.
4. No key remapping/settings screen in v1 (fixed controls only).
5. Level editing uses JSON only; no Tiled editor integration in v1.
6. Save/load persistence is out of scope; campaign progress resets on app restart.

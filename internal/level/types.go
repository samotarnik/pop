package level

type EnemySpawn struct {
	X      float64 `json:"x"`
	Kind   string  `json:"kind"`
	Health int     `json:"health"`
}

type HazardSpawn struct {
	Type  string  `json:"type"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Speed float64 `json:"speed"`
}

type PickupSpawn struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

type LevelData struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Width       float64       `json:"width"`
	Height      float64       `json:"height"`
	GroundY     float64       `json:"ground_y"`
	PlayerSpawn float64       `json:"player_spawn_x"`
	GoalX       float64       `json:"goal_x"`
	NextStoryID string        `json:"next_story_id"`
	BGM         string        `json:"bgm"`
	Enemies     []EnemySpawn  `json:"enemies"`
	Hazards     []HazardSpawn `json:"hazards"`
	Pickups     []PickupSpawn `json:"pickups"`
}

type StoryCard struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

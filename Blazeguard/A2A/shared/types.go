package shared

type A2AMessage struct {
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	EventType string                 `json:"event_type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp string                 `json:"timestamp"`
}

type FireEvent struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Confidence string  `json:"confidence"`
	Brightness float64 `json:"brightness"`
	ZoneID     string  `json:"zone_id"`
}

type PredictionEvent struct {
	ZoneID    string     `json:"zone_id"`
	Corridors []Corridor `json:"corridors"`
	WindSpeed float64    `json:"wind_speed"`
}

type Corridor struct {
	Path        string  `json:"path"`
	ETA         string  `json:"eta_minutes"`
	Probability float64 `json:"probability"`
}
type LogisticsEvent struct {
	ZoneID    string     `json:"zone_id"`
	Resources []Resource `json:"resources"`
}


type Resource struct {
	UnitID string `json:"unit_id"`
	Type   string `json:"type"`
	Route  string `json:"route"`
	ETA    string `json:"eta"`
}
type Report struct {
	Report      string `json:"report"`
	Severity    string `json:"severity"`
	GeneratedAt string `json:"generated_at"`
}

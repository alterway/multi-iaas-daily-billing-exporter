package app

type TargetResponse struct {
	Status string     `json:"status"`
	Data   TargetData `json:"data"`
}

type TargetData struct {
	ActiveTargets  []TargetSpec `json:"activeTargets"`
	DroppedTargets []TargetSpec `json:"droppedTargets"`
}

type TargetSpec struct {
	DiscoveredLabels   map[string]interface{} `json:"discoveredLabels"`
	Labels             map[string]interface{} `json:"labels"`
	ScrapePool         string                 `json:"scrapePool"`
	ScrapeUrl          string                 `json:"scrapeUrl"`
	GlobalUrl          string                 `json:"globalUrl"`
	LastError          string                 `json:"lastError"`
	LastScrape         string                 `json:"lastScrape"`
	LastScrapeDuration string                 `json:"lastScrapeDuration"`
	Health             string                 `json:"health"`
}

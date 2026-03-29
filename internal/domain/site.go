package domain

import "time"

type Site struct {
	ID            uint64    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	URL           string    `json:"url"`
	Domain        string    `json:"domain"`
	Status        string    `json:"status"`
	AddedBy       string    `json:"added_by"`
	Tags          []string  `json:"tags"`
	CreatedAt     time.Time `json:"created_at"`
	BookmarkedBy  int       `json:"bookmarked_by"`
	StrategyLabel string    `json:"strategy_label,omitempty"`
}

type SiteFilter struct {
	Strategy string
	Query    string
	Tag      string
	Status   string
}

type SiteListResult struct {
	Sites                []Site                `json:"sites"`
	Strategy             string                `json:"strategy"`
	AvailableTags        []string              `json:"available_tags"`
	MonthlyRegistrations []MonthlyRegistration `json:"monthly_registrations"`
}

type MonthlyRegistration struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

package service

import (
	"slices"
	"strings"

	"push-link-v2/internal/domain"
)

type SiteStrategy interface {
	Name() string
	Apply([]domain.Site, domain.SiteFilter) []domain.Site
}

type SiteStrategyRegistry struct {
	strategies map[string]SiteStrategy
}

func NewSiteStrategyRegistry() *SiteStrategyRegistry {
	items := []SiteStrategy{
		DefaultSiteStrategy{},
		PublishedSiteStrategy{},
		TrendingSiteStrategy{},
	}

	strategies := make(map[string]SiteStrategy, len(items))
	for _, strategy := range items {
		strategies[strategy.Name()] = strategy
	}

	return &SiteStrategyRegistry{strategies: strategies}
}

func (r *SiteStrategyRegistry) Resolve(name string) SiteStrategy {
	if strategy, ok := r.strategies[name]; ok {
		return strategy
	}
	return r.strategies["default"]
}

type DefaultSiteStrategy struct{}

func (DefaultSiteStrategy) Name() string {
	return "default"
}

func (DefaultSiteStrategy) Apply(sites []domain.Site, filter domain.SiteFilter) []domain.Site {
	result := make([]domain.Site, 0, len(sites))
	for _, site := range sites {
		if !matchesQuery(site, filter.Query) {
			continue
		}
		if filter.Tag != "" && !contains(site.Tags, filter.Tag) {
			continue
		}
		if filter.Status != "" && site.Status != filter.Status {
			continue
		}
		site.StrategyLabel = "Default"
		result = append(result, site)
	}

	slices.SortFunc(result, func(a, b domain.Site) int {
		return strings.Compare(a.Title, b.Title)
	})

	return result
}

type PublishedSiteStrategy struct{}

func (PublishedSiteStrategy) Name() string {
	return "published"
}

func (PublishedSiteStrategy) Apply(sites []domain.Site, filter domain.SiteFilter) []domain.Site {
	filter.Status = "published"
	result := DefaultSiteStrategy{}.Apply(sites, filter)
	for index := range result {
		result[index].StrategyLabel = "Published Only"
	}
	return result
}

type TrendingSiteStrategy struct{}

func (TrendingSiteStrategy) Name() string {
	return "trending"
}

func (TrendingSiteStrategy) Apply(sites []domain.Site, filter domain.SiteFilter) []domain.Site {
	filter.Status = ""
	result := DefaultSiteStrategy{}.Apply(sites, filter)
	slices.SortFunc(result, func(a, b domain.Site) int {
		if a.BookmarkedBy == b.BookmarkedBy {
			if a.CreatedAt.Equal(b.CreatedAt) {
				return strings.Compare(a.Title, b.Title)
			}
			if a.CreatedAt.After(b.CreatedAt) {
				return -1
			}
			return 1
		}
		if a.BookmarkedBy > b.BookmarkedBy {
			return -1
		}
		return 1
	})
	for index := range result {
		result[index].StrategyLabel = "Trending"
	}
	return result
}

func matchesQuery(site domain.Site, query string) bool {
	if query == "" {
		return true
	}

	haystack := strings.ToLower(strings.Join([]string{
		site.Title,
		site.Description,
		site.Domain,
		strings.Join(site.Tags, " "),
	}, " "))

	return strings.Contains(haystack, strings.ToLower(query))
}

func contains(items []string, expected string) bool {
	for _, item := range items {
		if item == expected {
			return true
		}
	}
	return false
}

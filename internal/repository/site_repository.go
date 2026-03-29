package repository

import (
	"slices"
	"time"

	"push-link-v2/internal/domain"
)

type SiteRepository interface {
	List() ([]domain.Site, error)
	Tags() ([]string, error)
}

type MemorySiteRepository struct {
	sites []domain.Site
	tags  []string
}

func NewMemorySiteRepository(sites []domain.Site) *MemorySiteRepository {
	tagSet := make(map[string]struct{})
	for _, site := range sites {
		for _, tag := range site.Tags {
			tagSet[tag] = struct{}{}
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	slices.Sort(tags)

	return &MemorySiteRepository{sites: sites, tags: tags}
}

func (r *MemorySiteRepository) List() ([]domain.Site, error) {
	return slices.Clone(r.sites), nil
}

func (r *MemorySiteRepository) Tags() ([]string, error) {
	return slices.Clone(r.tags), nil
}

func SeedSites() []domain.Site {
	return []domain.Site{
		{
			ID:           1,
			Title:        "MDN Web Docs",
			Description:  "HTML、CSS、JavaScript の一次情報を横断できる定番ドキュメント。",
			URL:          "https://developer.mozilla.org/",
			Domain:       "developer.mozilla.org",
			Status:       "published",
			AddedBy:      "Toyo Admin",
			Tags:         []string{"development", "learning"},
			CreatedAt:    time.Date(2026, 1, 15, 10, 0, 0, 0, time.Local),
			BookmarkedBy: 14,
		},
		{
			ID:           2,
			Title:        "Stack Overflow",
			Description:  "実運用の詰まりどころを素早く解決するための Q&A 集積地。",
			URL:          "https://stackoverflow.com/",
			Domain:       "stackoverflow.com",
			Status:       "published",
			AddedBy:      "Mika Curator",
			Tags:         []string{"development", "productivity"},
			CreatedAt:    time.Date(2026, 1, 28, 18, 30, 0, 0, time.Local),
			BookmarkedBy: 20,
		},
		{
			ID:           3,
			Title:        "Figma",
			Description:  "UI 設計とレビューを一気通貫で進めるための共同編集ツール。",
			URL:          "https://www.figma.com/",
			Domain:       "www.figma.com",
			Status:       "published",
			AddedBy:      "Mika Curator",
			Tags:         []string{"design", "productivity"},
			CreatedAt:    time.Date(2026, 2, 10, 9, 0, 0, 0, time.Local),
			BookmarkedBy: 11,
		},
		{
			ID:           4,
			Title:        "freeCodeCamp",
			Description:  "体系的に学習できる無料教材。初学者向け導線が強い。",
			URL:          "https://www.freecodecamp.org/",
			Domain:       "www.freecodecamp.org",
			Status:       "published",
			AddedBy:      "Toyo Admin",
			Tags:         []string{"development", "learning"},
			CreatedAt:    time.Date(2026, 2, 23, 13, 0, 0, 0, time.Local),
			BookmarkedBy: 9,
		},
		{
			ID:           5,
			Title:        "OpenAI Developers",
			Description:  "AI プロダクト実装のための API ドキュメントと設計ガイド。",
			URL:          "https://platform.openai.com/docs/overview",
			Domain:       "platform.openai.com",
			Status:       "draft",
			AddedBy:      "Toyo Admin",
			Tags:         []string{"ai", "development"},
			CreatedAt:    time.Date(2026, 3, 18, 21, 0, 0, 0, time.Local),
			BookmarkedBy: 5,
		},
	}
}

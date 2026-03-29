package usecase

import (
	"fmt"
	"sort"
	"time"

	"push-link-v2/internal/domain"
	"push-link-v2/internal/repository"
	"push-link-v2/internal/service"
)

type SiteUsecase struct {
	repository repository.SiteRepository
	registry   *service.SiteStrategyRegistry
}

func NewSiteUsecase(repository repository.SiteRepository, registry *service.SiteStrategyRegistry) *SiteUsecase {
	return &SiteUsecase{
		repository: repository,
		registry:   registry,
	}
}

func (u *SiteUsecase) List(filter domain.SiteFilter) (domain.SiteListResult, error) {
	sites, err := u.repository.List()
	if err != nil {
		return domain.SiteListResult{}, err
	}

	tags, err := u.repository.Tags()
	if err != nil {
		return domain.SiteListResult{}, err
	}

	strategy := u.registry.Resolve(filter.Strategy)
	filter.Strategy = strategy.Name()

	return domain.SiteListResult{
		Sites:                strategy.Apply(sites, filter),
		Strategy:             filter.Strategy,
		AvailableTags:        tags,
		MonthlyRegistrations: monthlyRegistrations(sites),
	}, nil
}

func monthlyRegistrations(sites []domain.Site) []domain.MonthlyRegistration {
	counts := make(map[string]int)
	months := make(map[string]time.Time)

	for _, site := range sites {
		month := time.Date(site.CreatedAt.Year(), site.CreatedAt.Month(), 1, 0, 0, 0, 0, site.CreatedAt.Location())
		key := month.Format("2006-01")
		counts[key]++
		months[key] = month
	}

	keys := make([]string, 0, len(months))
	for key := range months {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return months[keys[i]].Before(months[keys[j]])
	})

	result := make([]domain.MonthlyRegistration, 0, len(keys))
	for _, key := range keys {
		result = append(result, domain.MonthlyRegistration{
			Label: fmt.Sprintf("%d/%02d", months[key].Year(), months[key].Month()),
			Count: counts[key],
		})
	}

	return result
}

package services

import (
	"context"
	"strings"
	"time"

	"github.com/skygenesisenterprise/kami-sama/server/src/interfaces"
	"github.com/skygenesisenterprise/kami-sama/server/src/models"
	"github.com/skygenesisenterprise/kami-sama/server/src/utils"
)

type AnimeService struct {
	repos *Repositories
}

func NewAnimeService(repos *Repositories) *AnimeService {
	return &AnimeService{repos: repos}
}

func (s *AnimeService) Create(ctx context.Context, userID, title, jpTitle, synopsis, coverImg, bannerImg, trailer, status string, totalEpisodes, releaseYear int, season, source, ageRating string, genreIDs, studioIDs []string) (*models.Anime, error) {
	if strings.TrimSpace(title) == "" {
		return nil, utils.ErrValidationFailed
	}
	slug := generateSlug(title)
	now := time.Now().UTC()
	anime := &models.Anime{
		Common:         models.Common{ID: utils.NewID(), CreatedAt: now, UpdatedAt: now},
		Slug:           slug,
		Title:          strings.TrimSpace(title),
		JapaneseTitle:  jpTitle,
		Synopsis:       synopsis,
		CoverImageUrl:  coverImg,
		BannerImageUrl: bannerImg,
		TrailerUrl:     trailer,
		Status:         defaultString(status, "upcoming"),
		TotalEpisodes:  totalEpisodes,
		ReleaseYear:    releaseYear,
		Season:         season,
		Source:         source,
		AgeRating:      ageRating,
	}
	if err := s.repos.Anime().Create(ctx, anime); err != nil {
		return nil, err
	}
	for _, genreID := range genreIDs {
		s.repos.db.Exec("INSERT INTO anime_genres (anime_id, genre_id) VALUES (?, ?) ON CONFLICT DO NOTHING", anime.ID, genreID)
	}
	for _, studioID := range studioIDs {
		s.repos.db.Exec("INSERT INTO anime_studios (anime_id, studio_id) VALUES (?, ?) ON CONFLICT DO NOTHING", anime.ID, studioID)
	}
	return anime, nil
}

func (s *AnimeService) GetByID(ctx context.Context, id string) (*models.Anime, error) {
	return s.repos.Anime().GetByID(ctx, id)
}

func (s *AnimeService) GetBySlug(ctx context.Context, slug string) (*models.Anime, error) {
	return s.repos.Anime().GetBySlug(ctx, slug)
}

func (s *AnimeService) List(ctx context.Context, opts interfaces.ListAnimeOpts) ([]models.Anime, int64, error) {
	return s.repos.Anime().List(ctx, opts)
}

type UpdateAnimeInput struct {
	Title          *string  `json:"title"`
	JapaneseTitle  *string  `json:"japaneseTitle"`
	Synopsis       *string  `json:"synopsis"`
	CoverImageUrl  *string  `json:"coverImageUrl"`
	BannerImageUrl *string  `json:"bannerImageUrl"`
	TrailerUrl     *string  `json:"trailerUrl"`
	Status         *string  `json:"status"`
	Rating         *float64 `json:"rating"`
	TotalEpisodes  *int     `json:"totalEpisodes"`
	ReleaseYear    *int     `json:"releaseYear"`
	Season         *string  `json:"season"`
	Source         *string  `json:"source"`
	AgeRating      *string  `json:"ageRating"`
	IsFeatured     *bool    `json:"isFeatured"`
	IsTrending     *bool    `json:"isTrending"`
	GenreIDs       []string `json:"genreIds"`
	StudioIDs      []string `json:"studioIds"`
}

func (s *AnimeService) Update(ctx context.Context, userID, id string, input UpdateAnimeInput) (*models.Anime, error) {
	anime, err := s.repos.Anime().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Title != nil {
		anime.Title = strings.TrimSpace(*input.Title)
	}
	if input.JapaneseTitle != nil {
		anime.JapaneseTitle = *input.JapaneseTitle
	}
	if input.Synopsis != nil {
		anime.Synopsis = *input.Synopsis
	}
	if input.CoverImageUrl != nil {
		anime.CoverImageUrl = *input.CoverImageUrl
	}
	if input.BannerImageUrl != nil {
		anime.BannerImageUrl = *input.BannerImageUrl
	}
	if input.TrailerUrl != nil {
		anime.TrailerUrl = *input.TrailerUrl
	}
	if input.Status != nil {
		anime.Status = *input.Status
	}
	if input.Rating != nil {
		anime.Rating = *input.Rating
	}
	if input.TotalEpisodes != nil {
		anime.TotalEpisodes = *input.TotalEpisodes
	}
	if input.ReleaseYear != nil {
		anime.ReleaseYear = *input.ReleaseYear
	}
	if input.Season != nil {
		anime.Season = *input.Season
	}
	if input.Source != nil {
		anime.Source = *input.Source
	}
	if input.AgeRating != nil {
		anime.AgeRating = *input.AgeRating
	}
	if input.IsFeatured != nil {
		anime.IsFeatured = *input.IsFeatured
	}
	if input.IsTrending != nil {
		anime.IsTrending = *input.IsTrending
	}
	anime.UpdatedAt = time.Now().UTC()
	if err := s.repos.Anime().Update(ctx, anime); err != nil {
		return nil, err
	}
	if input.GenreIDs != nil {
		s.repos.db.Exec("DELETE FROM anime_genres WHERE anime_id = ?", anime.ID)
		for _, genreID := range input.GenreIDs {
			s.repos.db.Exec("INSERT INTO anime_genres (anime_id, genre_id) VALUES (?, ?) ON CONFLICT DO NOTHING", anime.ID, genreID)
		}
	}
	if input.StudioIDs != nil {
		s.repos.db.Exec("DELETE FROM anime_studios WHERE anime_id = ?", anime.ID)
		for _, studioID := range input.StudioIDs {
			s.repos.db.Exec("INSERT INTO anime_studios (anime_id, studio_id) VALUES (?, ?) ON CONFLICT DO NOTHING", anime.ID, studioID)
		}
	}
	return anime, nil
}

func (s *AnimeService) Delete(ctx context.Context, id string) error {
	return s.repos.Anime().Delete(ctx, id)
}

func (s *AnimeService) Search(ctx context.Context, query string, limit int) ([]models.Anime, error) {
	return s.repos.Anime().Search(ctx, query, limit)
}

func (s *AnimeService) ListEpisodes(ctx context.Context, animeID string) ([]models.Episode, error) {
	return s.repos.Episodes().ListByAnime(ctx, animeID, nil)
}

func (s *AnimeService) ListSeasons(ctx context.Context, animeID string) ([]models.AnimeSeason, error) {
	var seasons []models.AnimeSeason
	err := s.repos.db.WithContext(ctx).Where("anime_id = ?", animeID).Order("number asc").Find(&seasons).Error
	return seasons, err
}

func (s *AnimeService) ListReviews(ctx context.Context, animeID string, limit, offset int) ([]models.Review, int64, error) {
	return s.repos.Reviews().ListByAnime(ctx, animeID, limit, offset)
}

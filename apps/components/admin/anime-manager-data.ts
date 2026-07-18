import type { AnimeStatus } from '@/types/anime'
import { getAllAnime } from '@/lib/mock-data'

export type AnimeManagerItem = {
  id: string
  slug: string
  title: string
  japaneseTitle: string
  studio: string
  year: number
  status: AnimeStatus
  genres: string[]
  rating: number
  ratingCount: number
  totalEpisodes: number
  ageRating: string
  cover: string
}

export function getAnimeManagerData(): AnimeManagerItem[] {
  return getAllAnime().map((a) => ({
    id: a.id,
    slug: a.slug,
    title: a.title,
    japaneseTitle: a.japaneseTitle,
    studio: a.studio.name,
    year: a.year,
    status: a.status,
    genres: a.genres.map((g) => g.name),
    rating: a.rating,
    ratingCount: a.ratingCount,
    totalEpisodes: a.totalEpisodes,
    ageRating: a.ageRating,
    cover: a.cover,
  }))
}

export const statusLabels: Record<AnimeStatus, string> = {
  airing: 'En cours',
  completed: 'Terminé',
  upcoming: 'À venir',
  hiatus: 'En pause',
}

export const statusStyles: Record<AnimeStatus, string> = {
  airing: 'bg-chart-2/15 text-chart-2',
  completed: 'bg-primary/15 text-primary',
  upcoming: 'bg-chart-3/15 text-chart-3',
  hiatus: 'bg-muted text-muted-foreground',
}

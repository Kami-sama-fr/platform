'use client'

import { use, useMemo } from 'react'
import { useSearchParams } from 'next/navigation'
import Image from 'next/image'
import Link from 'next/link'
import {
  ChevronLeft,
  ChevronRight,
  Clock,
  List,
  Maximize,
  MessageSquare,
  Play,
  Share2,
  SkipForward,
  Volume2,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Separator } from '@/components/ui/separator'
import { ScrollArea } from '@/components/ui/scroll-area'
import { GenreTag } from '@/components/kami/genre-tag'
import { RatingBadge } from '@/components/kami/rating-badge'
import { cn } from '@/lib/utils'
import {
  getAnime,
  getEpisode,
  getEpisodes,
  getReviews,
  formatDuration,
} from '@/lib/mock-data'

export default function WatchPage({
  params,
}: {
  params: Promise<{ slug: string }>
}) {
  const { slug } = use(params)
  const searchParams = useSearchParams()
  const epId = searchParams.get('ep')

  const anime = getAnime(slug)
  const allEpisodes = useMemo(() => (anime ? getEpisodes(anime.id) : []), [anime])

  const currentEpisode = useMemo(() => {
    if (!anime) return null
    if (epId) return getEpisode(anime.id, epId) ?? allEpisodes[0]
    return allEpisodes[0]
  }, [anime, epId, allEpisodes])

  const currentIndex = currentEpisode
    ? allEpisodes.findIndex((e) => e.id === currentEpisode.id)
    : -1
  const prevEpisode = currentIndex > 0 ? allEpisodes[currentIndex - 1] : null
  const nextEpisode =
    currentIndex < allEpisodes.length - 1 ? allEpisodes[currentIndex + 1] : null

  const reviews = useMemo(
    () => (anime ? getReviews(anime.id) : []),
    [anime],
  )

  // Group by season for the episode list
  const seasons = useMemo(() => {
    if (!anime) return []
    const map = new Map<number, typeof allEpisodes>()
    for (const ep of allEpisodes) {
      const list = map.get(ep.season) ?? []
      list.push(ep)
      map.set(ep.season, list)
    }
    return Array.from(map.entries()).map(([season, eps]) => ({
      season,
      title: anime.seasons.find((s) => s.number === season)?.title ?? `Season ${season}`,
      episodes: eps,
    }))
  }, [anime, allEpisodes])

  // Current season index for prev/next season navigation
  const currentSeasonIdx = currentEpisode
    ? seasons.findIndex((s) => s.season === currentEpisode.season)
    : -1

  if (!anime || !currentEpisode) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <h1 className="font-display text-2xl font-bold">Episode not found</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            The episode you&apos;re looking for doesn&apos;t exist.
          </p>
          <Button asChild variant="secondary" className="mt-4">
            <Link href="/">Go Home</Link>
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      {/* ── Video Player ──────────────────────────────────────────── */}
      <div className="relative aspect-video w-full bg-black">
        {/* Placeholder poster */}
        <img
          src={anime.banner || '/placeholder.svg'}
          alt=""
          className="absolute inset-0 size-full object-cover opacity-25"
        />

        {/* Center play */}
        <div className="absolute inset-0 z-10 flex flex-col items-center justify-center">
          <button
            type="button"
            className="flex size-16 items-center justify-center rounded-full bg-primary/90 text-primary-foreground shadow-2xl shadow-primary/30 transition-transform hover:scale-110"
          >
            <Play className="ml-1 size-7 fill-current" />
          </button>
          <p className="mt-3 text-sm text-white/60">
            {formatDuration(currentEpisode.duration)}
          </p>
        </div>

        {/* Controls bar */}
        <div className="absolute inset-x-0 bottom-0 z-20 bg-gradient-to-t from-black/80 to-transparent p-4">
          <div className="mx-auto max-w-5xl">
            <Progress
              value={35}
              className="h-1.5 cursor-pointer rounded-full bg-white/20 [&>div]:bg-primary"
            />
            <div className="mt-2 flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-8 text-white hover:text-white/80"
                >
                  <Play className="size-4 fill-current" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-8 text-white hover:text-white/80"
                >
                  <SkipForward className="size-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-8 text-white hover:text-white/80"
                >
                  <Volume2 className="size-4" />
                </Button>
                <span className="ml-1 text-xs text-white/60">
                  8:24 / {formatDuration(currentEpisode.duration)}
                </span>
              </div>
              <div className="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-8 text-white hover:text-white/80"
                >
                  <MessageSquare className="size-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-8 text-white hover:text-white/80"
                >
                  <Maximize className="size-4" />
                </Button>
              </div>
            </div>
          </div>
        </div>

        {/* Prev / Next episode overlay */}
        {prevEpisode && (
          <Link
            href={`/watch/${slug}?ep=${prevEpisode.id}`}
            className="absolute left-3 top-1/2 z-30 -translate-y-1/2 rounded-full bg-black/50 p-2.5 text-white backdrop-blur-sm transition-colors hover:bg-black/70"
            aria-label={`Previous episode: ${prevEpisode.title}`}
          >
            <ChevronLeft className="size-5" />
          </Link>
        )}
        {nextEpisode && (
          <Link
            href={`/watch/${slug}?ep=${nextEpisode.id}`}
            className="absolute right-3 top-1/2 z-30 -translate-y-1/2 rounded-full bg-black/50 p-2.5 text-white backdrop-blur-sm transition-colors hover:bg-black/70"
            aria-label={`Next episode: ${nextEpisode.title}`}
          >
            <ChevronRight className="size-5" />
          </Link>
        )}
      </div>

      {/* ── Content ───────────────────────────────────────────────── */}
      <div className="mx-auto max-w-7xl px-4 py-6 md:px-8">
        <div className="grid gap-8 lg:grid-cols-[1fr_340px]">
          {/* Main */}
          <div className="min-w-0">
            {/* Header */}
            <div className="flex items-start justify-between gap-4">
              <div className="min-w-0">
                <div className="flex items-center gap-2 text-sm">
                  <Link
                    href={`/anime/${slug}`}
                    className="font-medium text-primary hover:underline"
                  >
                    {anime.title}
                  </Link>
                  <span className="text-muted-foreground">·</span>
                  <span className="text-muted-foreground">
                    Season {currentEpisode.season}
                  </span>
                </div>
                <h1 className="mt-1.5 font-display text-xl font-bold md:text-2xl">
                  Episode {currentEpisode.number}: {currentEpisode.title}
                </h1>
              </div>

              <div className="flex shrink-0 gap-2">
                {prevEpisode && (
                  <Button variant="secondary" size="sm" asChild>
                    <Link href={`/watch/${slug}?ep=${prevEpisode.id}`}>
                      <ChevronLeft className="mr-1 size-4" />
                      Prev
                    </Link>
                  </Button>
                )}
                {nextEpisode && (
                  <Button size="sm" asChild>
                    <Link href={`/watch/${slug}?ep=${nextEpisode.id}`}>
                      Next
                      <ChevronRight className="ml-1 size-4" />
                    </Link>
                  </Button>
                )}
              </div>
            </div>

            {/* Meta */}
            <div className="mt-3 flex flex-wrap items-center gap-2.5">
              <RatingBadge rating={anime.rating} />
              <Badge variant="secondary" className="text-xs">
                {anime.ageRating}
              </Badge>
              <Badge variant="secondary" className="text-xs">
                {formatDuration(currentEpisode.duration)}
              </Badge>
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <Clock className="size-3" />
                {currentEpisode.releaseDate}
              </span>
              <Button variant="ghost" size="icon" className="ml-auto size-8">
                <Share2 className="size-4" />
              </Button>
            </div>

            {/* Description */}
            {currentEpisode.description && (
              <p className="mt-4 text-sm leading-relaxed text-muted-foreground">
                {currentEpisode.description}
              </p>
            )}

            <Separator className="my-6" />

            {/* Episode list by season */}
            <div className="space-y-6">
              {seasons.map((season, si) => (
                <section key={season.season}>
                  <div className="mb-3 flex items-center justify-between">
                    <div>
                      <h3 className="font-display text-sm font-semibold">
                        {season.title}
                      </h3>
                      <p className="text-xs text-muted-foreground">
                        {season.episodes.length} episodes
                      </p>
                    </div>
                    {si > 0 && (
                      <Button variant="ghost" size="sm" className="text-xs" asChild>
                        <Link
                          href={`/watch/${slug}?ep=${season.episodes[0].id}`}
                        >
                          Watch Season
                        </Link>
                      </Button>
                    )}
                  </div>
                  <div className="space-y-1.5">
                    {season.episodes.map((ep) => {
                      const isActive = ep.id === currentEpisode.id
                      return (
                        <Link
                          key={ep.id}
                          href={`/watch/${slug}?ep=${ep.id}`}
                          className={cn(
                            'group flex items-center gap-3 rounded-lg border p-2 transition-colors',
                            isActive
                              ? 'border-primary/40 bg-primary/10'
                              : 'border-transparent hover:bg-accent',
                          )}
                        >
                          <div className="relative h-12 w-20 shrink-0 overflow-hidden rounded-md">
                            <img
                              src={ep.thumbnail || anime.cover || '/placeholder.svg'}
                              alt={ep.title}
                              className="size-full object-cover"
                            />
                            {isActive && (
                              <div className="absolute inset-0 flex items-center justify-center bg-primary/30">
                                <Play className="size-4 fill-white text-white" />
                              </div>
                            )}
                          </div>
                          <div className="min-w-0 flex-1">
                            <span
                              className={cn(
                                'text-xs font-medium',
                                isActive ? 'text-primary' : 'text-primary/60',
                              )}
                            >
                              Ep {ep.number}
                              {isActive && ' — Now Playing'}
                            </span>
                            <h4 className="truncate text-sm font-medium">
                              {ep.title}
                            </h4>
                          </div>
                          <span className="shrink-0 text-xs text-muted-foreground">
                            {formatDuration(ep.duration)}
                          </span>
                        </Link>
                      )
                    })}
                  </div>
                </section>
              ))}
            </div>
          </div>

          {/* ── Sidebar ─────────────────────────────────────────────── */}
          <aside className="hidden lg:block">
            <div className="sticky top-20 space-y-4">
              {/* Anime card */}
              <Link
                href={`/anime/${slug}`}
                className="group block overflow-hidden rounded-xl border border-border/60 bg-card"
              >
                <div className="relative aspect-[2/3] overflow-hidden">
                  <Image
                    src={anime.cover || '/placeholder.svg'}
                    alt={anime.title}
                    fill
                    sizes="340px"
                    className="object-cover transition-transform duration-500 group-hover:scale-105"
                  />
                  <div className="absolute inset-0 bg-gradient-to-t from-background/80 via-transparent to-transparent" />
                  <div className="absolute inset-x-0 bottom-0 p-4">
                    <h3 className="font-display text-base font-bold">
                      {anime.title}
                    </h3>
                    <p className="text-xs text-muted-foreground">
                      {anime.japaneseTitle}
                    </p>
                  </div>
                </div>
              </Link>

              {/* Info */}
              <div className="rounded-xl border border-border/60 bg-card p-4">
                <div className="grid grid-cols-2 gap-3 text-sm">
                  <div>
                    <span className="text-muted-foreground">Rating</span>
                    <p className="mt-0.5 font-semibold text-gold">
                      ★ {anime.rating}
                    </p>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Year</span>
                    <p className="mt-0.5 font-semibold">{anime.year}</p>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Status</span>
                    <p className="mt-0.5 font-semibold capitalize">
                      {anime.status}
                    </p>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Episodes</span>
                    <p className="mt-0.5 font-semibold">
                      {anime.totalEpisodes}
                    </p>
                  </div>
                </div>
                <Separator className="my-3" />
                <div className="flex flex-wrap gap-2">
                  {anime.genres.map((genre) => (
                    <GenreTag key={genre.id} genre={genre} asLink={false} />
                  ))}
                </div>
              </div>

              {/* Quick reviews */}
              {reviews.length > 0 && (
                <div className="rounded-xl border border-border/60 bg-card p-4">
                  <h4 className="mb-3 text-sm font-semibold">Recent Reviews</h4>
                  <div className="space-y-3">
                    {reviews.slice(0, 2).map((r) => (
                      <div key={r.id}>
                        <div className="flex items-center gap-2">
                          <div className="flex size-6 items-center justify-center rounded-full bg-secondary text-[10px] font-semibold">
                            {r.user.displayName[0]}
                          </div>
                          <span className="text-xs font-medium">
                            {r.user.displayName}
                          </span>
                          <Badge
                            variant="secondary"
                            className="ml-auto text-[10px] text-gold"
                          >
                            ★ {r.rating}
                          </Badge>
                        </div>
                        <p className="mt-1 line-clamp-2 text-xs text-muted-foreground">
                          {r.body}
                        </p>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </aside>
        </div>
      </div>
    </div>
  )
}

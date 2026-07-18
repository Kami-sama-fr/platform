import { Flame, Sparkles, TrendingUp, Tv } from 'lucide-react'
import Link from 'next/link'
import { HeroBanner } from '@/components/kami/hero-banner'
import { CarouselSection } from '@/components/kami/carousel-section'
import { AnimeCard } from '@/components/kami/anime-card'
import { ContinueWatchingCard } from '@/components/kami/continue-watching-card'
import { CommunityCard } from '@/components/kami/community-card'
import { RecommendationCard } from '@/components/kami/recommendation-card'
import { SimulcastCard } from '@/components/kami/simulcast-card'
import { SeasonalPickCard } from '@/components/kami/seasonal-pick-card'
import { SectionDivider } from '@/components/kami/section-divider'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  GENRES,
  getAnime,
  getCommunityPicks,
  getContinueWatching,
  getRecentlyAdded,
  getRecommendations,
  getSeasonalPicks,
  getSimulcasts,
  getTopPicks,
  getTrending,
} from '@/lib/mock-data'

export default function HomePage() {
  const featured = ['crimson-blade', 'eternal-frost', 'neon-orbit'].map((id) => getAnime(id)!)
  const topPicks = getTopPicks()
  const continueWatching = getContinueWatching()
  const simulcasts = getSimulcasts()
  const trending = getTrending()
  const recentlyAdded = getRecentlyAdded()
  const seasonalPicks = getSeasonalPicks()
  const communityPicks = getCommunityPicks()
  const recommendations = getRecommendations()

  return (
    <div className="pb-12">
      {/* ── Hero Banner ─────────────────────────────────────────────── */}
      <HeroBanner items={featured} />

      {/* ── Genre Navigation Strip ──────────────────────────────────── */}
      <nav className="border-b border-line/60 bg-paper/95 backdrop-blur supports-backdrop-filter:bg-paper/80">
        <div className="scrollbar-hide flex gap-1 overflow-x-auto px-4 py-2.5 md:px-8">
          {GENRES.map((genre) => (
            <Link
              key={genre.id}
              href={`/browse?genre=${genre.id}`}
              className="shrink-0 rounded-full border border-line/60 bg-paper-dim/50 px-3.5 py-1.5 text-xs font-medium text-ink-wash transition-colors hover:border-stamp/40 hover:bg-stamp/10 hover:text-ink"
            >
              {genre.name}
            </Link>
          ))}
        </div>
      </nav>

      <div className="relative z-10">
        {/* ── Top Picks for You ──────────────────────────────────────── */}
        <section className="py-6">
          <div className="mb-4 px-4 md:px-8">
            <div className="flex items-center gap-2">
              <Sparkles className="size-5 text-stamp" />
              <h2 className="font-display text-lg font-semibold tracking-tight md:text-xl">
                Top Picks for You
              </h2>
            </div>
            <p className="mt-0.5 text-sm text-ink-wash">
              Personalized picks based on your taste
            </p>
          </div>
          <div className="scrollbar-hide flex gap-3 overflow-x-auto px-4 pb-2 md:gap-4 md:px-8">
            {topPicks.map((anime) => (
              <AnimeCard key={anime.id} anime={anime} className="w-37.5 shrink-0 sm:w-42.5 lg:w-46.25" />
            ))}
          </div>
        </section>

        <SectionDivider />

        {/* ── Continue Watching ──────────────────────────────────────── */}
        <CarouselSection
          title="Continue Watching"
          subtitle="Pick up where you left off"
          href="/library"
          itemClassName="w-[280px] sm:w-[310px]"
        >
          {continueWatching.map((item) => (
            <ContinueWatchingCard key={item.anime.id} item={item} />
          ))}
        </CarouselSection>

        <SectionDivider />

        {/* ── Simulcasting Now ───────────────────────────────────────── */}
        <CarouselSection
          title="Simulcasting Now"
          subtitle="Catch new episodes as they air"
          itemClassName="w-[150px] sm:w-[170px] lg:w-[185px]"
        >
          {simulcasts.map((item) => (
            <SimulcastCard key={item.anime.id} item={item} />
          ))}
        </CarouselSection>

        <SectionDivider />

        {/* ── Trending Now ───────────────────────────────────────────── */}
        <CarouselSection
          title="Trending Now"
          subtitle="What everyone is watching this week"
          href="/discover"
        >
          {trending.map((item) => (
            <AnimeCard
              key={item.anime.id}
              anime={item.anime}
              badge={
                <span className="flex size-7 items-center justify-center rounded-md bg-stamp text-sm font-bold text-paper shadow">
                  {item.rank}
                </span>
              }
            />
          ))}
        </CarouselSection>

        <SectionDivider />

        {/* ── Top 10 ────────────────────────────────────────────────── */}
        <section className="py-6">
          <div className="mb-4 flex items-end justify-between gap-4 px-4 md:px-8">
            <div className="min-w-0">
              <div className="flex items-center gap-2">
                <Flame className="size-5 text-stamp" />
                <h2 className="font-display text-lg font-semibold tracking-tight md:text-xl">
                  Top 10
                </h2>
              </div>
              <p className="mt-0.5 line-clamp-1 text-sm text-ink-wash">
                Most popular anime right now
              </p>
            </div>
            <Button
              variant="ghost"
              size="sm"
              className="shrink-0 text-ink-wash"
              asChild
            >
              <Link href="/discover">See all</Link>
            </Button>
          </div>
          <div className="scrollbar-hide flex gap-3 overflow-x-auto px-4 pb-2 md:gap-4 md:px-8">
            {trending.map((item, i) => (
              <Link
                key={item.anime.id}
                href={`/anime/${item.anime.slug}`}
                className="group relative w-30 shrink-0 sm:w-35 md:w-40"
              >
                <div className="pointer-events-none select-none text-[5rem] font-black leading-none text-ink/4 md:text-[7rem]">
                  {i + 1}
                </div>
                <div className="relative -mt-8 aspect-2/3 overflow-hidden rounded-lg border border-line/40 bg-paper-raised md:-mt-12 md:rounded-xl">
                  <img
                    src={item.anime.cover || '/placeholder.svg'}
                    alt={item.anime.title}
                    className="absolute inset-0 size-full object-cover transition-transform duration-500 group-hover:scale-105"
                  />
                  <div className="absolute inset-0 bg-linear-to-t from-paper/80 via-transparent to-transparent" />
                  <div className="absolute right-1.5 top-1.5 z-10">
                    <span className="flex size-6 items-center justify-center rounded bg-stamp text-[10px] font-bold text-paper shadow">
                      {item.rank}
                    </span>
                  </div>
                  <div className="absolute inset-0 flex items-center justify-center opacity-0 transition-opacity duration-300 group-hover:opacity-100">
                    <span className="flex size-10 items-center justify-center rounded-full bg-stamp text-paper shadow-lg">
                      <svg
                        className="size-4 fill-current"
                        viewBox="0 0 24 24"
                        aria-hidden="true"
                      >
                        <path d="M8 5v14l11-7z" />
                      </svg>
                    </span>
                  </div>
                  <div className="absolute inset-x-0 bottom-0 z-10 p-2">
                    <h3 className="line-clamp-2 text-xs font-semibold leading-tight text-balance">
                      {item.anime.title}
                    </h3>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </section>

        <SectionDivider />

        {/* ── New Episodes ───────────────────────────────────────────── */}
        <CarouselSection
          title="New Episodes"
          subtitle="Fresh episodes just added"
          href="/browse"
        >
          {recentlyAdded.map((item) => (
            <AnimeCard
              key={item.anime.id}
              anime={item.anime}
              badge={
                <Badge className="bg-stamp text-paper">
                  EP {item.episode.number}
                </Badge>
              }
            />
          ))}
        </CarouselSection>

        <SectionDivider />

        {/* ── Seasonal Picks ─────────────────────────────────────────── */}
        <CarouselSection
          title="Seasonal Picks"
          subtitle="Curated picks for this season"
          itemClassName="w-[150px] sm:w-[170px] lg:w-[185px]"
        >
          {seasonalPicks.map((pick) => (
            <SeasonalPickCard key={pick.anime.id} pick={pick} />
          ))}
        </CarouselSection>

        <SectionDivider />

        {/* ── Community Picks ────────────────────────────────────────── */}
        <section className="py-6">
          <div className="mb-4 flex items-end justify-between gap-4 px-4 md:px-8">
            <div className="min-w-0">
              <div className="flex items-center gap-2">
                <TrendingUp className="size-5 text-stamp" />
                <h2 className="font-display text-lg font-semibold tracking-tight md:text-xl">
                  Most Requested by the Community
                </h2>
              </div>
              <p className="mt-0.5 line-clamp-1 text-sm text-ink-wash">
                Titles our members are voting to bring to Kami-Sama
              </p>
            </div>
          </div>
          <div className="grid gap-3 px-4 md:grid-cols-2 md:px-8 lg:grid-cols-3">
            {communityPicks.slice(0, 6).map((pick, i) => (
              <CommunityCard key={pick.anime.id} pick={pick} rank={i + 1} />
            ))}
          </div>
        </section>

        <SectionDivider />

        {/* ── Recommended For You ────────────────────────────────────── */}
        <section className="py-6">
          <div className="mb-4 px-4 md:px-8">
            <div className="flex items-center gap-2">
              <Tv className="size-5 text-stamp" />
              <h2 className="font-display text-lg font-semibold tracking-tight md:text-xl">
                Recommended For You
              </h2>
            </div>
            <p className="mt-0.5 text-sm text-ink-wash">
              More titles picked just for you
            </p>
          </div>
          <div className="grid gap-4 px-4 sm:grid-cols-2 md:px-8 lg:grid-cols-4">
            {recommendations.map((rec) => (
              <RecommendationCard key={rec.anime.id} rec={rec} />
            ))}
          </div>
        </section>
      </div>
    </div>
  )
}

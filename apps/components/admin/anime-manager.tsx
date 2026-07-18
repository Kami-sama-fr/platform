'use client'

import * as React from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { cn } from '@/lib/utils'
import {
  getAnimeManagerData,
  statusLabels,
  statusStyles,
  type AnimeManagerItem,
} from './anime-manager-data'
import { GENRES } from '@/lib/mock-data'
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Checkbox } from '@/components/ui/checkbox'
import { Switch } from '@/components/ui/switch'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogAction,
  AlertDialogCancel,
} from '@/components/ui/alert-dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from '@/components/ui/form'
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationPrevious,
  PaginationNext,
  PaginationEllipsis,
} from '@/components/ui/pagination'
import {
  Search,
  Plus,
  MoreHorizontal,
  Pencil,
  Copy,
  Trash2,
  Star,
  Eye,
  Clapperboard,
  TrendingUp,
  ArrowUpDown,
  ChevronDown,
  Image,
  X,
} from 'lucide-react'

const animeSchema = z.object({
  title: z.string().min(1, 'Le titre est requis.'),
  japaneseTitle: z.string().optional(),
  studio: z.string().min(1, 'Le studio est requis.'),
  year: z.coerce.number().min(1990, 'Année invalide.').max(2030, 'Année invalide.'),
  status: z.enum(['airing', 'completed', 'upcoming', 'hiatus']),
  genres: z.array(z.string()).min(1, 'Sélectionnez au moins un genre.'),
  rating: z.coerce.number().min(0, 'Note minimale 0.').max(10, 'Note maximale 10.').optional(),
  totalEpisodes: z.coerce.number().min(1, 'Minimum 1 épisode.').max(9999, 'Trop d\'épisodes.').optional(),
  ageRating: z.string().optional(),
  synopsis: z.string().optional(),
})

type AnimeFormValues = z.infer<typeof animeSchema>

const ITEMS_PER_PAGE = 8

type SortKey = 'title' | 'studio' | 'year' | 'status' | 'rating'
type SortDir = 'asc' | 'desc'

export function AnimeManager() {
  const [items, setItems] = React.useState<AnimeManagerItem[]>(() => getAnimeManagerData())
  const [search, setSearch] = React.useState('')
  const [statusFilter, setStatusFilter] = React.useState<string>('all')
  const [genreFilter, setGenreFilter] = React.useState<string>('all')
  const [sortKey, setSortKey] = React.useState<SortKey>('title')
  const [sortDir, setSortDir] = React.useState<SortDir>('asc')
  const [page, setPage] = React.useState(1)
  const [selected, setSelected] = React.useState<Set<string>>(new Set())
  const [editItem, setEditItem] = React.useState<AnimeManagerItem | null>(null)
  const [createOpen, setCreateOpen] = React.useState(false)
  const [deleteItem, setDeleteItem] = React.useState<AnimeManagerItem | null>(null)
  const [deleteBulkOpen, setDeleteBulkOpen] = React.useState(false)

  const form = useForm<AnimeFormValues>({
    resolver: zodResolver(animeSchema),
    defaultValues: {
      title: '',
      japaneseTitle: '',
      studio: '',
      year: new Date().getFullYear(),
      status: 'airing',
      genres: [],
      rating: 8,
      totalEpisodes: 12,
      ageRating: 'PG-13',
      synopsis: '',
    },
  })

  const filtered = React.useMemo(() => {
    let result = [...items]

    if (search.trim()) {
      const q = search.toLowerCase()
      result = result.filter(
        (a) =>
          a.title.toLowerCase().includes(q) ||
          a.japaneseTitle.toLowerCase().includes(q) ||
          a.studio.toLowerCase().includes(q) ||
          a.genres.some((g) => g.toLowerCase().includes(q)),
      )
    }

    if (statusFilter !== 'all') {
      result = result.filter((a) => a.status === statusFilter)
    }

    if (genreFilter !== 'all') {
      result = result.filter((a) => a.genres.includes(genreFilter))
    }

    result.sort((a, b) => {
      const av = a[sortKey]
      const bv = b[sortKey]
      if (typeof av === 'string' && typeof bv === 'string') {
        return sortDir === 'asc' ? av.localeCompare(bv) : bv.localeCompare(av)
      }
      return sortDir === 'asc' ? (av as number) - (bv as number) : (bv as number) - (av as number)
    })

    return result
  }, [items, search, statusFilter, genreFilter, sortKey, sortDir])

  const totalPages = Math.max(1, Math.ceil(filtered.length / ITEMS_PER_PAGE))
  const pageItems = filtered.slice((page - 1) * ITEMS_PER_PAGE, page * ITEMS_PER_PAGE)

  React.useEffect(() => {
    setPage(1)
  }, [search, statusFilter, genreFilter])

  const toggleSort = (key: SortKey) => {
    if (sortKey === key) {
      setSortDir((d) => (d === 'asc' ? 'desc' : 'asc'))
    } else {
      setSortKey(key)
      setSortDir('asc')
    }
  }

  const toggleSelect = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id)
      else next.add(id)
      return next
    })
  }

  const toggleSelectAll = () => {
    if (selected.size === pageItems.length && pageItems.length > 0) {
      setSelected(new Set())
    } else {
      setSelected(new Set(pageItems.map((a) => a.id)))
    }
  }

  const openCreate = () => {
    setEditItem(null)
    form.reset({
      title: '',
      japaneseTitle: '',
      studio: '',
      year: new Date().getFullYear(),
      status: 'airing',
      genres: [],
      rating: 8,
      totalEpisodes: 12,
      ageRating: 'PG-13',
      synopsis: '',
    })
    setCreateOpen(true)
  }

  const openEdit = (item: AnimeManagerItem) => {
    setEditItem(item)
    form.reset({
      title: item.title,
      japaneseTitle: item.japaneseTitle,
      studio: item.studio,
      year: item.year,
      status: item.status,
      genres: item.genres,
      rating: item.rating,
      totalEpisodes: item.totalEpisodes,
      ageRating: item.ageRating,
      synopsis: '',
    })
    setCreateOpen(true)
  }

  const handleSubmit = (values: AnimeFormValues) => {
    if (editItem) {
      setItems((prev) =>
        prev.map((a) =>
          a.id === editItem.id
            ? {
                ...a,
                title: values.title,
                japaneseTitle: values.japaneseTitle || a.japaneseTitle,
                studio: values.studio,
                year: values.year,
                status: values.status,
                genres: values.genres,
                rating: values.rating ?? a.rating,
                totalEpisodes: values.totalEpisodes ?? a.totalEpisodes,
                ageRating: values.ageRating ?? a.ageRating,
              }
            : a,
        ),
      )
    } else {
      const newId = values.title.toLowerCase().replace(/\s+/g, '-') + '-' + Date.now()
      const newItem: AnimeManagerItem = {
        id: newId,
        slug: newId,
        title: values.title,
        japaneseTitle: values.japaneseTitle || '',
        studio: values.studio,
        year: values.year,
        status: values.status,
        genres: values.genres,
        rating: values.rating ?? 0,
        ratingCount: 0,
        totalEpisodes: values.totalEpisodes ?? 12,
        ageRating: values.ageRating ?? 'PG-13',
        cover: '/covers/placeholder.png',
      }
      setItems((prev) => [newItem, ...prev])
    }
    setCreateOpen(false)
    setEditItem(null)
  }

  const handleDelete = (item: AnimeManagerItem) => {
    setItems((prev) => prev.filter((a) => a.id !== item.id))
    setSelected((prev) => {
      const next = new Set(prev)
      next.delete(item.id)
      return next
    })
    setDeleteItem(null)
  }

  const handleBulkDelete = () => {
    setItems((prev) => prev.filter((a) => !selected.has(a.id)))
    setSelected(new Set())
    setDeleteBulkOpen(false)
  }

  const handleDuplicate = (item: AnimeManagerItem) => {
    const newId = item.slug + '-copy-' + Date.now()
    const dup: AnimeManagerItem = {
      ...item,
      id: newId,
      slug: newId,
      title: item.title + ' (copie)',
      ratingCount: 0,
    }
    setItems((prev) => [dup, ...prev])
  }

  const stats = React.useMemo(() => {
    const total = items.length
    const airing = items.filter((a) => a.status === 'airing').length
    const completed = items.filter((a) => a.status === 'completed').length
    const upcoming = items.filter((a) => a.status === 'upcoming').length
    return { total, airing, completed, upcoming }
  }, [items])

  return (
    <div className="space-y-6">
      {/* KPI Cards */}
      <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <KpiCard icon={Clapperboard} label="Total animés" value={stats.total} hint="Dans le catalogue" />
        <KpiCard icon={Eye} label="En cours" value={stats.airing} hint="Diffusion active" accent="chart-2" />
        <KpiCard icon={Star} label="Terminés" value={stats.completed} hint="Séries complètes" accent="chart-1" />
        <KpiCard icon={TrendingUp} label="À venir" value={stats.upcoming} hint="Prochainement" accent="chart-3" />
      </section>

      {/* Toolbar */}
      <section className="rounded-xl border border-border bg-card p-4">
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex flex-1 items-center gap-2">
            <div className="relative flex-1 sm:max-w-xs">
              <Search className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Rechercher un animé…"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9"
              />
            </div>
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger size="sm" className="w-[140px]">
                <SelectValue placeholder="Statut" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Tous les statuts</SelectItem>
                <SelectItem value="airing">En cours</SelectItem>
                <SelectItem value="completed">Terminé</SelectItem>
                <SelectItem value="upcoming">À venir</SelectItem>
                <SelectItem value="hiatus">En pause</SelectItem>
              </SelectContent>
            </Select>
            <Select value={genreFilter} onValueChange={setGenreFilter}>
              <SelectTrigger size="sm" className="w-[140px]">
                <SelectValue placeholder="Genre" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Tous les genres</SelectItem>
                {GENRES.map((g) => (
                  <SelectItem key={g.id} value={g.name}>
                    {g.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="flex items-center gap-2">
            {selected.size > 0 && (
              <Button
                variant="destructive"
                size="sm"
                onClick={() => setDeleteBulkOpen(true)}
              >
                <Trash2 className="size-4" />
                Supprimer ({selected.size})
              </Button>
            )}
            <Button size="sm" onClick={openCreate}>
              <Plus className="size-4" />
              Ajouter
            </Button>
          </div>
        </div>

        {filtered.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-3 py-16 text-center">
            <div className="flex size-10 items-center justify-center rounded-lg bg-muted text-muted-foreground">
              <Clapperboard className="size-5" />
            </div>
            <p className="text-sm font-medium">Aucun animé trouvé</p>
            <p className="text-xs text-muted-foreground">
              Essayez de modifier vos filtres ou ajoutez un nouvel animé.
            </p>
          </div>
        ) : (
          <>
            {/* Table */}
            <div className="mt-4 overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow className="hover:bg-transparent">
                    <TableHead className="w-10">
                      <Checkbox
                        checked={pageItems.length > 0 && selected.size === pageItems.length}
                        onCheckedChange={toggleSelectAll}
                      />
                    </TableHead>
                    <TableHead>Animé</TableHead>
                    <TableHead>
                      <button
                        onClick={() => toggleSort('studio')}
                        className="inline-flex items-center gap-1 hover:text-foreground"
                      >
                        Studio
                        <ArrowUpDown className="size-3" />
                      </button>
                    </TableHead>
                    <TableHead>
                      <button
                        onClick={() => toggleSort('year')}
                        className="inline-flex items-center gap-1 hover:text-foreground"
                      >
                        Année
                        <ArrowUpDown className="size-3" />
                      </button>
                    </TableHead>
                    <TableHead>
                      <button
                        onClick={() => toggleSort('status')}
                        className="inline-flex items-center gap-1 hover:text-foreground"
                      >
                        Statut
                        <ArrowUpDown className="size-3" />
                      </button>
                    </TableHead>
                    <TableHead>Genres</TableHead>
                    <TableHead>
                      <button
                        onClick={() => toggleSort('rating')}
                        className="inline-flex items-center gap-1 hover:text-foreground"
                      >
                        Note
                        <ArrowUpDown className="size-3" />
                      </button>
                    </TableHead>
                    <TableHead className="w-10" />
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {pageItems.map((anime) => (
                    <TableRow
                      key={anime.id}
                      data-state={selected.has(anime.id) ? 'selected' : undefined}
                    >
                      <TableCell>
                        <Checkbox
                          checked={selected.has(anime.id)}
                          onCheckedChange={() => toggleSelect(anime.id)}
                        />
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-3">
                          <div className="flex size-10 shrink-0 items-center justify-center overflow-hidden rounded-md bg-muted">
                            <Image className="size-5 text-muted-foreground" />
                          </div>
                          <div className="min-w-0">
                            <p className="truncate text-sm font-medium">{anime.title}</p>
                            <p className="truncate text-xs text-muted-foreground">
                              {anime.japaneseTitle}
                            </p>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell className="text-muted-foreground">{anime.studio}</TableCell>
                      <TableCell className="text-muted-foreground">{anime.year}</TableCell>
                      <TableCell>
                        <span
                          className={cn(
                            'inline-flex rounded-full px-2 py-0.5 text-xs font-medium',
                            statusStyles[anime.status],
                          )}
                        >
                          {statusLabels[anime.status]}
                        </span>
                      </TableCell>
                      <TableCell>
                        <div className="flex flex-wrap gap-1">
                          {anime.genres.slice(0, 2).map((g) => (
                            <span
                              key={g}
                              className="inline-flex rounded border border-border px-1.5 py-0.5 text-[10px] text-muted-foreground"
                            >
                              {g}
                            </span>
                          ))}
                          {anime.genres.length > 2 && (
                            <span className="inline-flex rounded border border-border px-1.5 py-0.5 text-[10px] text-muted-foreground">
                              +{anime.genres.length - 2}
                            </span>
                          )}
                        </div>
                      </TableCell>
                      <TableCell>
                        <span className="inline-flex items-center gap-1 text-sm font-medium">
                          <Star className="size-3 fill-primary text-primary" />
                          {anime.rating.toFixed(1)}
                        </span>
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <button className="inline-flex size-8 items-center justify-center rounded-md text-muted-foreground hover:bg-accent hover:text-accent-foreground">
                              <MoreHorizontal className="size-4" />
                            </button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => openEdit(anime)}>
                              <Pencil className="size-4" />
                              Modifier
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => handleDuplicate(anime)}>
                              <Copy className="size-4" />
                              Dupliquer
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem
                              variant="destructive"
                              onClick={() => setDeleteItem(anime)}
                            >
                              <Trash2 className="size-4" />
                              Supprimer
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>

            {/* Pagination */}
            <div className="flex items-center justify-between pt-4">
              <p className="text-xs text-muted-foreground">
                {filtered.length} animé{filtered.length > 1 ? 's' : ''} trouvé{filtered.length > 1 ? 's' : ''}
                {selected.size > 0 && ` · ${selected.size} sélectionné${selected.size > 1 ? 's' : ''}`}
              </p>
              {totalPages > 1 && (
                <Pagination>
                  <PaginationContent>
                    <PaginationItem>
                      <PaginationPrevious
                        onClick={() => setPage((p) => Math.max(1, p - 1))}
                        aria-disabled={page === 1}
                        className={cn(page === 1 && 'pointer-events-none opacity-50')}
                      />
                    </PaginationItem>
                    {Array.from({ length: totalPages }, (_, i) => i + 1)
                      .filter((p) => p === 1 || p === totalPages || Math.abs(p - page) <= 1)
                      .reduce<(number | 'ellipsis')[]>((acc, p, i, arr) => {
                        if (i > 0 && p - (arr[i - 1] as number) > 1) acc.push('ellipsis')
                        acc.push(p)
                        return acc
                      }, [])
                      .map((item, i) =>
                        item === 'ellipsis' ? (
                          <PaginationItem key={`e-${i}`}>
                            <PaginationEllipsis />
                          </PaginationItem>
                        ) : (
                          <PaginationItem key={item}>
                            <PaginationLink
                              isActive={item === page}
                              onClick={() => setPage(item)}
                            >
                              {item}
                            </PaginationLink>
                          </PaginationItem>
                        ),
                      )}
                    <PaginationItem>
                      <PaginationNext
                        onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                        aria-disabled={page === totalPages}
                        className={cn(page === totalPages && 'pointer-events-none opacity-50')}
                      />
                    </PaginationItem>
                  </PaginationContent>
                </Pagination>
              )}
            </div>
          </>
        )}
      </section>

      {/* Create / Edit Dialog */}
      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>{editItem ? 'Modifier l\'animé' : 'Ajouter un animé'}</DialogTitle>
            <DialogDescription>
              {editItem
                ? 'Modifiez les informations de cet animé.'
                : 'Remplissez les informations pour ajouter un nouvel animé au catalogue.'}
            </DialogDescription>
          </DialogHeader>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <FormField
                  control={form.control}
                  name="title"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Titre</FormLabel>
                      <FormControl>
                        <Input placeholder="Ex: Eternal Frost" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="japaneseTitle"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Titre japonais</FormLabel>
                      <FormControl>
                        <Input placeholder="永遠の霜" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <div className="grid gap-4 sm:grid-cols-2">
                <FormField
                  control={form.control}
                  name="studio"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Studio</FormLabel>
                      <FormControl>
                        <Input placeholder="Ex: Studio Aurora" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="year"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Année</FormLabel>
                      <FormControl>
                        <Input type="number" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <div className="grid gap-4 sm:grid-cols-3">
                <FormField
                  control={form.control}
                  name="status"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Statut</FormLabel>
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="airing">En cours</SelectItem>
                          <SelectItem value="completed">Terminé</SelectItem>
                          <SelectItem value="upcoming">À venir</SelectItem>
                          <SelectItem value="hiatus">En pause</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="rating"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Note</FormLabel>
                      <FormControl>
                        <Input type="number" step="0.1" min="0" max="10" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="totalEpisodes"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Épisodes</FormLabel>
                      <FormControl>
                        <Input type="number" min="1" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <FormField
                control={form.control}
                name="ageRating"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Classification d&apos;âge</FormLabel>
                    <Select onValueChange={field.onChange} defaultValue={field.value}>
                      <FormControl>
                        <SelectTrigger className="w-[180px]">
                          <SelectValue />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="G">G</SelectItem>
                        <SelectItem value="PG">PG</SelectItem>
                        <SelectItem value="PG-13">PG-13</SelectItem>
                        <SelectItem value="R">R</SelectItem>
                        <SelectItem value="NC-17">NC-17</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="genres"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Genres</FormLabel>
                    <div className="flex flex-wrap gap-2">
                      {GENRES.map((g) => {
                        const checked = field.value.includes(g.name)
                        return (
                          <button
                            key={g.id}
                            type="button"
                            onClick={() => {
                              if (checked) {
                                field.onChange(field.value.filter((v) => v !== g.name))
                              } else {
                                field.onChange([...field.value, g.name])
                              }
                            }}
                            className={cn(
                              'inline-flex items-center rounded-md border px-2.5 py-1 text-xs font-medium transition-colors',
                              checked
                                ? 'border-primary bg-primary/15 text-primary'
                                : 'border-border text-muted-foreground hover:bg-accent hover:text-accent-foreground',
                            )}
                          >
                            {g.name}
                          </button>
                        )
                      })}
                    </div>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="synopsis"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Synopsis</FormLabel>
                    <FormControl>
                      <Textarea placeholder="Description de l'animé…" rows={3} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setCreateOpen(false)}>
                  Annuler
                </Button>
                <Button type="submit">
                  {editItem ? 'Enregistrer' : 'Créer'}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Delete Single Confirmation */}
      <AlertDialog open={!!deleteItem} onOpenChange={() => setDeleteItem(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Supprimer cet animé ?</AlertDialogTitle>
            <AlertDialogDescription>
              Êtes-vous sûr de vouloir supprimer <strong>{deleteItem?.title}</strong> ?
              Cette action est irréversible.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Annuler</AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive text-white hover:bg-destructive/90"
              onClick={() => deleteItem && handleDelete(deleteItem)}
            >
              Supprimer
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Delete Bulk Confirmation */}
      <AlertDialog open={deleteBulkOpen} onOpenChange={setDeleteBulkOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Supprimer la sélection ?</AlertDialogTitle>
            <AlertDialogDescription>
              Êtes-vous sûr de vouloir supprimer <strong>{selected.size} animé{selected.size > 1 ? 's' : ''}</strong> ?
              Cette action est irréversible.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Annuler</AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive text-white hover:bg-destructive/90"
              onClick={handleBulkDelete}
            >
              Tout supprimer
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}

/* ── KPI Card ────────────────────────────────────────────────────────────── */

function KpiCard({
  icon: Icon,
  label,
  value,
  hint,
  accent,
}: {
  icon: React.ElementType
  label: string
  value: number
  hint: string
  accent?: string
}) {
  return (
    <div className="rounded-xl border border-border bg-card p-5">
      <div className="flex items-center justify-between">
        <span className="flex size-10 items-center justify-center rounded-lg bg-accent text-accent-foreground">
          <Icon className="size-5" />
        </span>
        {accent && (
          <span className={cn('size-2 rounded-full', `bg-${accent}`)} />
        )}
      </div>
      <p className="mt-4 text-2xl font-semibold tracking-tight">{value}</p>
      <p className="mt-1 text-sm text-muted-foreground">{label}</p>
      <p className="mt-0.5 text-xs text-muted-foreground/70">{hint}</p>
    </div>
  )
}

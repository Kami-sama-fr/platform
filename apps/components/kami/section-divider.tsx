import { cn } from '@/lib/utils'

const ACCENT_COLORS = [
  'bg-primary/30',
  'bg-gold/25',
  'bg-chart-3/20',
  'bg-primary/20',
  'bg-gold/15',
] as const

let colorIndex = 0

function nextColor() {
  const color = ACCENT_COLORS[colorIndex % ACCENT_COLORS.length]
  colorIndex++
  return color
}

export function SectionDivider({ className }: { className?: string }) {
  const color = nextColor()
  return (
    <div className={cn('mx-4 h-px md:mx-8', className)}>
      <div className={cn('h-full w-full rounded-full', color)} />
    </div>
  )
}

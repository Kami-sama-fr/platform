'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import type { Episode } from '@/types/anime'
import { formatTime } from './format-time'

interface VideoPlayerProps {
  episode: Episode
}

export default function VideoPlayer({ episode }: VideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null)
  const containerRef = useRef<HTMLDivElement>(null)
  const hideTimer = useRef<ReturnType<typeof setTimeout> | null>(null)

  const [playing, setPlaying] = useState(false)
  const [muted, setMuted] = useState(false)
  const [volume, setVolume] = useState(1)
  const [currentTime, setCurrentTime] = useState(0)
  const [duration, setDuration] = useState(0)
  const [buffered, setBuffered] = useState(0)
  const [showControls, setShowControls] = useState(true)
  const [isFullscreen, setIsFullscreen] = useState(false)
  const [isSeeking, setIsSeeking] = useState(false)
  const [seekPreview, setSeekPreview] = useState<number | null>(null)

  const video = videoRef.current

  // Auto-hide controls
  const resetHideTimer = useCallback(() => {
    if (hideTimer.current) clearTimeout(hideTimer.current)
    setShowControls(true)
    if (playing) {
      hideTimer.current = setTimeout(() => setShowControls(false), 3000)
    }
  }, [playing])

  useEffect(() => {
    resetHideTimer()
    return () => {
      if (hideTimer.current) clearTimeout(hideTimer.current)
    }
  }, [playing, resetHideTimer])

  // Sync fullscreen state
  useEffect(() => {
    const onFSChange = () => setIsFullscreen(!!document.fullscreenElement)
    document.addEventListener('fullscreenchange', onFSChange)
    return () => document.removeEventListener('fullscreenchange', onFSChange)
  }, [])

  // Keyboard shortcuts
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (!video) return
      switch (e.key) {
        case ' ':
        case 'k':
          e.preventDefault()
          video.paused ? video.play() : video.pause()
          break
        case 'ArrowLeft':
          e.preventDefault()
          video.currentTime = Math.max(0, video.currentTime - 10)
          break
        case 'ArrowRight':
          e.preventDefault()
          video.currentTime = Math.min(video.duration || 0, video.currentTime + 10)
          break
        case 'ArrowUp':
          e.preventDefault()
          video.volume = Math.min(1, video.volume + 0.1)
          video.muted = false
          break
        case 'ArrowDown':
          e.preventDefault()
          video.volume = Math.max(0, video.volume - 0.1)
          break
        case 'f':
          e.preventDefault()
          toggleFullscreen()
          break
        case 'm':
          e.preventDefault()
          toggleMute()
          break
      }
      resetHideTimer()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [video, resetHideTimer])

  // Video event handlers
  const handleTimeUpdate = () => {
    if (video && !isSeeking) {
      setCurrentTime(video.currentTime)
    }
  }

  const handleProgress = () => {
    if (video && video.buffered.length > 0) {
      setBuffered(video.buffered.end(video.buffered.length - 1))
    }
  }

  const togglePlay = () => {
    if (!video) return
    video.paused ? video.play() : video.pause()
  }

  const toggleMute = () => {
    if (!video) return
    video.muted = !video.muted
    setMuted(video.muted)
  }

  const handleVolumeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!video) return
    const val = parseFloat(e.target.value)
    video.volume = val
    video.muted = val === 0
    setVolume(val)
    setMuted(val === 0)
  }

  const seekBack = () => {
    if (video) video.currentTime = Math.max(0, video.currentTime - 10)
  }

  const seekForward = () => {
    if (video) video.currentTime = Math.min(video.duration || 0, video.currentTime + 10)
  }

  const toggleFullscreen = async () => {
    const el = containerRef.current
    if (!el) return
    if (document.fullscreenElement) {
      await document.exitFullscreen()
    } else {
      await el.requestFullscreen()
    }
  }

  // Progress bar seeking
  const handleProgressMouseDown = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!video || !duration) return
    setIsSeeking(true)
    seekFromEvent(e)

    const onMove = (ev: MouseEvent) => seekFromEvent(ev)
    const onUp = () => {
      setIsSeeking(false)
      document.removeEventListener('mousemove', onMove)
      document.removeEventListener('mouseup', onUp)
    }
    document.addEventListener('mousemove', onMove)
    document.addEventListener('mouseup', onUp)
  }

  const seekFromEvent = (e: MouseEvent | React.MouseEvent) => {
    const bar = (e.target as HTMLElement).closest('[data-progress-bar]') as HTMLElement | null
    if (!bar || !video || !duration) return
    const rect = bar.getBoundingClientRect()
    const pct = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width))
    video.currentTime = pct * duration
    setCurrentTime(pct * duration)
  }

  const handleProgressHover = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!duration) return
    const rect = e.currentTarget.getBoundingClientRect()
    const pct = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width))
    setSeekPreview(pct)
  }

  const progress = duration > 0 ? (currentTime / duration) * 100 : 0
  const bufferedPct = duration > 0 ? (buffered / duration) * 100 : 0

  const VolumeIcon = muted || volume === 0 ? (
    <svg className="size-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5" />
      <line x1="23" y1="9" x2="17" y2="15" />
      <line x1="17" y1="9" x2="23" y2="15" />
    </svg>
  ) : volume < 0.5 ? (
    <svg className="size-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5" />
      <path d="M15.54 8.46a5 5 0 0 1 0 7.07" />
    </svg>
  ) : (
    <svg className="size-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5" />
      <path d="M15.54 8.46a5 5 0 0 1 0 7.07" />
      <path d="M19.07 4.93a10 10 0 0 1 0 14.14" />
    </svg>
  )

  return (
    <div
      ref={containerRef}
      className="group/player relative w-full bg-black select-none"
      onMouseMove={resetHideTimer}
      onMouseLeave={() => playing && setShowControls(false)}
    >
      {/* Video element */}
      <video
        ref={videoRef}
        className="aspect-24/9 w-full bg-black object-contain"
        poster={episode.thumbnail || episode.cover}
        preload="metadata"
        playsInline
        crossOrigin="anonymous"
        onPlay={() => setPlaying(true)}
        onPause={() => setPlaying(false)}
        onTimeUpdate={handleTimeUpdate}
        onDurationChange={(e) => setDuration((e.target as HTMLVideoElement).duration)}
        onProgress={handleProgress}
        onClick={togglePlay}
      >
        <source src={episode.videoUrl} type={episode.videoUrl.endsWith('.m3u8') ? 'application/x-mpegURL' : 'video/mp4'} />
        {episode.tracks.map((track) => (
          <track
            key={track.src}
            src={track.src}
            kind="subtitles"
            srcLang={track.lang}
            label={track.label}
            default={track.default}
          />
        ))}
      </video>

      {/* Big play button when paused */}
      {!playing && (
        <button
          type="button"
          onClick={togglePlay}
          className="absolute inset-0 z-20 flex items-center justify-center bg-black/30 transition-opacity duration-300"
          aria-label="Lire"
        >
          <span className="flex size-16 items-center justify-center rounded-full bg-white/90 text-black shadow-2xl transition-transform duration-300 hover:scale-110">
            <svg className="size-7 ml-1 fill-current" viewBox="0 0 24 24">
              <polygon points="5 3 19 12 5 21 5 3" />
            </svg>
          </span>
        </button>
      )}

      {/* Controls overlay */}
      <div
        className={`absolute inset-0 z-30 flex flex-col justify-end transition-opacity duration-300 ${
          showControls || !playing ? 'opacity-100' : 'opacity-0 pointer-events-none'
        }`}
      >
        {/* Gradient background */}
        <div className="absolute inset-0 bg-linear-to-t from-black/90 via-black/20 to-transparent pointer-events-none" />

        <div className="relative z-10 px-4 pb-3 pt-20 md:px-6">
          {/* Progress bar */}
          <div
            data-progress-bar
            className="group/bar relative h-1 w-full cursor-pointer rounded-full bg-white/20 transition-all duration-200 hover:h-1.5"
            onMouseDown={handleProgressMouseDown}
            onMouseMove={handleProgressHover}
            onMouseLeave={() => setSeekPreview(null)}
          >
            {/* Buffered */}
            <div
              className="absolute inset-y-0 left-0 rounded-full bg-white/20"
              style={{ width: `${bufferedPct}%` }}
            />
            {/* Played */}
            <div
              className="absolute inset-y-0 left-0 rounded-full bg-[#e50914] transition-[width] duration-75"
              style={{ width: `${progress}%` }}
            />
            {/* Seek preview */}
            {seekPreview !== null && (
              <div
                className="absolute inset-y-0 left-0 rounded-full bg-white/10"
                style={{ width: `${seekPreview * 100}%` }}
              />
            )}
            {/* Scrubber dot */}
            <div
              className="absolute top-1/2 size-3 -translate-x-1/2 -translate-y-1/2 rounded-full bg-[#e50914] opacity-0 transition-opacity duration-200 group-hover/bar:opacity-100"
              style={{ left: `${progress}%` }}
            />
          </div>

          {/* Control row */}
          <div className="mt-2 flex items-center justify-between gap-2">
            {/* Left: Play, skip, volume */}
            <div className="flex items-center gap-1">
              {/* Play / Pause */}
              <button type="button" onClick={togglePlay} className="p-1.5 text-white/90 hover:text-white transition-colors rounded" aria-label={playing ? 'Pause' : 'Lire'}>
                {playing ? (
                  <svg className="size-5" viewBox="0 0 24 24" fill="currentColor">
                    <rect x="6" y="4" width="4" height="16" rx="1" />
                    <rect x="14" y="4" width="4" height="16" rx="1" />
                  </svg>
                ) : (
                  <svg className="size-5 ml-0.5" viewBox="0 0 24 24" fill="currentColor">
                    <polygon points="5 3 19 12 5 21 5 3" />
                  </svg>
                )}
              </button>

              {/* Skip back 10s */}
              <button type="button" onClick={seekBack} className="p-1.5 text-white/90 hover:text-white transition-colors rounded" aria-label="Reculer 10 secondes">
                <svg className="size-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M1 4v6h6" />
                  <path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10" />
                </svg>
              </button>

              {/* Skip forward 10s */}
              <button type="button" onClick={seekForward} className="p-1.5 text-white/90 hover:text-white transition-colors rounded" aria-label="Avancer 10 secondes">
                <svg className="size-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M23 4v6h-6" />
                  <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10" />
                </svg>
              </button>

              {/* Volume */}
              <div className="flex items-center gap-1 group/vol">
                <button type="button" onClick={toggleMute} className="p-1.5 text-white/90 hover:text-white transition-colors rounded" aria-label={muted ? 'Activer le son' : 'Couper le son'}>
                  {VolumeIcon}
                </button>
                <input
                  type="range"
                  min={0}
                  max={1}
                  step={0.01}
                  value={muted ? 0 : volume}
                  onChange={handleVolumeChange}
                  className="h-1 w-0 cursor-pointer appearance-none rounded-full bg-white/30 accent-white transition-all duration-200 group-hover/vol:w-20 [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:size-2.5 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:cursor-pointer"
                  aria-label="Volume"
                />
              </div>
            </div>

            {/* Center: Episode title */}
            <span className="hidden text-sm font-medium text-white/80 md:inline truncate max-w-xs px-4">
              EP {episode.number} — {episode.title}
            </span>

            {/* Right: Time, CC, fullscreen */}
            <div className="flex items-center gap-1">
              {/* Time */}
              <span className="text-xs tabular-nums text-white/60 mr-2 hidden sm:inline">
                {formatTime(currentTime)} / {formatTime(duration)}
              </span>

              {/* Subtitles */}
              {episode.tracks.length > 0 && (
                <button type="button" className="p-1.5 text-white/90 hover:text-white transition-colors rounded" aria-label="Sous-titres">
                  <svg className="size-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <rect x="1" y="4" width="22" height="16" rx="2" />
                    <path d="M7 12h2" />
                    <path d="M11 12h4" />
                    <path d="M7 16h10" />
                  </svg>
                </button>
              )}

              {/* Fullscreen */}
              <button type="button" onClick={toggleFullscreen} className="p-1.5 text-white/90 hover:text-white transition-colors rounded" aria-label={isFullscreen ? 'Quitter le plein écran' : 'Plein écran'}>
                {isFullscreen ? (
                  <svg className="size-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M8 3v3a2 2 0 0 1-2 2H3" />
                    <path d="M21 8h-3a2 2 0 0 1-2-2V3" />
                    <path d="M3 16h3a2 2 0 0 1 2 2v3" />
                    <path d="M16 21v-3a2 2 0 0 1 2-2h3" />
                  </svg>
                ) : (
                  <svg className="size-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M8 3H5a2 2 0 0 0-2 2v3" />
                    <path d="M21 8V5a2 2 0 0 0-2-2h-3" />
                    <path d="M3 16v3a2 2 0 0 0 2 2h3" />
                    <path d="M16 21h3a2 2 0 0 0 2-2v-3" />
                  </svg>
                )}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

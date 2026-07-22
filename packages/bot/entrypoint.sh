#!/bin/sh
set -e

# ══════════════════════════════════════════════════════════════════════════════
# Kami-Sama Bot — Development Entrypoint
# ══════════════════════════════════════════════════════════════════════════════

export PATH="/usr/local/bin:/usr/bin:/bin:${PATH}"
export NODE_ENV="${NODE_ENV:-development}"

# ── Logging ───────────────────────────────────────────────────────────────────
log_info() {
    echo "[INFO] $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_warn() {
    echo "[WARN] $(date '+%Y-%m-%d %H:%M:%S') - $1" >&2
}

log_error() {
    echo "[ERROR] $(date '+%Y-%m-%d %H:%M:%S') - $1" >&2
}

# ── Configuration de la base de données ───────────────────────────────────────
configure_database() {
    if [ -n "${POSTGRESQL__HOST:-}" ] && [ -z "${DB_HOST:-}" ]; then
        export DB_HOST="${POSTGRESQL__HOST}"
    fi
    if [ -n "${POSTGRESQL__PORT:-}" ] && [ -z "${DB_PORT:-}" ]; then
        export DB_PORT="${POSTGRESQL__PORT}"
    fi
    if [ -n "${POSTGRESQL__USER:-}" ] && [ -z "${DB_USER:-}" ]; then
        export DB_USER="${POSTGRESQL__USER}"
    fi
    if [ -n "${POSTGRESQL__NAME:-}" ] && [ -z "${DB_NAME:-}" ]; then
        export DB_NAME="${POSTGRESQL__NAME}"
    fi
    if [ -n "${POSTGRESQL__PASSWORD:-}" ] && [ -z "${DB_PASSWORD:-}" ]; then
        export DB_PASSWORD="${POSTGRESQL__PASSWORD}"
    fi

    if [ -n "${PG_HOST:-}" ] && [ -z "${DB_HOST:-}" ]; then
        export DB_HOST="${PG_HOST}"
    fi
    if [ -n "${PG_PORT:-}" ] && [ -z "${DB_PORT:-}" ]; then
        export DB_PORT="${PG_PORT}"
    fi
    if [ -n "${PG_USER:-}" ] && [ -z "${DB_USER:-}" ]; then
        export DB_USER="${PG_USER}"
    fi
    if [ -n "${PG_DB:-}" ] && [ -z "${DB_NAME:-}" ]; then
        export DB_NAME="${PG_DB}"
    fi
    if [ -n "${PG_PASS:-}" ] && [ -z "${DB_PASSWORD:-}" ]; then
        export DB_PASSWORD="${PG_PASS}"
    fi

    export DB_HOST="${DB_HOST:-postgresql}"
    export DB_PORT="${DB_PORT:-5432}"
    export DB_USER="${DB_USER:-postgres}"
    export DB_NAME="${DB_NAME:-postgres}"
    export DB_PASSWORD="${DB_PASSWORD:-${POSTGRES_PASSWORD:-postgres}}"

    if [ -z "${DATABASE_URL:-}" ]; then
        export DATABASE_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
    fi
}

# ── Configuration Redis ───────────────────────────────────────────────────────
configure_redis() {
    if [ -n "${REDIS_URL:-}" ]; then
        redis_url="${REDIS_URL#redis://}"
        redis_url="${redis_url#rediss://}"
        redis_authority="${redis_url%%/*}"
        redis_db="${redis_url#*/}"
        redis_db="${redis_db%%\?*}"

        credentials=""
        redis_host_port="${redis_authority}"
        if [ "${redis_authority#*@}" != "${redis_authority}" ]; then
            credentials="${redis_authority%@*}"
            redis_host_port="${redis_authority#*@}"
        fi

        if [ -n "${credentials}" ]; then
            case "${credentials}" in
                *:*)
                    export REDIS_PASSWORD="${credentials#*:}"
                    ;;
                *)
                    export REDIS_PASSWORD="${credentials}"
                    ;;
            esac
        fi

        case "${redis_host_port}" in
            \[*\]:*)
                export REDIS_HOST="${redis_host_port%\]:*}"
                export REDIS_HOST="${REDIS_HOST#\[}"
                export REDIS_PORT="${redis_host_port##*\]:}"
                ;;
            *:*)
                export REDIS_HOST="${redis_host_port%%:*}"
                export REDIS_PORT="${redis_host_port#*:}"
                ;;
            *)
                export REDIS_HOST="${redis_host_port}"
                ;;
        esac

        if [ -n "${redis_db}" ] && [ "${redis_db}" != "${redis_url}" ]; then
            export REDIS_DB="${redis_db}"
        fi
    fi

    export REDIS_HOST="${REDIS_HOST:-redis}"
    export REDIS_PORT="${REDIS_PORT:-6379}"
    export REDIS_DB="${REDIS_DB:-0}"
    export REDIS_ENABLED="${REDIS_ENABLED:-true}"
    export REDIS_REQUIRED="${REDIS_REQUIRED:-false}"
    export REDIS_KEY_PREFIX="${REDIS_KEY_PREFIX:-kami-bot:v1}"
}

# ── Configuration du runtime ──────────────────────────────────────────────────
configure_runtime() {
    configure_database
    configure_redis

    if [ -n "${SECRET_KEY:-}" ] && [ -z "${SYSTEM_KEY:-}" ]; then
        export SYSTEM_KEY="${SECRET_KEY}"
    fi
}

# ── Prisma ────────────────────────────────────────────────────────────────────
find_prisma_bin() {
    for bin in ./prisma/node_modules/.bin/prisma ./node_modules/.bin/prisma; do
        if [ -x "${bin}" ]; then
            echo "${bin}"
            return 0
        fi
    done

    if command -v npx >/dev/null 2>&1; then
        echo "npx prisma"
        return 0
    fi

    return 1
}

run_prisma_generate() {
    if [ ! -f ./prisma/schema.prisma ]; then
        log_warn "Prisma schema not found at ./prisma/schema.prisma; skipping"
        return 0
    fi

    prisma_bin="$(find_prisma_bin || true)"

    if [ -z "${prisma_bin}" ]; then
        log_warn "Prisma CLI is not available; skipping"
        return 0
    fi

    log_info "Generating Prisma client..."
    if DATABASE_URL="${DATABASE_URL}" ${prisma_bin} generate; then
        log_info "Prisma client generated"
    else
        log_warn "Prisma generate failed; continuing"
    fi
}

# ── Lancement du bot ──────────────────────────────────────────────────────────
run_bot() {
    log_info "Kami-Sama Bot starting (development mode)"

    if [ ! -f /app/index.js ]; then
        log_error "Bot entrypoint not found at /app/index.js"
        return 1
    fi

    if [ ! -d /app/node_modules ]; then
        log_error "Bot dependencies not found at /app/node_modules"
        return 1
    fi

    exec node --watch /app/index.js
}

# ── Point d'entrée ────────────────────────────────────────────────────────────
configure_runtime
run_prisma_generate

role="${1:-bot}"

case "${role}" in
    bot)
        shift || true
        run_bot "$@"
        ;;
    *)
        exec "$@"
        ;;
esac

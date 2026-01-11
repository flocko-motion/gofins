# Bugs (fix bugs first before working on TODOs)

## Weekly YoY data has N/A blocks in 2020-2021
Weekly prices show blocks of N/A YoY values during 2020-2021 period (e.g., AAPL from July 2020 to April 2021). Pattern shows every other week has N/A, suggesting missing weekly price data or calculation issue for that period. Need to investigate why YoY calculation fails for these specific weeks.

# TODO list for gofins, in order of priority

## FMP - IMPLEMENTED ✓

**Implementation completed:**
- Dual currency storage: prices stored in both original currency (*_orig columns) and USD
- Update threshold changed from monthly (1st of month) to 30-day rolling window
- Original prices preserved in DB, allowing reconversion without refetching from FMP
- Migration: deployment/migrations/001_add_original_currency_prices.sql

**To deploy:**
1. Run migration: `./gofins db sql -q "$(cat deployment/migrations/001_add_original_currency_prices.sql)"`
2. Next price update will populate original currency columns
3. Monitor FMP query reduction (should be ~12x fewer queries per year)

**Future enhancements (optional):**
- Daily watchlist-only price updates (not yet implemented)
- Forex conversion still uses monthly cached data (already optimized)

## Journal / Notebook Feature

**Concept**: Unified research journal mixing ratings and freeform notes

**Database Design**:
- Keep `user_ratings` table as-is (optimized for ticker-specific queries)
- Add new `user_journal` table for freeform notes:
  - id, title, content, type ('note'|'idea'|'news'|'strategy'), tags (JSON), created_at, updated_at
- Add `journal_tickers` junction table (many-to-many):
  - journal_id, ticker
  - Allows notes to link to 0, 1, or multiple tickers

**UI/UX**:
- New "Journal" tab showing unified timeline of ratings + notes
- Each entry shows type icon (⭐ rating, 📝 note, 💡 idea, 📰 news)
- Ratings displayed as special note type with numeric badge
- Linked tickers shown as clickable badges
- Filter by: type, ticker, tag, date range
- Full-text search across all content
- Quick capture: floating "+" button or hotkey

**Features**:
- Create standalone notes (research, ideas, strategies)
- Link notes to multiple tickers
- Tag/categorize entries
- View all entries related to a ticker (in symbol detail view)
- Export journal to markdown/PDF
- Auto-link ticker mentions in text (e.g., $AAPL)

**Implementation Phases**:
1. Phase 1: Basic notes (create, edit, delete, link to tickers)
2. Phase 2: Tags, categories, rich text/markdown editor
3. Phase 3: Unified timeline view mixing ratings + notes
4. Phase 4: Advanced features (attachments, export, auto-linking)

**API Endpoints**:
- GET/POST /api/journal - List/create entries
- GET/PUT/DELETE /api/journal/{id} - CRUD operations
- GET /api/journal/ticker/{ticker} - All entries mentioning ticker
- POST/DELETE /api/journal/{id}/tickers - Link/unlink tickers
- GET /api/timeline - Unified view of ratings + journal entries

## Beta correlation

- beta correlation should be added to analysis module (can't be in profile itself, because we need to specify a time range and a reference index)
- we should offer a dropdown of a few handselected indices as reference index: NXP, SPX .. maybe a few more
- beta should be added as another column
- the score formula needs beta as third ingredient..so we need three sliders, one for each weight: stddev, mü, beta

## Backup Strategy

**What to back up**: Only user-generated data (everything else is reproducible from market data)

**User tables to backup**:
- `user_ratings` - Stock ratings and notes
- `user_favorites` - Favorite stocks
- `analysis_packages` - Custom analysis configurations
- `users` - User accounts (if multi-user)

**Strategy **:

### Option 1: PostgreSQL pg_dump (Recommended)
```bash
# Daily backup script
#!/bin/bash
BACKUP_DIR="/opt/backups/gofins"
DATE=$(date +%Y%m%d_%H%M%S)
docker exec gofins-db pg_dump -U gofins -d gofins \
  --table=user_ratings \
  --table=user_favorites \
  --table=analysis_packages \
  --table=users \
  > "$BACKUP_DIR/user_data_$DATE.sql"

# Keep last 30 days, compress older
find "$BACKUP_DIR" -name "*.sql" -mtime +7 -exec gzip {} \;
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +30 -delete
```

### Option 2: Docker Volume Backup
```bash
# Backup entire PostgreSQL data directory
docker run --rm \
  --volumes-from gofins-db \
  -v /opt/backups:/backup \
  alpine tar czf /backup/gofins-db-$(date +%Y%m%d).tar.gz /var/lib/postgresql/data
```

### Option 3: Export to JSON (Application-level)
- Add `/api/backup/export` endpoint that returns JSON of all user data
- Add `/api/backup/import` endpoint to restore from JSON
- Pros: Database-agnostic, human-readable, easy to version control
- Cons: Requires application code, slower for large datasets

**Recommended Setup**:
1. Daily pg_dump of user tables (Option 1)
2. Weekly full database backup (Option 2)
3. Cron job: `0 2 * * * /opt/gofins/backup.sh`
4. Store backups on separate volume/server
5. Test restore procedure monthly

**Restore procedure**:
```bash
# Restore from SQL dump
docker exec -i gofins-db psql -U gofins -d gofins < user_data_backup.sql

# Or restore full volume
docker stop gofins-db
tar xzf gofins-db-backup.tar.gz -C /var/lib/docker/volumes/gofins-db/_data
docker start gofins-db
```

## Deployment Setup

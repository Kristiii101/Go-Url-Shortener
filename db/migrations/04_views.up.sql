-- Convenience views so your API queries are simple and consistent

-- All-time totals per link (clicks, unique visitors, last clicked timestamp)
CREATE OR REPLACE VIEW v_link_stats_totals AS
SELECT
  l.id AS link_id,
  l.key,
  COUNT(c.id) AS clicks_total,
  COUNT(DISTINCT c.visitor_hash) AS unique_visitors_total,
  MAX(c.occurred_at) AS last_clicked_at
FROM links l
LEFT JOIN clicks c ON c.link_id = l.id
GROUP BY l.id, l.key;

-- Daily totals (UTC) with unique visitors
CREATE OR REPLACE VIEW v_link_stats_daily AS
SELECT
  c.link_id,
  date(timezone('UTC', c.occurred_at)) AS day,
  COUNT(*) AS clicks,
  COUNT(DISTINCT c.visitor_hash) AS unique_visitors
FROM clicks c
GROUP BY c.link_id, date(timezone('UTC', c.occurred_at));

-- Regional (country) totals (optional)
CREATE OR REPLACE VIEW v_link_stats_regions AS
SELECT
  c.link_id,
  COALESCE(c.country_code, 'ZZ') AS country_code,
  COUNT(*) AS clicks,
  COUNT(DISTINCT c.visitor_hash) AS unique_visitors
FROM clicks c
GROUP BY c.link_id, COALESCE(c.country_code, 'ZZ');
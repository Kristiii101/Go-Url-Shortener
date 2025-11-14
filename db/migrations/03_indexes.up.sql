-- Expression index to speed up daily aggregations (UTC day)
CREATE INDEX IF NOT EXISTS idx_clicks_link_day
  ON clicks (link_id, (date(timezone('UTC', occurred_at))));

-- Unique visitor aggregations benefit from this
CREATE INDEX IF NOT EXISTS idx_clicks_link_visitor
  ON clicks (link_id, visitor_hash);

-- Regional breakdowns
CREATE INDEX IF NOT EXISTS idx_clicks_link_country
  ON clicks (link_id, country_code);
WITH matching_timestamps AS (
  SELECT 
    ph.id AS play_id,
    ph.release_id,
    ph.played_at,
    ch.id AS cleaning_id,
    ch.cleaned_at
  FROM play_history ph
  JOIN cleaning_history ch ON ph.release_id = ch.release_id AND ph.played_at = ch.cleaned_at
)

UPDATE play_history
SET played_at = datetime(played_at, '+1 seconds')
WHERE id IN (SELECT play_id FROM matching_timestamps);

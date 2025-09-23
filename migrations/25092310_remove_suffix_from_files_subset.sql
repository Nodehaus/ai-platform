-- Remove _eng suffix from files_subset column values
-- This updates all filenames in the files_subset array to remove the _eng suffix

UPDATE corpus
SET files_subset = array(
    SELECT regexp_replace(unnest(files_subset), '_eng\.json$', '.json', 'g')
)
WHERE files_subset IS NOT NULL;
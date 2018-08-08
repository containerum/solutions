CREATE OR REPLACE FUNCTION pc_chartoint(chartoconvert character varying)
  RETURNS integer AS
$BODY$
SELECT CASE WHEN regexp_replace($1, '[^0-9]+', '', 'g') SIMILAR TO '[0-9]+'
                 THEN CAST(regexp_replace($1, '[^0-9]+', '', 'g') AS integer)
            ELSE 100 END;

$BODY$
LANGUAGE 'sql' IMMUTABLE STRICT;

ALTER TABLE templates
  ALTER COLUMN cpu TYPE integer USING pc_chartoint(cpu),
  ALTER COLUMN ram TYPE integer USING pc_chartoint(ram);

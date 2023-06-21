CREATE INDEX fs_name_idx ON users(first_name, second_name text_pattern_ops) ;

CREATE EXTENSION pg_trgm;
UPDATE pg_opclass SET opcdefault = TRUE WHERE opcname='gin_trgm_ops';
CREATE INDEX CONCURRENTLY fs_name_gin_idx ON users USING gin (first_name, second_name);
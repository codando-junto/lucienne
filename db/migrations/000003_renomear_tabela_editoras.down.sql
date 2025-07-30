ALTER TABLE publishers DROP CONSTRAINT publishers_name_key;
ALTER TABLE publishers RENAME TO editoras;
ALTER TABLE editoras RENAME COLUMN id TO editora_id;
ALTER TABLE editoras RENAME COLUMN name TO nome;

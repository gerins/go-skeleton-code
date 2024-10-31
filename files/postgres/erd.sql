SET TIMEZONE = 'Etc/GMT-7';

CREATE OR REPLACE FUNCTION update_modified_column() 
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW; 
END;
$$ language 'plpgsql';

---------------------------------------------------------------------------------------------------------------------

CREATE TABLE fuels (
    id                              SERIAL PRIMARY KEY,
    type                            VARCHAR(255) DEFAULT '',
    description                     VARCHAR(255) DEFAULT '',
    unit                            VARCHAR(255) DEFAULT '',
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at                      TIMESTAMP WITH TIME ZONE 
);

CREATE TRIGGER fuels BEFORE UPDATE ON fuels FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO fuels (
    "id",
    "deleted_at"
) VALUES (0, now());

---------------------------------------------------------------------------------------------------------------------

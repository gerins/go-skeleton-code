SET TIMEZONE = 'Etc/GMT-7';

CREATE OR REPLACE FUNCTION update_modified_column() 
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW; 
END;
$$ language 'plpgsql';

---------------------------------------------------------------------------------------------------------------------

CREATE TABLE users (
    id                              SERIAL PRIMARY KEY,
    full_name                       VARCHAR(128) NOT NULL DEFAULT '',
    email                           VARCHAR(128) NOT NULL,
    phone_number                    VARCHAR(128) NOT NULL DEFAULT '',
    password                        VARCHAR(512) NOT NULL,
    status                          BOOLEAN NOT NULL DEFAULT true,
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at                      TIMESTAMP WITH TIME ZONE 
);

CREATE TRIGGER users BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO users (
    "id",
    "full_name",
    "email",
    "phone_number",
    "password",
    "status",
    "deleted_at"
) VALUES (0, '', '', '', '', false, now());

---------------------------------------------------------------------------------------------------------------------

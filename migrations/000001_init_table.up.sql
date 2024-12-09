-- Note: This was managed by GORM auto migration, but it is included here for reference.
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ, -- Soft delete
    service_name VARCHAR(255) NOT NULL,
    service_description TEXT
);

CREATE TABLE IF NOT EXISTS service_versions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ, -- Soft delete
    service_id INT NOT NULL REFERENCES services(id) ON UPDATE CASCADE ON DELETE CASCADE,
    service_version_name VARCHAR(255),
    service_version_url TEXT,
    service_version_description TEXT
);

-- Validate if the existing services match the schema
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'services') THEN
        RAISE EXCEPTION 'Table "services" does not exist';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'services' AND column_name = 'service_name') THEN
        RAISE EXCEPTION 'Column "service_name" does not exist in table "services"';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'service_versions') THEN
        RAISE EXCEPTION 'Table "service_versions" does not exist';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'service_versions' AND column_name = 'service_version_name') THEN
        RAISE EXCEPTION 'Column "service_version_name" does not exist in table "service_versions"';
    END IF;
END $$;
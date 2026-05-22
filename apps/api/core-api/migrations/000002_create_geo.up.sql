-- Гео-таблицы (PostGIS): районы ЗКО и земельные участки
-- Создаётся до farms/borrowers, т.к. farms ссылается на districts через поле district

CREATE TABLE IF NOT EXISTS districts (
    id      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name    VARCHAR(255) NOT NULL,
    name_kz VARCHAR(255),
    geom    GEOMETRY(MULTIPOLYGON, 4326)
);

CREATE INDEX IF NOT EXISTS idx_districts_geom ON districts USING GIST (geom);

-- Фермерские хозяйства с геопривязкой
CREATE TABLE IF NOT EXISTS farms (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255) NOT NULL,
    district      VARCHAR(255),
    activity_type VARCHAR(50)  CHECK (activity_type IN ('crop_farming', 'livestock', 'mixed')),
    land_area_ha  NUMERIC(10, 2),
    location      GEOMETRY(POINT, 4326),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_farms_location ON farms USING GIST (location);

-- Земельные участки — кадастровые полигоны
CREATE TABLE IF NOT EXISTS land_parcels (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    farm_id          UUID        NOT NULL REFERENCES farms(id) ON DELETE CASCADE,
    land_type        VARCHAR(50) CHECK (land_type IN ('cropland', 'pasture', 'fallow')),
    area_ha          NUMERIC(10, 2),
    cadastral_number VARCHAR(255) UNIQUE,
    geom             GEOMETRY(POLYGON, 4326)
);

CREATE INDEX IF NOT EXISTS idx_land_parcels_farm_id ON land_parcels (farm_id);
CREATE INDEX IF NOT EXISTS idx_land_parcels_geom    ON land_parcels USING GIST (geom);

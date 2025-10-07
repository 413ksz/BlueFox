CREATE TABLE "user" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(30) NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(70),
    last_name VARCHAR(70),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    last_online TIMESTAMP WITH TIME ZONE,
    bio VARCHAR(255),
    date_of_birth DATE NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    profile_picture_asset_id UUID,
    FOREIGN KEY (profile_picture_asset_id) REFERENCES media_asset(id)
);

CREATE INDEX idx_username ON "user" (username);

ALTER TABLE media_asset
ADD CONSTRAINT fk_uploaded_by_user
FOREIGN KEY (uploaded_by_user_id) REFERENCES "user"(id);
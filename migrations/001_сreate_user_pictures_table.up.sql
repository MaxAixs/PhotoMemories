CREATE TABLE user_pictures (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    pic_key VARCHAR(255) NOT NULL,
    tag VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uniq_user_tag UNIQUE (user_id, tag)
);


CREATE INDEX idx_user_id ON user_pictures(user_id);
CREATE INDEX idx_tag ON user_pictures(tag);

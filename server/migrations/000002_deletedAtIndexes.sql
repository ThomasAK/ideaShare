-- +goose Up

CREATE INDEX idx_users_deleted_at ON users (deleted_at);
CREATE INDEX idx_user_roles_deleted_at ON user_roles (deleted_at);
CREATE INDEX idx_ideas_deleted_at ON ideas (deleted_at);
CREATE INDEX idx_idea_comments_deleted_at ON idea_comments (deleted_at);
CREATE INDEX idx_idea_likes_deleted_at ON idea_likes (deleted_at);

-- +goose Down

DROP TABLE idea_likes;
DROP TABLE idea_comments;
DROP TABLE ideas;
DROP TABLE user_roles;
DROP TABLE users;
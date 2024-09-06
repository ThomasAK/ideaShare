-- +goose Up
CREATE TABLE users
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    external_id VARCHAR(255)                        NOT NULL,
    first_name  VARCHAR(50),
    last_name   VARCHAR(50),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP DEFAULT NULL
);
CREATE TABLE user_roles
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT                                 NOT NULL,
    role       VARCHAR(20)                         NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    UNIQUE KEY (user_id, role)
);
CREATE TABLE ideas
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    title       VARCHAR(255)                        NOT NULL,
    description TEXT,
    status      VARCHAR(20)                         NOT NULL,
    created_by  INT                                 NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (created_by) REFERENCES users (id)
);
CREATE TABLE idea_comments
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    idea_id    INT                                 NOT NULL,
    comment    TEXT                                NOT NULL check ( length(comment) > 3 ),
    created_by INT                                 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (idea_id) REFERENCES ideas (id),
    FOREIGN KEY (created_by) REFERENCES users (id)
);
CREATE TABLE idea_likes
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    idea_id    INT                                 NOT NULL,
    created_by INT                                 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (idea_id) REFERENCES ideas (id),
    FOREIGN KEY (created_by) REFERENCES users (id),
    UNIQUE KEY (idea_id, created_by)
);

CREATE INDEX idx_idea_status ON ideas (status);
CREATE INDEX idx_idea_created_by ON ideas (created_by);

CREATE INDEX idx_user_external_id ON users (external_id);

CREATE INDEX idx_user_role_user_id ON user_roles (user_id);

CREATE INDEX idx_idea_comment_idea_id ON idea_comments (idea_id);

CREATE INDEX idx_idea_like_idea_id ON idea_likes (idea_id);
CREATE INDEX idx_idea_like_created_by ON idea_likes (created_by);

INSERT INTO users (external_id, first_name, last_name)
VALUES ('12354', 'admin', 'admin');
INSERT INTO user_roles (user_id, role)
VALUES (1, 'ADMIN');


-- +goose Down
DROP TABLE idea_likes;
DROP TABLE idea_comments;
DROP TABLE ideas;
DROP TABLE user_roles;
DROP TABLE users;
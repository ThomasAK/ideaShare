-- +goose Up
CREATE TABLE user (
                      id INT AUTO_INCREMENT PRIMARY KEY,
                      external_id VARCHAR(255),
                      first_name VARCHAR(50),
                      last_name VARCHAR(50)
);
CREATE TABLE user_role (
                           id INT AUTO_INCREMENT PRIMARY KEY,
                           user_id INT NOT NULL,
                           role VARCHAR(20),
                           FOREIGN KEY (user_id) REFERENCES user(id)
);
CREATE TABLE idea (
                      id INT AUTO_INCREMENT PRIMARY KEY,
                      title VARCHAR(255) NOT NULL,
                      description TEXT,
                      status VARCHAR(20),
                      created_by INT,
                      created_ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                      updated_ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                      FOREIGN KEY (created_by) REFERENCES user(id)
);
CREATE TABLE idea_comment (
                              id INT AUTO_INCREMENT PRIMARY KEY,
                              idea_id INT NOT NULL,
                              comment TEXT,
                              created_by INT,
                              created_ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              updated_ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                              FOREIGN KEY (idea_id) REFERENCES idea(id),
                              FOREIGN KEY (created_by) REFERENCES user(id)
);
CREATE TABLE idea_like (
                           id INT AUTO_INCREMENT PRIMARY KEY,
                           idea_id INT NOT NULL,
                           created_by INT,
                           created_ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                           FOREIGN KEY (idea_id) REFERENCES idea(id),
                           FOREIGN KEY (created_by) REFERENCES user(id),
                           UNIQUE KEY (idea_id, created_by)
);

CREATE INDEX idx_idea_status ON idea (status);
CREATE INDEX idx_idea_created_by ON idea (created_by);

CREATE INDEX idx_user_external_id ON user (external_id);

CREATE INDEX idx_user_role_user_id ON user_role (user_id);

CREATE INDEX idx_idea_comment_idea_id ON idea_comment (idea_id);

CREATE INDEX idx_idea_like_idea_id ON idea_like (idea_id);
CREATE INDEX idx_idea_like_created_by ON idea_like (created_by);


-- +goose Down
DROP TABLE idea_like;
DROP TABLE idea_comment;
DROP TABLE idea;
DROP TABLE user_role;
DROP TABLE user;
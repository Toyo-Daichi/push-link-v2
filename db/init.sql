DROP DATABASE IF EXISTS push_link_v2;
CREATE DATABASE push_link_v2
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE push_link_v2;

CREATE TABLE users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(255) NOT NULL,
  role ENUM('admin', 'editor', 'viewer') NOT NULL DEFAULT 'viewer',
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  create_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE tags (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) NOT NULL,
  create_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_tags_name (name),
  UNIQUE KEY uk_tags_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE sites (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  url VARCHAR(2048) NOT NULL,
  domain VARCHAR(255) NOT NULL,
  status ENUM('draft', 'published', 'archived') NOT NULL DEFAULT 'published',
  added_by_user_id BIGINT UNSIGNED NOT NULL,
  create_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_sites_url (url),
  KEY idx_sites_domain (domain),
  KEY idx_sites_status (status),
  KEY idx_sites_added_by_user_id (added_by_user_id),
  CONSTRAINT fk_sites_added_by_user_id
    FOREIGN KEY (added_by_user_id) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE site_tags (
  site_id BIGINT UNSIGNED NOT NULL,
  tag_id BIGINT UNSIGNED NOT NULL,
  create_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (site_id, tag_id),
  KEY idx_site_tags_tag_id (tag_id),
  CONSTRAINT fk_site_tags_site_id
    FOREIGN KEY (site_id) REFERENCES sites (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_site_tags_tag_id
    FOREIGN KEY (tag_id) REFERENCES tags (id)
    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE user_site_bookmarks (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  site_id BIGINT UNSIGNED NOT NULL,
  note VARCHAR(255) NULL,
  create_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_user_site_bookmarks_user_site (user_id, site_id),
  KEY idx_user_site_bookmarks_site_id (site_id),
  CONSTRAINT fk_user_site_bookmarks_user_id
    FOREIGN KEY (user_id) REFERENCES users (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_user_site_bookmarks_site_id
    FOREIGN KEY (site_id) REFERENCES sites (id)
    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO users (id, name, email, role, is_active) VALUES
  (1, 'Toyo Admin', 'admin@example.com', 'admin', 1),
  (2, 'Mika Curator', 'mika@example.com', 'editor', 1),
  (3, 'Ken Viewer', 'ken@example.com', 'viewer', 1);

INSERT INTO tags (id, name, slug) VALUES
  (1, 'AI', 'ai'),
  (2, 'Development', 'development'),
  (3, 'Design', 'design'),
  (4, 'Learning', 'learning'),
  (5, 'Productivity', 'productivity');

INSERT INTO sites (id, title, description, url, domain, status, added_by_user_id) VALUES
  (
    1,
    'MDN Web Docs',
    'Web standards, HTML, CSS, JavaScript and API references for developers.',
    'https://developer.mozilla.org/',
    'developer.mozilla.org',
    'published',
    1
  ),
  (
    2,
    'Stack Overflow',
    'Q&A platform for software development and related technical topics.',
    'https://stackoverflow.com/',
    'stackoverflow.com',
    'published',
    2
  ),
  (
    3,
    'Figma',
    'Collaborative interface design and prototyping platform.',
    'https://www.figma.com/',
    'www.figma.com',
    'published',
    2
  ),
  (
    4,
    'freeCodeCamp',
    'Hands-on programming tutorials and learning resources.',
    'https://www.freecodecamp.org/',
    'www.freecodecamp.org',
    'published',
    1
  ),
  (
    5,
    'OpenAI Developers',
    'API documentation, guides and examples for building AI products.',
    'https://platform.openai.com/docs/overview',
    'platform.openai.com',
    'draft',
    1
  );

INSERT INTO site_tags (site_id, tag_id) VALUES
  (1, 2),
  (1, 4),
  (2, 2),
  (2, 5),
  (3, 3),
  (3, 5),
  (4, 2),
  (4, 4),
  (5, 1),
  (5, 2);

INSERT INTO user_site_bookmarks (id, user_id, site_id, note) VALUES
  (1, 1, 1, 'Reference docs used frequently'),
  (2, 1, 5, 'Review before publishing'),
  (3, 2, 3, 'Useful for UI inspiration'),
  (4, 3, 4, 'Good learning path for beginners');

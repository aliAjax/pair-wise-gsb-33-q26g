-- gbplantwiki 植物养护知识百科平台 初始化脚本（MySQL 8.0）
-- 由 MySQL 官方镜像 /docker-entrypoint-initdb.d 首次启动时自动执行。

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(64) NOT NULL UNIQUE,
  email VARCHAR(128) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  nickname VARCHAR(64),
  avatar VARCHAR(255),
  bio VARCHAR(512),
  role VARCHAR(16) NOT NULL DEFAULT 'user',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS plant_species (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  family VARCHAR(64) DEFAULT '',
  genus VARCHAR(64) DEFAULT '',
  name VARCHAR(128) NOT NULL UNIQUE,
  alias VARCHAR(128) DEFAULT '',
  type VARCHAR(32) NOT NULL,
  origin VARCHAR(128) DEFAULT '',
  temp_min DOUBLE DEFAULT 0,
  temp_max DOUBLE DEFAULT 0,
  light_requirement VARCHAR(255) DEFAULT '',
  water_frequency VARCHAR(255) DEFAULT '',
  description TEXT,
  image_urls JSON,
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS care_articles (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(255) NOT NULL,
  content LONGTEXT,
  cover VARCHAR(255) DEFAULT '',
  topic_tag VARCHAR(32) NOT NULL,
  status VARCHAR(16) NOT NULL DEFAULT 'published',
  view_count INT NOT NULL DEFAULT 0,
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  KEY idx_articles_user (user_id),
  KEY idx_articles_topic (topic_tag)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS disease_pests (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  plant_species_id BIGINT UNSIGNED DEFAULT 0,
  name VARCHAR(128) NOT NULL,
  symptoms TEXT,
  cause TEXT,
  treatment TEXT,
  recommended_medicine VARCHAR(255) DEFAULT '',
  images JSON,
  keywords VARCHAR(255) DEFAULT '',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_pests_plant (plant_species_id),
  KEY idx_pests_keywords (keywords)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS care_reminders (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  plant_species_id BIGINT UNSIGNED DEFAULT 0,
  task_title VARCHAR(255) NOT NULL,
  remind_date DATE,
  frequency VARCHAR(32) DEFAULT '',
  status VARCHAR(16) NOT NULL DEFAULT 'pending',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_reminders_user (user_id),
  KEY idx_reminders_date (remind_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS favorites (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  target_type VARCHAR(16) NOT NULL,
  target_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_fav_user_target (user_id, target_type, target_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS user_gardens (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  plant_species_id BIGINT UNSIGNED NOT NULL,
  nickname VARCHAR(64) DEFAULT '',
  owned_since DATE,
  location VARCHAR(128) DEFAULT '',
  care_reminder_id BIGINT UNSIGNED DEFAULT 0,
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_garden_user_plant (user_id, plant_species_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS questions (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT,
  images JSON,
  status VARCHAR(16) NOT NULL DEFAULT 'open',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_questions_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS answers (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  question_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  content TEXT NOT NULL,
  is_best TINYINT(1) NOT NULL DEFAULT 0,
  like_count INT NOT NULL DEFAULT 0,
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_answers_question (question_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 种子数据
INSERT INTO users (username, email, password_hash, nickname, bio, role) VALUES
  ('admin', 'admin@gbplantwiki.local', '$2a$10$92HNAGfeO3qr7w17GkmGaOaBDxCQ7Q73gbeQ.dGGfgnIuhpPbZH4a', '园艺管理员', '平台内容维护管理员', 'admin'),
  ('gardener', 'gardener@gbplantwiki.local', '$2a$10$txSqFgLTRQZHGsde2i9vPuyh0WeH3adS0BTHSc..i8Y8FbtF4/rri', '绿手指', '热爱园艺的普通用户', 'user');

INSERT INTO plant_species (family, genus, name, alias, type, origin, temp_min, temp_max, light_requirement, water_frequency, description, image_urls) VALUES
  ('天南星科', '龟背竹属', '龟背竹', '蓬莱蕉', 'foliage', '墨西哥', 18, 30, '散射光', '每周2次', '耐阴的大型观叶植物，叶片具独特裂孔。', JSON_ARRAY('https://images.unsplash.com/photo-1524594152303-9fd13543fe6e?w=600')),
  ('景天科', '拟石莲属', '多肉吉娃娃', '吉娃娃', 'succulent', '墨西哥', 10, 28, '充足直射光', '每月2次', '叶片莲座状排列，日照充足时叶尖泛红。', JSON_ARRAY('https://images.unsplash.com/photo-1509423350716-97f9360b4e09?w=600')),
  ('睡莲科', '睡莲属', '碗莲', '微型荷花', 'aquatic', '中国', 15, 35, '全日照', '保持水位', '小型水生花卉，适合庭院水缸栽培。', JSON_ARRAY('https://images.unsplash.com/photo-1508766917616-d22f3f1eea14?w=600')),
  ('蔷薇科', '月季属', '月季', '月月红', 'flower', '中国', 5, 30, '全日照', '每周3次', '花型丰富、花期长的经典观赏花卉。', JSON_ARRAY('https://images.unsplash.com/photo-1496062031456-07b8f162a322?w=600')),
  ('百合科', '芦荟属', '库拉索芦荟', '真芦荟', 'succulent', '非洲', 10, 32, '明亮散射光', '每两周1次', '多年生常绿多肉植物，具有美容护肤价值。', JSON_ARRAY('https://images.unsplash.com/photo-1512418418704-1e26dc9c8a7a?w=600')),
  ('柏科', '圆柏属', '清香木', '细叶清香木', 'foliage', '中国西南', 8, 30, '半日照', '每周1次', '常绿灌木，叶片揉碎有清香，适合盆栽。', JSON_ARRAY('https://images.unsplash.com/photo-1463320726281-696a485928c7?w=600'));

INSERT INTO care_articles (user_id, title, content, cover, topic_tag, status, view_count) VALUES
  (1, '春季换盆全攻略：时机、方法与注意事项', '春季气温回升后是换盆的最佳时机。换盆前停止浇水3天，小心脱盆，修剪烂根并消毒，选择比原盆大1-2号的透气花盆，底部垫陶粒排水层……', 'https://images.unsplash.com/photo-1459156212016-c812468e2115?w=800', 'repotting', 'published', 128),
  (1, '多肉植物施肥要点：薄肥勤施', '多肉施肥宜稀薄，生长季每月一次稀释液肥即可，休眠期停止施肥，避免肥害烧根……', 'https://images.unsplash.com/photo-1485955900006-10f4d324d411?w=800', 'fertilizing', 'published', 96),
  (2, '月季夏季修剪与控旺', '月季夏季修剪以轻剪为主，剪除残花和细弱枝，保留健壮枝条促进复花……', 'https://images.unsplash.com/photo-1496062031456-07b8f162a322?w=800', 'pruning', 'published', 210),
  (1, '常见介壳虫的识别与防治', '介壳虫常附着在叶背和枝干，可用酒精棉擦拭，严重时喷洒矿物油乳剂……', 'https://images.unsplash.com/photo-1530836369250-ef72a3f5cda8?w=800', 'pest_control', 'published', 154),
  (2, '龟背竹扦插繁殖实操', '选取带气生根的健壮枝条，切口晾干后插入湿润的蛭石中，保持湿度约三周生根……', 'https://images.unsplash.com/photo-1524594152303-9fd13543fe6e?w=800', 'propagation', 'published', 67);

INSERT INTO disease_pests (plant_species_id, name, symptoms, cause, treatment, recommended_medicine, keywords, images) VALUES
  (4, '月季黑斑病', '叶片出现黑色圆形斑点，边缘呈放射状，严重时叶片脱落。', '高温高湿、通风不良，病原为蔷薇黑斑菌。', '及时摘除病叶，喷施代森锰锌或苯醚甲环唑，每周一次连续2-3次。', '代森锰锌、苯醚甲环唑', '黑斑,黄叶,月季', JSON_ARRAY()),
  (2, '多肉介壳虫', '叶腋处出现白色棉絮状物，叶片发黏发黄。', '通风差、湿度大，虫源为蚧壳虫若虫。', '人工刮除后用酒精擦拭，严重时喷施噻嗪酮。', '噻嗪酮、矿物油乳剂', '介壳虫,白色,黏', JSON_ARRAY()),
  (1, '龟背竹叶斑病', '叶片出现褐色水渍状病斑，逐渐扩大干枯。', '浇水过多、长期积水，病原真菌感染。', '控水通风，剪除病叶，喷施多菌灵。', '多菌灵', '叶斑,烂叶', JSON_ARRAY()),
  (0, '红蜘蛛', '叶面出现细密黄白色斑点，叶背有蛛网。', '空气干燥、高温，螨虫滋生。', '增加湿度，喷施阿维菌素或哒螨灵。', '阿维菌素、哒螨灵', '红蜘蛛,螨,黄点', JSON_ARRAY());

INSERT INTO care_reminders (user_id, plant_species_id, task_title, remind_date, frequency, status) VALUES
  (2, 4, '给月季补充缓释肥', DATE_ADD(CURDATE(), INTERVAL 3 DAY), 'monthly', 'pending'),
  (2, 1, '龟背竹叶片擦拭除尘', DATE_ADD(CURDATE(), INTERVAL 1 DAY), 'weekly', 'pending');

INSERT INTO questions (user_id, title, content, images, status) VALUES
  (2, '新买的月季叶子发黄怎么办？', '刚上盆一周，叶片边缘发黄，是不是浇水太多？', JSON_ARRAY(), 'open'),
  (2, '多肉徒长了如何补救？', '冬季光照不足，多肉长高了，可以砍头吗？', JSON_ARRAY(), 'open');

INSERT INTO answers (question_id, user_id, content, is_best, like_count) VALUES
  (1, 1, '新上盆植物根系未恢复，建议先放在散射光处缓苗，见干见湿浇水，避免积水。', 0, 5),
  (2, 1, '可以砍头繁殖，砍下的头部晾干后重新扦插，母株会萌发侧芽。', 0, 8);

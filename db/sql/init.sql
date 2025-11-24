GRANT ALL PRIVILEGES ON test_develop_ecs-express-mode.* TO 'user'@'%';
GRANT ALL PRIVILEGES ON develop_develop_ecs-express-mode.* TO 'user'@'%';
FLUSH PRIVILEGES;

CREATE DATABASE IF NOT EXISTS test_develop_ecs-express-mode;
CREATE DATABASE IF NOT EXISTS develop_develop_ecs-express-mode;

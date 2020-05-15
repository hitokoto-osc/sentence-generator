# 一言句子集生成工具
本工具用于定时同步句子，并推送给指定 Git 仓库（需要预先配置仓库验证权限）
## 使用
```
yarn && yarn start
```

## 涉及的 数据库结构（SQL）
`hitokoto_sentence`：
```SQL
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for hitokoto_sentence
-- ----------------------------
DROP TABLE IF EXISTS `hitokoto_sentence`;
CREATE TABLE `hitokoto_sentence`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `hitokoto` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `type` char(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `from` char(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `from_who` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `creator` char(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'hitokoto',
  `creator_uid` int(11) NULL DEFAULT NULL,
  `reviewer` int(32) NOT NULL DEFAULT 1,
  `commit_from` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'web',
  `assessor` char(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `owner` char(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0',
  `created_at` char(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uuid`(`uuid`) USING BTREE,
  INDEX `reviewer`(`reviewer`) USING BTREE,
  INDEX `creator_uid`(`creator_uid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6017 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

SET FOREIGN_KEY_CHECKS = 1;
```

`hitokoto_categories`:
```SQL
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for hitokoto_categories
-- ----------------------------
DROP TABLE IF EXISTS `hitokoto_categories`;
CREATE TABLE `hitokoto_categories`  (
  `id` int(255) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `desc` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '描述分类',
  `key` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '分类键名',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`, `key`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

SET FOREIGN_KEY_CHECKS = 1;

```

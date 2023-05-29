
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL';

-- -----------------------------------------------------
-- Table `users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `users` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`)
);

-- -----------------------------------------------------
-- Table `assets`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `assets` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `user_id` INT NOT NULL,
  `type` ENUM('charts', 'insights', 'audiences') NOT NULL,
  `description` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_user_id` (`user_id`), 
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);

-- -----------------------------------------------------
-- Table `charts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `charts` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `assets_id` INT NOT NULL,
  `title` VARCHAR(255),
  `axes_titles` VARCHAR(255),
  `data` VARCHAR(255),
  PRIMARY KEY (`id`),
  INDEX `idx_assets_id` (`assets_id`),
  FOREIGN KEY (`assets_id`) REFERENCES `assets` (`id`) ON DELETE CASCADE
);

-- -----------------------------------------------------
-- Table `insights`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `insights` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `assets_id` INT NOT NULL,
  `text` VARCHAR(255),
  PRIMARY KEY (`id`),
  INDEX `idx_assets_id` (`assets_id`),
  FOREIGN KEY (`assets_id`) REFERENCES `assets` (`id`) ON DELETE CASCADE
);

-- -----------------------------------------------------
-- Table `audiences`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `audiences` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `assets_id` INT NOT NULL,
  `characteristics` VARCHAR(255),
  PRIMARY KEY (`id`),
  INDEX `idx_assets_id` (`assets_id`),
  FOREIGN KEY (`assets_id`) REFERENCES `assets` (`id`) ON DELETE CASCADE
);

-- -----------------------------------------------------
-- Table `permissions`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `permissions` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `token` VARCHAR(255),
  `create` BOOLEAN,
  `read` BOOLEAN,
  `update` BOOLEAN,
  `delete` BOOLEAN,
  PRIMARY KEY (`id`)
);

-- -----------------------------------------------------
-- Table `users_permissions`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `users_permissions` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `users_id` INT NOT NULL,
  `permissions_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`users_id`) REFERENCES `users` (`id`),
  FOREIGN KEY (`permissions_id`) REFERENCES `permissions` (`id`)
);

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
-- -----------------------------------------------------
-- Drop Tables
-- -----------------------------------------------------

-- Drop table `audiences`
DROP TABLE IF EXISTS `audiences`;

-- Drop table `insights`
DROP TABLE IF EXISTS `insights`;

-- Drop table `charts`
DROP TABLE IF EXISTS `charts`;

-- Drop table `assets`
DROP TABLE IF EXISTS `assets`;

-- Drop table `users`
DROP TABLE IF EXISTS `users`;

-- Drop table `users`
DROP TABLE IF EXISTS `users`;

-- Drop table `permissions`
DROP TABLE IF EXISTS `permissions`;

-- Drop table `users_permissions`
DROP TABLE IF EXISTS `users_permissions`;


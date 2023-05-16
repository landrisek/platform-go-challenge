
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
  `type` ENUM('Chart', 'Insight', 'Audience') NOT NULL,
  `title` VARCHAR(100) NOT NULL,
  `user_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);

-- -----------------------------------------------------
-- Table `charts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `charts` (
  `id` INT NOT NULL,
  `axes_titles` VARCHAR(255),
  `data` VARCHAR(255),
  `description` VARCHAR(255),
  PRIMARY KEY (`id`),
  FOREIGN KEY (`id`) REFERENCES `assets` (`id`)
);

-- -----------------------------------------------------
-- Table `insights`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `insights` (
  `id` INT NOT NULL,
  `text` VARCHAR(255),
  `description` VARCHAR(255),
  PRIMARY KEY (`id`),
  FOREIGN KEY (`id`) REFERENCES `assets` (`id`)
);

-- -----------------------------------------------------
-- Table `audiences`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `audiences` (
  `id` INT NOT NULL,
  `characteristics` VARCHAR(255),
  `description` VARCHAR(255),
  PRIMARY KEY (`id`),
  FOREIGN KEY (`id`) REFERENCES `assets` (`id`)
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


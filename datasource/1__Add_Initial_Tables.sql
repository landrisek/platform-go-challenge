
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL';

-- -----------------------------------------------------
-- Table `user`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `user` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`)
);

-- -----------------------------------------------------
-- Table `asset`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `asset` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `type_id` INT NOT NULL,
  `description` TEXT,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`type_id`) REFERENCES `asset_type` (`id`)
);

-- -----------------------------------------------------
-- Table `asset_type`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `asset_type` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `type` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`)
);

-- -----------------------------------------------------
-- Table `user_favorite`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_favorite` (
  `user_id` INT NOT NULL,
  `asset_id` INT NOT NULL,
  PRIMARY KEY (`user_id`, `asset_id`),
  FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  FOREIGN KEY (`asset_id`) REFERENCES `asset` (`id`)
);

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `user`;
DROP TABLE `asset`;
DROP TABLE `asset_type`;
DROP TABLE `user_favorite`;
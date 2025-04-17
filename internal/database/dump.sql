-- MySQL dump 10.13  Distrib 8.0.42, for Linux (x86_64)
--
-- Host: localhost    Database: forumdb
-- ------------------------------------------------------
-- Server version	8.0.42

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `category`
--

DROP TABLE IF EXISTS `category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `category` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `category`
--

LOCK TABLES `category` WRITE;
/*!40000 ALTER TABLE `category` DISABLE KEYS */;
INSERT INTO `category` VALUES (4,'Adaptations de Livres en Films'),(2,'Classiques du Cinéma'),(3,'Critiques et Recommandations'),(5,'Films Cultes'),(1,'Genres de Films');
/*!40000 ALTER TABLE `category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comment_likes`
--

DROP TABLE IF EXISTS `comment_likes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comment_likes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `comment_id` int NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `value` int DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `comment_id` (`comment_id`,`user_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `comment_likes_ibfk_1` FOREIGN KEY (`comment_id`) REFERENCES `comments` (`id`),
  CONSTRAINT `comment_likes_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=114 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comment_likes`
--

LOCK TABLES `comment_likes` WRITE;
/*!40000 ALTER TABLE `comment_likes` DISABLE KEYS */;
INSERT INTO `comment_likes` VALUES (1,1,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','2025-03-25 10:05:23',-1),(2,4,'308e9405-784d-48bc-ab7f-393c25778f43','2025-04-08 10:05:23',1),(3,5,'359b955d-f28f-4a44-b759-60e3509010ff','2025-02-18 10:05:23',-1),(4,7,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','2025-03-31 10:05:23',-1),(5,8,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','2025-03-31 10:05:23',-1),(6,10,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','2025-03-24 10:05:23',-1),(7,11,'df2c278e-f920-473e-b93d-ca31c9641e07','2025-03-24 10:05:23',-1),(8,12,'308e9405-784d-48bc-ab7f-393c25778f43','2025-03-24 10:05:23',1),(9,13,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','2025-04-04 10:05:23',1),(10,14,'440a6b09-b04b-4959-b851-bb129ae70568','2025-04-04 10:05:23',1),(11,16,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','2025-04-04 10:05:23',1),(12,17,'440a6b09-b04b-4959-b851-bb129ae70568','2025-03-14 10:05:23',1),(13,18,'308e9405-784d-48bc-ab7f-393c25778f43','2025-02-18 10:05:23',1),(14,21,'308e9405-784d-48bc-ab7f-393c25778f43','2025-04-11 10:05:23',-1),(15,22,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','2025-04-11 10:05:23',-1),(16,23,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','2025-04-17 10:05:23',1),(17,27,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','2025-02-20 10:05:23',1),(18,28,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','2025-02-20 10:05:23',-1),(19,29,'359b955d-f28f-4a44-b759-60e3509010ff','2025-02-20 10:05:23',-1),(20,32,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','2025-03-14 10:05:23',1),(21,36,'df2c278e-f920-473e-b93d-ca31c9641e07','2025-04-12 10:05:23',-1),(22,41,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','2025-03-20 10:05:23',-1),(23,43,'359b955d-f28f-4a44-b759-60e3509010ff','2025-04-03 10:05:23',-1),(24,45,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','2025-03-12 10:05:23',1),(25,46,'308e9405-784d-48bc-ab7f-393c25778f43','2025-03-12 10:05:23',-1),(26,47,'308e9405-784d-48bc-ab7f-393c25778f43','2025-03-12 10:05:23',-1),(27,49,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','2025-02-27 10:05:23',-1),(28,53,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','2025-03-24 10:05:23',-1),(29,55,'df2c278e-f920-473e-b93d-ca31c9641e07','2025-02-22 10:05:23',1),(30,57,'df2c278e-f920-473e-b93d-ca31c9641e07','2025-02-19 10:05:23',1),(31,62,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','2025-04-05 10:05:23',-1),(32,65,'23467e97-565c-484c-a34b-bbd2849f2f21','2025-03-30 10:05:23',-1),(33,67,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','2025-03-02 10:05:23',-1),(34,73,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','2025-03-28 10:05:23',1),(35,75,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','2025-02-28 10:05:23',-1),(36,76,'df2c278e-f920-473e-b93d-ca31c9641e07','2025-04-01 10:05:23',1),(37,77,'308e9405-784d-48bc-ab7f-393c25778f43','2025-03-28 10:05:23',1),(38,79,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','2025-03-13 10:05:23',1),(39,80,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','2025-04-05 10:05:23',-1),(40,82,'440a6b09-b04b-4959-b851-bb129ae70568','2025-02-23 10:05:23',1),(41,83,'df2c278e-f920-473e-b93d-ca31c9641e07','2025-03-28 10:05:23',1),(42,84,'23467e97-565c-484c-a34b-bbd2849f2f21','2025-03-28 10:05:23',-1),(43,87,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','2025-03-15 10:05:23',-1),(44,90,'df2c278e-f920-473e-b93d-ca31c9641e07','2025-03-02 10:05:23',-1),(45,92,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','2025-04-07 10:05:23',1);
/*!40000 ALTER TABLE `comment_likes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comment_reports`
--

DROP TABLE IF EXISTS `comment_reports`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comment_reports` (
  `id` int NOT NULL AUTO_INCREMENT,
  `comment_id` int NOT NULL,
  `reporter_id` varchar(255) NOT NULL,
  `reason` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `comment_id` (`comment_id`),
  KEY `reporter_id` (`reporter_id`),
  CONSTRAINT `comment_reports_ibfk_1` FOREIGN KEY (`comment_id`) REFERENCES `comments` (`id`),
  CONSTRAINT `comment_reports_ibfk_2` FOREIGN KEY (`reporter_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comment_reports`
--

LOCK TABLES `comment_reports` WRITE;
/*!40000 ALTER TABLE `comment_reports` DISABLE KEYS */;
INSERT INTO `comment_reports` VALUES (1,1,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Insulte','2025-03-25 10:05:23'),(2,2,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Insulte','2025-04-08 10:05:23'),(3,3,'308e9405-784d-48bc-ab7f-393c25778f43','Spam','2025-04-08 10:05:23'),(4,4,'df2c278e-f920-473e-b93d-ca31c9641e07','Insulte','2025-04-08 10:05:23'),(5,8,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Troll','2025-03-31 10:05:23'),(6,9,'440a6b09-b04b-4959-b851-bb129ae70568','Troll','2025-02-18 10:05:23'),(7,15,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Spam','2025-04-04 10:05:23'),(8,16,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Insulte','2025-04-04 10:05:23'),(9,17,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Spam','2025-03-14 10:05:23'),(10,18,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Troll','2025-02-18 10:05:23'),(11,19,'df2c278e-f920-473e-b93d-ca31c9641e07','Insulte','2025-02-18 10:05:23'),(12,20,'359b955d-f28f-4a44-b759-60e3509010ff','Spam','2025-02-18 10:05:23'),(13,22,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Insulte','2025-04-11 10:05:23'),(14,29,'23467e97-565c-484c-a34b-bbd2849f2f21','Troll','2025-02-20 10:05:23'),(15,31,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Troll','2025-03-14 10:05:23'),(16,32,'df2c278e-f920-473e-b93d-ca31c9641e07','Spam','2025-03-14 10:05:23'),(17,34,'308e9405-784d-48bc-ab7f-393c25778f43','Insulte','2025-04-12 10:05:23'),(18,35,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Spam','2025-04-12 10:05:23'),(19,36,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Insulte','2025-04-12 10:05:23'),(20,37,'440a6b09-b04b-4959-b851-bb129ae70568','Troll','2025-03-25 10:05:23'),(21,42,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Spam','2025-04-03 10:05:23'),(22,45,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Insulte','2025-03-12 10:05:23'),(23,47,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Spam','2025-03-12 10:05:23'),(24,48,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Troll','2025-02-27 10:05:23'),(25,49,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Insulte','2025-02-27 10:05:23'),(26,50,'440a6b09-b04b-4959-b851-bb129ae70568','Spam','2025-03-04 10:05:23'),(27,52,'308e9405-784d-48bc-ab7f-393c25778f43','Insulte','2025-03-24 10:05:23'),(28,53,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Insulte','2025-03-24 10:05:23'),(29,54,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Spam','2025-03-24 10:05:23'),(30,55,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Insulte','2025-02-22 10:05:23'),(31,57,'440a6b09-b04b-4959-b851-bb129ae70568','Insulte','2025-02-19 10:05:23'),(32,64,'308e9405-784d-48bc-ab7f-393c25778f43','Spam','2025-03-09 10:05:23'),(33,66,'359b955d-f28f-4a44-b759-60e3509010ff','Insulte','2025-03-02 10:05:23'),(34,67,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Spam','2025-03-02 10:05:23'),(35,68,'440a6b09-b04b-4959-b851-bb129ae70568','Troll','2025-03-18 10:05:23'),(36,72,'308e9405-784d-48bc-ab7f-393c25778f43','Spam','2025-03-28 10:05:23'),(37,75,'359b955d-f28f-4a44-b759-60e3509010ff','Insulte','2025-02-28 10:05:23'),(38,76,'440a6b09-b04b-4959-b851-bb129ae70568','Spam','2025-04-01 10:05:23'),(39,80,'23467e97-565c-484c-a34b-bbd2849f2f21','Insulte','2025-04-05 10:05:23'),(40,81,'23467e97-565c-484c-a34b-bbd2849f2f21','Spam','2025-02-23 10:05:23'),(41,82,'df2c278e-f920-473e-b93d-ca31c9641e07','Insulte','2025-02-23 10:05:23'),(42,85,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Troll','2025-03-28 10:05:23'),(43,86,'359b955d-f28f-4a44-b759-60e3509010ff','Spam','2025-03-19 10:05:23'),(44,89,'308e9405-784d-48bc-ab7f-393c25778f43','Troll','2025-03-02 10:05:23'),(45,92,'df2c278e-f920-473e-b93d-ca31c9641e07','Insulte','2025-04-07 10:05:23'),(46,93,'359b955d-f28f-4a44-b759-60e3509010ff','Spam','2025-04-16 10:05:23');
/*!40000 ALTER TABLE `comment_reports` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comments`
--

DROP TABLE IF EXISTS `comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `post_id` int NOT NULL,
  `author_id` varchar(36) NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `post_id` (`post_id`),
  KEY `author_id` (`author_id`),
  CONSTRAINT `comments_ibfk_1` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`),
  CONSTRAINT `comments_ibfk_2` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=101 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comments`
--

LOCK TABLES `comments` WRITE;
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
INSERT INTO `comments` VALUES (1,1,'440a6b09-b04b-4959-b851-bb129ae70568','Un film que je recommande chaudement.','2025-03-25 10:05:23'),(2,1,'23467e97-565c-484c-a34b-bbd2849f2f21','J’ai quelques réserves, mais dans l’ensemble c’est très réussi.','2025-03-25 10:05:23'),(3,2,'308e9405-784d-48bc-ab7f-393c25778f43','Un scénario solide, des personnages profonds.','2025-04-08 10:05:23'),(4,2,'23467e97-565c-484c-a34b-bbd2849f2f21','L’histoire est bien transposée, les acteurs sont incroyables.','2025-04-08 10:05:23'),(5,3,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Une adaptation fidèle et ambitieuse.','2025-02-18 10:05:23'),(6,4,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','C’est le genre de film qu’on n’oublie pas facilement.','2025-03-31 10:05:23'),(7,4,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Un film que je recommande chaudement.','2025-03-31 10:05:23'),(8,4,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','J’ai quelques réserves, mais dans l’ensemble c’est très réussi.','2025-03-31 10:05:23'),(9,5,'440a6b09-b04b-4959-b851-bb129ae70568','J’ai été ému par la simplicité et la puissance du message.','2025-02-18 10:05:23'),(10,5,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','J’ai quelques réserves, mais dans l’ensemble c’est très réussi.','2025-02-18 10:05:23'),(11,6,'308e9405-784d-48bc-ab7f-393c25778f43','Une adaptation fidèle et ambitieuse.','2025-03-24 10:05:23'),(12,7,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','J’ai été ému par la simplicité et la puissance du message.','2025-04-04 10:05:23'),(13,8,'440a6b09-b04b-4959-b851-bb129ae70568','Une adaptation fidèle et ambitieuse.','2025-04-04 10:05:23'),(14,8,'23467e97-565c-484c-a34b-bbd2849f2f21','Un film à voir au moins une fois dans sa vie.','2025-04-04 10:05:23'),(15,8,'df2c278e-f920-473e-b93d-ca31c9641e07','Un film à voir au moins une fois dans sa vie.','2025-04-04 10:05:23'),(16,9,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','L’histoire est bien transposée, les acteurs sont incroyables.','2025-03-14 10:05:23'),(17,10,'359b955d-f28f-4a44-b759-60e3509010ff','C’est le genre de film qu’on n’oublie pas facilement.','2025-02-18 10:05:23'),(18,11,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Un scénario solide, des personnages profonds.','2025-04-11 10:05:23'),(19,11,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','C’est le genre de film qu’on n’oublie pas facilement.','2025-04-11 10:05:23'),(20,12,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Un film que je recommande chaudement.','2025-04-17 10:05:23'),(21,12,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Un scénario solide, des personnages profonds.','2025-04-17 10:05:23'),(22,12,'23467e97-565c-484c-a34b-bbd2849f2f21','Un film que je recommande chaudement.','2025-04-17 10:05:23'),(23,13,'308e9405-784d-48bc-ab7f-393c25778f43','L’histoire est bien transposée, les acteurs sont incroyables.','2025-04-17 10:05:23'),(24,14,'440a6b09-b04b-4959-b851-bb129ae70568','Certains passages traînent en longueur, mais le fond est fort.','2025-03-04 10:05:23'),(25,15,'359b955d-f28f-4a44-b759-60e3509010ff','L’histoire est bien transposée, les acteurs sont incroyables.','2025-02-20 10:05:23'),(26,15,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Certains passages traînent en longueur, mais le fond est fort.','2025-02-20 10:05:23'),(27,15,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','J’ai été ému par la simplicité et la puissance du message.','2025-02-20 10:05:23'),(28,16,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Certains passages traînent en longueur, mais le fond est fort.','2025-03-14 10:05:23'),(29,17,'df2c278e-f920-473e-b93d-ca31c9641e07','Certains passages traînent en longueur, mais le fond est fort.','2025-03-01 10:05:23'),(30,17,'308e9405-784d-48bc-ab7f-393c25778f43','J’ai quelques réserves, mais dans l’ensemble c’est très réussi.','2025-03-01 10:05:23'),(31,18,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','J’ai quelques réserves, mais dans l’ensemble c’est très réussi.','2025-04-12 10:05:23'),(32,18,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Un film que je recommande chaudement.','2025-04-12 10:05:23'),(33,19,'df2c278e-f920-473e-b93d-ca31c9641e07','Certains passages traînent en longueur, mais le fond est fort.','2025-03-25 10:05:23'),(34,20,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Une adaptation fidèle et ambitieuse.','2025-03-06 10:05:23'),(35,21,'23467e97-565c-484c-a34b-bbd2849f2f21','Un scénario solide, des personnages profonds.','2025-03-08 10:05:23'),(36,22,'23467e97-565c-484c-a34b-bbd2849f2f21','C’est le genre de film qu’on n’oublie pas facilement.','2025-03-20 10:05:23'),(37,22,'308e9405-784d-48bc-ab7f-393c25778f43','L’histoire est bien transposée, les acteurs sont incroyables.','2025-03-20 10:05:23'),(38,23,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Un scénario solide, des personnages profonds.','2025-04-03 10:05:23'),(39,23,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','C’est le genre de film qu’on n’oublie pas facilement.','2025-04-03 10:05:23'),(40,24,'308e9405-784d-48bc-ab7f-393c25778f43','C’est le genre de film qu’on n’oublie pas facilement.','2025-03-12 10:05:23'),(41,24,'23467e97-565c-484c-a34b-bbd2849f2f21','Un film à voir au moins une fois dans sa vie.','2025-03-12 10:05:23'),(42,25,'23467e97-565c-484c-a34b-bbd2849f2f21','Certains passages traînent en longueur, mais le fond est fort.','2025-02-27 10:05:23'),(43,26,'440a6b09-b04b-4959-b851-bb129ae70568','Une adaptation fidèle et ambitieuse.','2025-03-04 10:05:23'),(44,26,'359b955d-f28f-4a44-b759-60e3509010ff','C’est le genre de film qu’on n’oublie pas facilement.','2025-03-04 10:05:23'),(45,27,'df2c278e-f920-473e-b93d-ca31c9641e07','J’ai quelques réserves, mais dans l’ensemble c’est très réussi.','2025-03-24 10:05:23'),(46,28,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Je viens de revoir ce film, et je suis encore bluffé par la réalisation.','2025-02-22 10:05:23'),(47,29,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Un film à voir au moins une fois dans sa vie.','2025-02-19 10:05:23'),(48,30,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Je viens de revoir ce film, et je suis encore bluffé par la réalisation.','2025-03-01 10:05:23'),(49,30,'359b955d-f28f-4a44-b759-60e3509010ff','Une adaptation fidèle et ambitieuse.','2025-03-01 10:05:23'),(50,30,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','J’ai été ému par la simplicité et la puissance du message.','2025-03-01 10:05:23'),(51,31,'440a6b09-b04b-4959-b851-bb129ae70568','Je viens de revoir ce film, et je suis encore bluffé par la réalisation.','2025-03-02 10:05:23'),(52,31,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Je viens de revoir ce film, et je suis encore bluffé par la réalisation.','2025-03-02 10:05:23'),(53,31,'308e9405-784d-48bc-ab7f-393c25778f43','J’ai été ému par la simplicité et la puissance du message.','2025-03-02 10:05:23'),(54,32,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Un film à voir au moins une fois dans sa vie.','2025-04-05 10:05:23'),(55,32,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','L’histoire est bien transposée, les acteurs sont incroyables.','2025-04-05 10:05:23'),(56,32,'308e9405-784d-48bc-ab7f-393c25778f43','Un film que je recommande chaudement.','2025-04-05 10:05:23'),(57,33,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','J’ai été ému par la simplicité et la puissance du message.','2025-03-09 10:05:23'),(58,34,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Un film que je recommande chaudement.','2025-03-30 10:05:23'),(59,34,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Un film à voir au moins une fois dans sa vie.','2025-03-30 10:05:23'),(60,35,'359b955d-f28f-4a44-b759-60e3509010ff','C’est le genre de film qu’on n’oublie pas facilement.','2025-03-02 10:05:23'),(61,36,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Certains passages traînent en longueur, mais le fond est fort.','2025-03-18 10:05:23'),(62,36,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Un scénario solide, des personnages profonds.','2025-03-18 10:05:23'),(63,36,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','J’ai été ému par la simplicité et la puissance du message.','2025-03-18 10:05:23'),(64,37,'df2c278e-f920-473e-b93d-ca31c9641e07','Certains passages traînent en longueur, mais le fond est fort.','2025-03-04 10:05:23'),(65,37,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','J’ai été ému par la simplicité et la puissance du message.','2025-03-04 10:05:23'),(66,37,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','L’histoire est bien transposée, les acteurs sont incroyables.','2025-03-04 10:05:23'),(67,38,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','L’histoire est bien transposée, les acteurs sont incroyables.','2025-03-28 10:05:23'),(68,38,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Un film que je recommande chaudement.','2025-03-28 10:05:23'),(69,38,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Je viens de revoir ce film, et je suis encore bluffé par la réalisation.','2025-03-28 10:05:23'),(70,39,'359b955d-f28f-4a44-b759-60e3509010ff','L’histoire est bien transposée, les acteurs sont incroyables.','2025-02-28 10:05:23'),(71,39,'359b955d-f28f-4a44-b759-60e3509010ff','Je viens de revoir ce film, et je suis encore bluffé par la réalisation.','2025-02-28 10:05:23'),(72,40,'23467e97-565c-484c-a34b-bbd2849f2f21','C’est le genre de film qu’on n’oublie pas facilement.','2025-04-01 10:05:23'),(73,40,'df2c278e-f920-473e-b93d-ca31c9641e07','Un film que je recommande chaudement.','2025-04-01 10:05:23'),(74,40,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Une adaptation fidèle et ambitieuse.','2025-04-01 10:05:23'),(75,41,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Un film que je recommande chaudement.','2025-03-28 10:05:23'),(76,41,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Un film que je recommande chaudement.','2025-03-28 10:05:23'),(77,41,'308e9405-784d-48bc-ab7f-393c25778f43','Un film à voir au moins une fois dans sa vie.','2025-03-28 10:05:23'),(78,42,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Une adaptation fidèle et ambitieuse.','2025-03-13 10:05:23'),(79,42,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','L’histoire est bien transposée, les acteurs sont incroyables.','2025-03-13 10:05:23'),(80,42,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Je viens de revoir ce film, et je suis encore bluffé par la réalisation.','2025-03-13 10:05:23'),(81,43,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Certains passages traînent en longueur, mais le fond est fort.','2025-04-05 10:05:23'),(82,43,'359b955d-f28f-4a44-b759-60e3509010ff','L’histoire est bien transposée, les acteurs sont incroyables.','2025-04-05 10:05:23'),(83,43,'308e9405-784d-48bc-ab7f-393c25778f43','Un scénario solide, des personnages profonds.','2025-04-05 10:05:23'),(84,44,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','L’histoire est bien transposée, les acteurs sont incroyables.','2025-02-23 10:05:23'),(85,44,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Certains passages traînent en longueur, mais le fond est fort.','2025-02-23 10:05:23'),(86,44,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','L’histoire est bien transposée, les acteurs sont incroyables.','2025-02-23 10:05:23'),(87,45,'df2c278e-f920-473e-b93d-ca31c9641e07','J’ai été ému par la simplicité et la puissance du message.','2025-03-28 10:05:23'),(88,46,'308e9405-784d-48bc-ab7f-393c25778f43','Un film à voir au moins une fois dans sa vie.','2025-03-19 10:05:23'),(89,46,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Un film à voir au moins une fois dans sa vie.','2025-03-19 10:05:23'),(90,46,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','J’ai quelques réserves, mais dans l’ensemble c’est très réussi.','2025-03-19 10:05:23'),(91,47,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Certains passages traînent en longueur, mais le fond est fort.','2025-03-15 10:05:23'),(92,47,'440a6b09-b04b-4959-b851-bb129ae70568','Un film que je recommande chaudement.','2025-03-15 10:05:23'),(93,47,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Une adaptation fidèle et ambitieuse.','2025-03-15 10:05:23'),(94,48,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','J’ai quelques réserves, mais dans l’ensemble c’est très réussi.','2025-03-02 10:05:23'),(95,48,'23467e97-565c-484c-a34b-bbd2849f2f21','Un scénario solide, des personnages profonds.','2025-03-02 10:05:23'),(96,48,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','J’ai quelques réserves, mais dans l’ensemble c’est très réussi.','2025-03-02 10:05:23'),(97,49,'440a6b09-b04b-4959-b851-bb129ae70568','Un film à voir au moins une fois dans sa vie.','2025-04-07 10:05:23'),(98,49,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Certains passages traînent en longueur, mais le fond est fort.','2025-04-07 10:05:23'),(99,50,'df2c278e-f920-473e-b93d-ca31c9641e07','Une adaptation fidèle et ambitieuse.','2025-04-16 10:05:23'),(100,50,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Je viens de revoir ce film, et je suis encore bluffé par la réalisation.','2025-04-16 10:05:23');
/*!40000 ALTER TABLE `comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `likes`
--

DROP TABLE IF EXISTS `likes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `likes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(36) NOT NULL,
  `post_id` int NOT NULL,
  `value` int NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `post_id` (`post_id`),
  CONSTRAINT `likes_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `likes_ibfk_2` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=115 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `likes`
--

LOCK TABLES `likes` WRITE;
/*!40000 ALTER TABLE `likes` DISABLE KEYS */;
INSERT INTO `likes` VALUES (1,'359b955d-f28f-4a44-b759-60e3509010ff',1,1,'2025-03-25 10:05:23'),(2,'ce973ab7-dfb0-41d0-9813-a16d6f30f925',1,1,'2025-03-25 10:05:23'),(3,'ce973ab7-dfb0-41d0-9813-a16d6f30f925',4,-1,'2025-03-31 10:05:23'),(4,'440a6b09-b04b-4959-b851-bb129ae70568',4,-1,'2025-03-31 10:05:23'),(5,'df2c278e-f920-473e-b93d-ca31c9641e07',6,1,'2025-03-24 10:05:23'),(6,'b7ef6093-5bd1-4cca-b1cf-b578cd569221',7,1,'2025-04-04 10:05:23'),(7,'ce973ab7-dfb0-41d0-9813-a16d6f30f925',9,-1,'2025-03-14 10:05:23'),(8,'359b955d-f28f-4a44-b759-60e3509010ff',9,-1,'2025-03-14 10:05:23'),(9,'440a6b09-b04b-4959-b851-bb129ae70568',11,1,'2025-04-11 10:05:23'),(10,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927',11,1,'2025-04-11 10:05:23'),(11,'b7ef6093-5bd1-4cca-b1cf-b578cd569221',13,1,'2025-04-17 10:05:23'),(12,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9',13,-1,'2025-04-17 10:05:23'),(13,'df2c278e-f920-473e-b93d-ca31c9641e07',14,-1,'2025-03-04 10:05:23'),(14,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1',15,-1,'2025-02-20 10:05:23'),(15,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9',15,-1,'2025-02-20 10:05:23'),(16,'359b955d-f28f-4a44-b759-60e3509010ff',17,-1,'2025-03-01 10:05:23'),(17,'23467e97-565c-484c-a34b-bbd2849f2f21',18,1,'2025-04-12 10:05:23'),(18,'440a6b09-b04b-4959-b851-bb129ae70568',19,-1,'2025-03-25 10:05:23'),(19,'308e9405-784d-48bc-ab7f-393c25778f43',20,-1,'2025-03-06 10:05:23'),(20,'308e9405-784d-48bc-ab7f-393c25778f43',20,1,'2025-03-06 10:05:23'),(21,'308e9405-784d-48bc-ab7f-393c25778f43',22,1,'2025-03-20 10:05:23'),(22,'23467e97-565c-484c-a34b-bbd2849f2f21',22,-1,'2025-03-20 10:05:23'),(23,'359b955d-f28f-4a44-b759-60e3509010ff',23,1,'2025-04-03 10:05:23'),(24,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927',23,-1,'2025-04-03 10:05:23'),(25,'df2c278e-f920-473e-b93d-ca31c9641e07',24,-1,'2025-03-12 10:05:23'),(26,'359b955d-f28f-4a44-b759-60e3509010ff',24,-1,'2025-03-12 10:05:23'),(27,'440a6b09-b04b-4959-b851-bb129ae70568',26,-1,'2025-03-04 10:05:23'),(28,'b7ef6093-5bd1-4cca-b1cf-b578cd569221',26,-1,'2025-03-04 10:05:23'),(29,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927',27,-1,'2025-03-24 10:05:23'),(30,'df2c278e-f920-473e-b93d-ca31c9641e07',27,-1,'2025-03-24 10:05:23'),(31,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1',28,-1,'2025-02-22 10:05:23'),(32,'440a6b09-b04b-4959-b851-bb129ae70568',29,-1,'2025-02-19 10:05:23'),(33,'359b955d-f28f-4a44-b759-60e3509010ff',30,1,'2025-03-01 10:05:23'),(34,'b7ef6093-5bd1-4cca-b1cf-b578cd569221',31,-1,'2025-03-02 10:05:23'),(35,'308e9405-784d-48bc-ab7f-393c25778f43',32,1,'2025-04-05 10:05:23'),(36,'440a6b09-b04b-4959-b851-bb129ae70568',32,-1,'2025-04-05 10:05:23'),(37,'359b955d-f28f-4a44-b759-60e3509010ff',33,-1,'2025-03-09 10:05:23'),(38,'308e9405-784d-48bc-ab7f-393c25778f43',33,-1,'2025-03-09 10:05:23'),(39,'440a6b09-b04b-4959-b851-bb129ae70568',34,-1,'2025-03-30 10:05:23'),(40,'440a6b09-b04b-4959-b851-bb129ae70568',34,1,'2025-03-30 10:05:23'),(41,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927',35,1,'2025-03-02 10:05:23'),(42,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9',35,1,'2025-03-02 10:05:23'),(43,'df2c278e-f920-473e-b93d-ca31c9641e07',36,-1,'2025-03-18 10:05:23'),(44,'359b955d-f28f-4a44-b759-60e3509010ff',37,-1,'2025-03-04 10:05:23'),(45,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927',37,-1,'2025-03-04 10:05:23'),(46,'23467e97-565c-484c-a34b-bbd2849f2f21',38,-1,'2025-03-28 10:05:23'),(47,'df2c278e-f920-473e-b93d-ca31c9641e07',38,1,'2025-03-28 10:05:23'),(48,'308e9405-784d-48bc-ab7f-393c25778f43',39,-1,'2025-02-28 10:05:23'),(49,'440a6b09-b04b-4959-b851-bb129ae70568',39,1,'2025-02-28 10:05:23'),(50,'ce973ab7-dfb0-41d0-9813-a16d6f30f925',40,1,'2025-04-01 10:05:23'),(51,'23467e97-565c-484c-a34b-bbd2849f2f21',40,1,'2025-04-01 10:05:23'),(52,'23467e97-565c-484c-a34b-bbd2849f2f21',44,1,'2025-02-23 10:05:23'),(53,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1',45,-1,'2025-03-28 10:05:23'),(54,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927',45,-1,'2025-03-28 10:05:23'),(55,'ce973ab7-dfb0-41d0-9813-a16d6f30f925',46,1,'2025-03-19 10:05:23'),(56,'440a6b09-b04b-4959-b851-bb129ae70568',48,1,'2025-03-02 10:05:23'),(57,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9',49,-1,'2025-04-07 10:05:23'),(58,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9',49,1,'2025-04-07 10:05:23'),(59,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927',50,1,'2025-04-16 10:05:23'),(60,'359b955d-f28f-4a44-b759-60e3509010ff',50,-1,'2025-04-16 10:05:23');
/*!40000 ALTER TABLE `likes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notifications`
--

DROP TABLE IF EXISTS `notifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notifications` (
  `id` int NOT NULL AUTO_INCREMENT,
  `recipient_id` text NOT NULL,
  `sender_id` text NOT NULL,
  `type` varchar(50) NOT NULL,
  `post_id` int DEFAULT NULL,
  `comment_id` int DEFAULT NULL,
  `message` text NOT NULL,
  `seen` tinyint(1) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=94 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notifications`
--

LOCK TABLES `notifications` WRITE;
/*!40000 ALTER TABLE `notifications` DISABLE KEYS */;
/*!40000 ALTER TABLE `notifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post_category`
--

DROP TABLE IF EXISTS `post_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_category` (
  `post_id` int DEFAULT NULL,
  `category_id` int DEFAULT NULL,
  KEY `post_id` (`post_id`),
  KEY `category_id` (`category_id`),
  CONSTRAINT `post_category_ibfk_1` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
  CONSTRAINT `post_category_ibfk_2` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post_category`
--

LOCK TABLES `post_category` WRITE;
/*!40000 ALTER TABLE `post_category` DISABLE KEYS */;
INSERT INTO `post_category` VALUES (1,4),(1,1),(2,5),(2,3),(2,1),(3,3),(3,5),(4,1),(4,2),(4,4),(5,5),(5,2),(5,4),(6,2),(6,5),(6,1),(7,3),(7,5),(8,5),(8,3),(8,4),(9,4),(9,3),(10,5),(11,1),(12,5),(12,3),(13,3),(13,1),(13,2),(14,3),(14,4),(15,4),(15,2),(16,2),(16,4),(17,5),(17,1),(17,2),(18,2),(19,3),(20,4),(20,5),(20,1),(21,2),(21,4),(22,2),(22,3),(23,2),(24,2),(25,1),(25,4),(26,1),(26,5),(27,1),(27,3),(28,3),(28,1),(29,3),(29,5),(29,4),(30,3),(30,2),(31,1),(31,4),(31,3),(32,1),(33,1),(34,1),(35,5),(36,1),(37,5),(37,3),(38,3),(39,2),(40,3),(40,5),(41,4),(42,2),(43,1),(43,5),(44,2),(44,4),(45,4),(45,1),(46,5),(46,2),(47,1),(47,4),(47,3),(48,3),(48,4),(49,1),(49,5),(49,4),(50,1),(50,3),(50,2);
/*!40000 ALTER TABLE `post_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `posts`
--

DROP TABLE IF EXISTS `posts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `posts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(36) NOT NULL,
  `content` text NOT NULL,
  `image` longblob,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `posts_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=68 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `posts`
--

LOCK TABLES `posts` WRITE;
/*!40000 ALTER TABLE `posts` DISABLE KEYS */;
INSERT INTO `posts` VALUES (1,'440a6b09-b04b-4959-b851-bb129ae70568','Top 5 des films cultes à voir absolument. Je viens de revoir ce film, et je suis encore bluffé par la réalisation.',NULL,'2025-03-25 10:05:23'),(2,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Petit retour sur La Ligne Verte – toujours aussi poignant. Un film à voir au moins une fois dans sa vie.',NULL,'2025-04-08 10:05:23'),(3,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Top 5 des films cultes à voir absolument. Un film à voir au moins une fois dans sa vie.',NULL,'2025-02-18 10:05:23'),(4,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Pourquoi Le Parrain reste un chef-dœuvre du cinéma. C’est le genre de film qu’on n’oublie pas facilement.',NULL,'2025-03-31 10:05:23'),(5,'359b955d-f28f-4a44-b759-60e3509010ff','Ce film m’a bouleversé : discussion autour de Requiem for a Dream. C’est le genre de film qu’on n’oublie pas facilement.',NULL,'2025-02-18 10:05:23'),(6,'308e9405-784d-48bc-ab7f-393c25778f43','Quels films vous ont marqué dans votre enfance ? Certains passages traînent en longueur, mais le fond est fort.',NULL,'2025-03-24 10:05:23'),(7,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Ce film m’a bouleversé : discussion autour de Requiem for a Dream. L’histoire est bien transposée, les acteurs sont incroyables.',NULL,'2025-04-04 10:05:23'),(8,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Pourquoi Le Parrain reste un chef-dœuvre du cinéma. L’histoire est bien transposée, les acteurs sont incroyables.',NULL,'2025-04-04 10:05:23'),(9,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Pourquoi j’adore les thrillers psychologiques. Un scénario solide, des personnages profonds.',NULL,'2025-03-14 10:05:23'),(10,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Mon analyse de Pulp Fiction, 30 ans plus tard. Certains passages traînent en longueur, mais le fond est fort.',NULL,'2025-02-18 10:05:23'),(11,'df2c278e-f920-473e-b93d-ca31c9641e07','Petit retour sur La Ligne Verte – toujours aussi poignant. Un scénario solide, des personnages profonds.',NULL,'2025-04-11 10:05:23'),(12,'308e9405-784d-48bc-ab7f-393c25778f43','J’ai découvert un classique : 12 Hommes en colère. Une adaptation fidèle et ambitieuse.',NULL,'2025-04-17 10:05:23'),(13,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Petit retour sur La Ligne Verte – toujours aussi poignant. C’est le genre de film qu’on n’oublie pas facilement.',NULL,'2025-04-17 10:05:23'),(14,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Les meilleures adaptations de livres en films selon moi. Certains passages traînent en longueur, mais le fond est fort.',NULL,'2025-03-04 10:05:23'),(15,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','J’ai découvert un classique : 12 Hommes en colère. C’est le genre de film qu’on n’oublie pas facilement.',NULL,'2025-02-20 10:05:23'),(16,'df2c278e-f920-473e-b93d-ca31c9641e07','Les meilleures adaptations de livres en films selon moi. J’ai été ému par la simplicité et la puissance du message.',NULL,'2025-03-14 10:05:23'),(17,'308e9405-784d-48bc-ab7f-393c25778f43','Top 5 des films cultes à voir absolument. Certains passages traînent en longueur, mais le fond est fort.',NULL,'2025-03-01 10:05:23'),(18,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Quels films vous ont marqué dans votre enfance ? Un film à voir au moins une fois dans sa vie.',NULL,'2025-04-12 10:05:23'),(19,'df2c278e-f920-473e-b93d-ca31c9641e07','Quels films vous ont marqué dans votre enfance ? Un film que je recommande chaudement.',NULL,'2025-03-25 10:05:23'),(20,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Pourquoi Le Parrain reste un chef-dœuvre du cinéma. Un film à voir au moins une fois dans sa vie.',NULL,'2025-03-06 10:05:23'),(21,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Mon analyse de Pulp Fiction, 30 ans plus tard. Je viens de revoir ce film, et je suis encore bluffé par la réalisation.',NULL,'2025-03-08 10:05:23'),(22,'359b955d-f28f-4a44-b759-60e3509010ff','Ce film m’a bouleversé : discussion autour de Requiem for a Dream. Certains passages traînent en longueur, mais le fond est fort.',NULL,'2025-03-20 10:05:23'),(23,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Mon avis sur le film Inception après re-visionnage. L’histoire est bien transposée, les acteurs sont incroyables.',NULL,'2025-04-03 10:05:23'),(24,'308e9405-784d-48bc-ab7f-393c25778f43','Pourquoi Le Parrain reste un chef-dœuvre du cinéma. J’ai quelques réserves, mais dans l’ensemble c’est très réussi.',NULL,'2025-03-12 10:05:23'),(25,'23467e97-565c-484c-a34b-bbd2849f2f21','Quels films vous ont marqué dans votre enfance ? Certains passages traînent en longueur, mais le fond est fort.',NULL,'2025-02-27 10:05:23'),(26,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Ce film m’a bouleversé : discussion autour de Requiem for a Dream. Une adaptation fidèle et ambitieuse.',NULL,'2025-03-04 10:05:23'),(27,'23467e97-565c-484c-a34b-bbd2849f2f21','Pourquoi j’adore les thrillers psychologiques. Un scénario solide, des personnages profonds.',NULL,'2025-03-24 10:05:23'),(28,'df2c278e-f920-473e-b93d-ca31c9641e07','Les meilleures adaptations de livres en films selon moi. C’est le genre de film qu’on n’oublie pas facilement.',NULL,'2025-02-22 10:05:23'),(29,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Pourquoi Le Parrain reste un chef-dœuvre du cinéma. Certains passages traînent en longueur, mais le fond est fort.',NULL,'2025-02-19 10:05:23'),(30,'440a6b09-b04b-4959-b851-bb129ae70568','Les meilleures adaptations de livres en films selon moi. Certains passages traînent en longueur, mais le fond est fort.',NULL,'2025-03-01 10:05:23'),(31,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Top 5 des films cultes à voir absolument. L’histoire est bien transposée, les acteurs sont incroyables.',NULL,'2025-03-02 10:05:23'),(32,'df2c278e-f920-473e-b93d-ca31c9641e07','Top 5 des films cultes à voir absolument. Un film à voir au moins une fois dans sa vie.',NULL,'2025-04-05 10:05:23'),(33,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Top 5 des films cultes à voir absolument. J’ai été ému par la simplicité et la puissance du message.',NULL,'2025-03-09 10:05:23'),(34,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Mon analyse de Pulp Fiction, 30 ans plus tard. Je viens de revoir ce film, et je suis encore bluffé par la réalisation.',NULL,'2025-03-30 10:05:23'),(35,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Pourquoi Le Parrain reste un chef-dœuvre du cinéma. Certains passages traînent en longueur, mais le fond est fort.',NULL,'2025-03-02 10:05:23'),(36,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Mon analyse de Pulp Fiction, 30 ans plus tard. L’histoire est bien transposée, les acteurs sont incroyables.',NULL,'2025-03-18 10:05:23'),(37,'359b955d-f28f-4a44-b759-60e3509010ff','Mon avis sur le film Inception après re-visionnage. J’ai quelques réserves, mais dans l’ensemble c’est très réussi.',NULL,'2025-03-04 10:05:23'),(38,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','J’ai découvert un classique : 12 Hommes en colère. C’est le genre de film qu’on n’oublie pas facilement.',NULL,'2025-03-28 10:05:23'),(39,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Mon analyse de Pulp Fiction, 30 ans plus tard. Un film que je recommande chaudement.',NULL,'2025-02-28 10:05:23'),(40,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Quels films vous ont marqué dans votre enfance ? J’ai quelques réserves, mais dans l’ensemble c’est très réussi.',NULL,'2025-04-01 10:05:23'),(41,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Petit retour sur La Ligne Verte – toujours aussi poignant. J’ai quelques réserves, mais dans l’ensemble c’est très réussi.',NULL,'2025-03-28 10:05:23'),(42,'df2c278e-f920-473e-b93d-ca31c9641e07','Ce film m’a bouleversé : discussion autour de Requiem for a Dream. L’histoire est bien transposée, les acteurs sont incroyables.',NULL,'2025-03-13 10:05:23'),(43,'91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Mon analyse de Pulp Fiction, 30 ans plus tard. Un film à voir au moins une fois dans sa vie.',NULL,'2025-04-05 10:05:23'),(44,'359b955d-f28f-4a44-b759-60e3509010ff','Pourquoi Le Parrain reste un chef-dœuvre du cinéma. Un scénario solide, des personnages profonds.',NULL,'2025-02-23 10:05:23'),(45,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Mon analyse de Pulp Fiction, 30 ans plus tard. L’histoire est bien transposée, les acteurs sont incroyables.',NULL,'2025-03-28 10:05:23'),(46,'308e9405-784d-48bc-ab7f-393c25778f43','J’ai découvert un classique : 12 Hommes en colère. J’ai quelques réserves, mais dans l’ensemble c’est très réussi.',NULL,'2025-03-19 10:05:23'),(47,'f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','J’ai découvert un classique : 12 Hommes en colère. Un film que je recommande chaudement.',NULL,'2025-03-15 10:05:23'),(48,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Petit retour sur La Ligne Verte – toujours aussi poignant. Une adaptation fidèle et ambitieuse.',NULL,'2025-03-02 10:05:23'),(49,'23467e97-565c-484c-a34b-bbd2849f2f21','Pourquoi j’adore les thrillers psychologiques. J’ai quelques réserves, mais dans l’ensemble c’est très réussi.',NULL,'2025-04-07 10:05:23'),(50,'308e9405-784d-48bc-ab7f-393c25778f43','Les meilleures adaptations de livres en films selon moi. C’est le genre de film qu’on n’oublie pas facilement.',NULL,'2025-04-16 10:05:23');
/*!40000 ALTER TABLE `posts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `reports`
--

DROP TABLE IF EXISTS `reports`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `reports` (
  `id` int NOT NULL AUTO_INCREMENT,
  `post_id` int NOT NULL,
  `reporter_id` varchar(255) NOT NULL,
  `reason` text,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `post_id` (`post_id`),
  KEY `reporter_id` (`reporter_id`),
  CONSTRAINT `reports_ibfk_1` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
  CONSTRAINT `reports_ibfk_2` FOREIGN KEY (`reporter_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `reports`
--

LOCK TABLES `reports` WRITE;
/*!40000 ALTER TABLE `reports` DISABLE KEYS */;
INSERT INTO `reports` VALUES (47,5,'b7ef6093-5bd1-4cca-b1cf-b578cd569221','Contenu offensant','2025-02-18 10:05:23'),(48,10,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Spam','2025-02-18 10:05:23'),(49,15,'df2c278e-f920-473e-b93d-ca31c9641e07','Contenu offensant','2025-02-20 10:05:23'),(50,20,'440a6b09-b04b-4959-b851-bb129ae70568','Hors sujet','2025-03-06 10:05:23'),(51,25,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Hors sujet','2025-02-27 10:05:23'),(52,30,'df2c278e-f920-473e-b93d-ca31c9641e07','Spam','2025-03-01 10:05:23'),(53,35,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Spam','2025-03-02 10:05:23'),(54,40,'ce973ab7-dfb0-41d0-9813-a16d6f30f925','Hors sujet','2025-04-01 10:05:23'),(55,45,'df2c278e-f920-473e-b93d-ca31c9641e07','Contenu offensant','2025-03-28 10:05:23'),(56,50,'a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Hors sujet','2025-04-16 10:05:23');
/*!40000 ALTER TABLE `reports` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sessions`
--

DROP TABLE IF EXISTS `sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sessions` (
  `token` varchar(255) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `user_agent` text NOT NULL,
  `expires_at` datetime NOT NULL,
  `role` enum('user','admin','moderator') NOT NULL DEFAULT 'user',
  PRIMARY KEY (`token`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `sessions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sessions`
--

LOCK TABLES `sessions` WRITE;
/*!40000 ALTER TABLE `sessions` DISABLE KEYS */;
/*!40000 ALTER TABLE `sessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` varchar(36) NOT NULL,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `password` varchar(255) DEFAULT '',
  `role` varchar(50) NOT NULL DEFAULT 'user',
  `google_id` varchar(255) DEFAULT NULL,
  `provider` varchar(50) DEFAULT 'local',
  `github_id` varchar(255) DEFAULT NULL,
  `banned` tinyint(1) NOT NULL DEFAULT '0',
  `is_public` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `google_id` (`google_id`),
  UNIQUE KEY `github_id` (`github_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES ('090b4790-644d-4d5a-92e1-64001583b0ad','moderateur','moderateur@gmail.com','$2a$10$EpjKvEkekXxvBXhNjSKuLu6dTQL.eTwqQdsB6EI7vK5XH98wGm/Um','moderator',NULL,'local',NULL,0,0),('23467e97-565c-484c-a34b-bbd2849f2f21','David','david@example.com','$2b$12$31nbGEfHa9EdigxQ13V25e5PWM7zGpQ4b/deAH51QYyJPTN.xZDYO','user',NULL,'local',NULL,0,1),('308e9405-784d-48bc-ab7f-393c25778f43','Alice','alice@example.com','$2b$12$P7pLhkkSya7YdUTQC6sEP.Si5DZa1XsPj12O2XkUjvc5qFPh2flzi','user',NULL,'local',NULL,0,1),('359b955d-f28f-4a44-b759-60e3509010ff','Jack','jack@example.com','$2b$12$3cubDuWE8kPpFoUgoaxcvuixdKbYJDYzUKYvTSfAY/yZwTh/2DMby','user',NULL,'local',NULL,0,1),('440a6b09-b04b-4959-b851-bb129ae70568','Emma','emma@example.com','$2b$12$sawQaBW/8xpn2uOv8CGRgOVv3LabKsToK5FXtbWOl.xZwhvZ5ZM76','user',NULL,'local',NULL,0,1),('64f07f12-8177-4d08-81d5-ef774a6fb551','user','user@gmail.com','$2a$10$NDlT3YoeRMI6WTQkkRX79O45rh9SwLdUWJBInvGZhNOiF98GHAcWm','user',NULL,'local',NULL,0,0),('91d6b5a2-0d3e-4adf-89d5-95a5e3f7c927','Bob','bob@example.com','$2b$12$vHsnH.7hnqKhIkw8C54QkeoHqLKdza3rY3c3l4ioRxUQRL/wZC8fS','user',NULL,'local',NULL,0,1),('a4ffc122-0ca9-4c84-9b9b-0cc12129caa9','Hugo','hugo@example.com','$2b$12$Q5qCYc12LBidPwgoRpMpBOIS4jirlOQDAzzK8Wp.dibDcd0/1qHJi','user',NULL,'local',NULL,0,1),('b7ef6093-5bd1-4cca-b1cf-b578cd569221','Ivy','ivy@example.com','$2b$12$V1OhHz4gzTRDhsmJ2qrEBOajuhbA3tOg2.uWE0sdF3b8UZbYA1FHq','user',NULL,'local',NULL,0,1),('ce973ab7-dfb0-41d0-9813-a16d6f30f925','Grace','grace@example.com','$2b$12$sVBCKcS1dq4muOLCLfOxbef1eo/3yrz9SFsVJ2wVP9HtzYIkclUR2','user',NULL,'local',NULL,0,1),('df2c278e-f920-473e-b93d-ca31c9641e07','Clara','clara@example.com','$2b$12$zhCre2lBIhwOC1uVHwUwCuRI0mGrDjzvuX.l8o7mWJ6mftOUVK1Vm','user',NULL,'local',NULL,0,1),('f5f5fba4-6462-4f07-8e5a-b0c8ca2e84a1','Frank','frank@example.com','$2b$12$/zE4ZMr/bE5LUcZneX7XteLKKee/WBh/KPQxT5MeySC3tHIPhxhBO','user',NULL,'local',NULL,0,1),('fc35f47c-633b-4438-9ed1-248113422950','admin','admin@gmail.com','$2a$10$SZtPcSeo5wPOlulcP4t0F.M285I0spcTM1zA24OmTbFLqkGAt038G','admin',NULL,'local',NULL,0,0);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `warns`
--

DROP TABLE IF EXISTS `warns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `warns` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(36) NOT NULL,
  `issued_by` varchar(36) NOT NULL,
  `reason` text NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `issued_by` (`issued_by`),
  CONSTRAINT `warns_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `warns_ibfk_2` FOREIGN KEY (`issued_by`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `warns`
--

LOCK TABLES `warns` WRITE;
/*!40000 ALTER TABLE `warns` DISABLE KEYS */;
/*!40000 ALTER TABLE `warns` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-04-17 10:38:37

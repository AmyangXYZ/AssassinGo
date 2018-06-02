-- MySQL dump 10.16  Distrib 10.2.13-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: ag
-- ------------------------------------------------------
-- Server version	10.2.13-MariaDB-10.2.13+maria~jessie

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `username` varchar(16) DEFAULT NULL,
  `password` varchar(64) DEFAULT NULL,
  `email` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=71 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'admin','adminn','a@a.com'),(2,'admin','adminn','a@a.com'),(3,'admin','adminn','a@a.com'),(4,'admin','adminn','a@a.com'),(5,'admin','adminn','a@a.com'),(6,'admin','adminn','a@a.com'),(7,'admin','adminn','a@a.com'),(8,'admin','adminn','a@a.com'),(9,'admin','adminn','a@a.com'),(10,'admin','adminn','a@a.com'),(11,'admin','adminn','a@a.com'),(12,'admin','adminn','a@a.com'),(13,'admin','adminn','a@a.com'),(14,'admin','adminn','a@a.com'),(15,'admin','adminn','a@a.com'),(16,'admin','adminn','a@a.com'),(17,'admin','adminn','a@a.com'),(18,'admin','adminn','a@a.com'),(19,'admin','adminn','a@a.com'),(20,'admin','adminn','a@a.com'),(21,'admin','adminn','a@a.com'),(22,'admin','adminn','a@a.com'),(23,'admin','adminn','a@a.com'),(24,'admin','adminn','a@a.com'),(25,'admin','adminn','a@a.com'),(26,'admin','adminn','a@a.com'),(27,'admin','adminn','a@a.com'),(28,'admin','adminn','a@a.com'),(29,'admin','adminn','a@a.com'),(30,'admin','adminn','a@a.com'),(31,'admin','adminn','a@a.com'),(32,'admin','adminn','a@a.com'),(33,'admin','adminn','a@a.com'),(34,'admin','adminn','a@a.com'),(35,'admin','adminn','a@a.com'),(36,'admin','adminn','a@a.com'),(37,'admin','adminn','a@a.com'),(38,'admin','adminn','a@a.com'),(39,'admin','adminn','a@a.com'),(40,'admin','adminn','a@a.com'),(41,'admin','adminn','a@a.com'),(42,'admin','adminn','a@a.com'),(43,'admin','adminn','a@a.com'),(44,'admin','adminn','a@a.com'),(45,'admin','adminn','a@a.com'),(46,'admin','adminn','a@a.com'),(47,'admin','adminn','a@a.com'),(48,'admin','adminn','a@a.com'),(49,'admin','adminn','a@a.com'),(50,'admin','adminn','a@a.com'),(51,'admin','adminn','a@a.com'),(52,'admin','adminn','a@a.com'),(53,'admin','adminn','a@a.com'),(54,'admin','adminn','a@a.com'),(55,'admin','adminn','a@a.com'),(56,'admin','adminn','a@a.com'),(57,'admin','adminn','a@a.com'),(58,'admin','adminn','a@a.com'),(59,'admin','adminn','a@a.com'),(60,'admin','adminn','a@a.com'),(61,'admin','adminn','a@a.com'),(62,'admin','adminn','a@a.com'),(63,'admin','adminn','a@a.com'),(64,'admin','adminn','a@a.com'),(65,'admin','adminn','a@a.com'),(66,'admin','adminn','a@a.com'),(67,'admin','adminn','a@a.com'),(68,'admin','adminn','a@a.com'),(69,'admin','adminn',NULL),(70,'admin','admin',NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-06-02  1:32:40

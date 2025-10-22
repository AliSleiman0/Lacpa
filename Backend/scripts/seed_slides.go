package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AliSleiman0/Lacpa/config"
	"github.com/AliSleiman0/Lacpa/models/admin"
	adminRepo "github.com/AliSleiman0/Lacpa/repository/admin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := config.ConnectMongoDB(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Initialize repository
	database := mongoClient.Database("LACPA")
	repo := adminRepo.NewHeroSlideRepository(database)

	// Clear existing slides (optional - comment out if you want to keep existing data)
	fmt.Println("Clearing existing slides...")
	collection := database.Collection("hero_slides")
	_, err = collection.DeleteMany(context.Background(), map[string]interface{}{})
	if err != nil {
		log.Printf("Warning: Could not clear existing slides: %v\n", err)
	}

	// Seed data
	slides := []admin.HeroSlide{
		{
			Title:             "Welcome to LACPA",
			Description:       "The Lebanese Association of Certified Public Accountants (LACPA) is the premier professional body for accountants in Lebanon, dedicated to upholding the highest standards of excellence and ethics.",
			ImgSrc:            "img_c1.png",
			ButtonTitle:       "Learn More",
			ButtonLink:        "/about",
			IsActive:          true,
			ImageActive:       true,
			ButtonActive:      true,
			TitleActive:       true,
			DescriptionActive: true,
			OrderIndex:        1,
		},
		{
			Title:             "Professional Excellence",
			Description:       "Join a community of over 1,500 certified professionals committed to maintaining the highest standards of accounting practice and professional development.",
			ImgSrc:            "img_c2.png",
			ButtonTitle:       "Become a Member",
			ButtonLink:        "/membership",
			IsActive:          true,
			ImageActive:       true,
			ButtonActive:      true,
			TitleActive:       true,
			DescriptionActive: true,
			OrderIndex:        2,
		},
		{
			Title:             "Continuous Learning",
			Description:       "Access exclusive training programs, workshops, and certifications designed to keep you at the forefront of the accounting profession.",
			ImgSrc:            "img_c3.png",
			ButtonTitle:       "View Programs",
			ButtonLink:        "/events",
			IsActive:          true,
			ImageActive:       true,
			ButtonActive:      true,
			TitleActive:       true,
			DescriptionActive: true,
			OrderIndex:        3,
		},
		{
			Title:             "Industry Leadership",
			Description:       "Stay informed with the latest industry news, regulations, and best practices through our comprehensive resources and publications.",
			ImgSrc:            "background.png",
			ButtonTitle:       "Read Publications",
			ButtonLink:        "/publications",
			IsActive:          true,
			ImageActive:       true,
			ButtonActive:      true,
			TitleActive:       true,
			DescriptionActive: true,
			OrderIndex:        4,
		},
		{
			Title:             "Network & Connect",
			Description:       "Build valuable connections with fellow professionals through our networking events, conferences, and online community platform.",
			ImgSrc:            "logo.png",
			ButtonTitle:       "Join Events",
			ButtonLink:        "/events",
			IsActive:          true,
			ImageActive:       true,
			ButtonActive:      true,
			TitleActive:       true,
			DescriptionActive: true,
			OrderIndex:        5,
		},
	}

	// Insert slides
	fmt.Println("Seeding hero slides...")
	for i, slide := range slides {
		if err := repo.CreateSlide(context.Background(), &slide); err != nil {
			log.Printf("Failed to create slide %d: %v\n", i+1, err)
			continue
		}
		fmt.Printf("✓ Created slide %d: %s\n", i+1, slide.Title)
	}

	fmt.Println("\n✅ Successfully seeded", len(slides), "hero slides!")
	fmt.Println("You can now view them at: http://localhost:3000/admin")
}

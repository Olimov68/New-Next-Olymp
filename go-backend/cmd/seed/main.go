package main

import (
	"fmt"
	"log"

	"github.com/nextolympservice/go-backend/config"
	"github.com/nextolympservice/go-backend/internal/database"
	"github.com/nextolympservice/go-backend/internal/models"
	"github.com/nextolympservice/go-backend/internal/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config load error:", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("DB connect error:", err)
	}

	// Check if superadmin already exists
	var count int64
	db.Model(&models.StaffUser{}).Where("role = ?", "superadmin").Count(&count)
	if count > 0 {
		fmt.Println("SuperAdmin already exists, skipping seed")
		// Show existing superadmins
		var admins []models.StaffUser
		db.Where("role = ?", "superadmin").Find(&admins)
		for _, a := range admins {
			fmt.Printf("  - ID: %d, Username: %s, FullName: %s, Status: %s\n", a.ID, a.Username, a.FullName, a.Status)
		}
		return
	}

	hash, err := utils.HashPassword("SuperAdmin123!")
	if err != nil {
		log.Fatal("Password hash error:", err)
	}

	superadmin := &models.StaffUser{
		Username:     "superadmin",
		PasswordHash: hash,
		FullName:     "Super Administrator",
		Role:         models.StaffRoleSuperAdmin,
		Status:       models.StaffStatusActive,
	}

	if err := db.Create(superadmin).Error; err != nil {
		log.Fatal("Create error:", err)
	}

	fmt.Println("SuperAdmin user created successfully!")
	fmt.Println("  Username: superadmin")
	fmt.Println("  Password: SuperAdmin123!")
	fmt.Println("  Role: superadmin")
	fmt.Printf("  ID: %d\n", superadmin.ID)
}

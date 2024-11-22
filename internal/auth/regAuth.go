package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/EmptyInsid/db_gui/internal/database"
)

func RegistrUser(db *database.Database, ctx context.Context, username, password, role string) error {
	passwordHash, err := HashPassword(password)
	if err != nil {
		log.Printf("Error while hash password: %v\n", err)
		return err
	}

	if err := db.RegistrUserDB(ctx, username, passwordHash, role); err != nil {
		log.Printf("Error while reqistry user: %v\n", err)
		return err
	}

	return err
}

func AuthenticateUser(db *database.Database, ctx context.Context, username, password string) (string, string, error) {
	storedPassword, role, err := db.AuthUser(ctx, username, password)
	if err != nil {
		log.Printf("Error while auth user: %v\n", err)
		return "", "", err
	}
	if !CheckPasswordHash(password, storedPassword) {
		log.Printf("Incorrect password")
		return "", "", fmt.Errorf("incorrect password")
	}
	return username, role, nil
}

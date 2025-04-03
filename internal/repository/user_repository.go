package repository

import (
	_ "errors"
	_ "fmt"
	_ "time"

	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"gorm.io/gorm"
)

// UserRepository interface defines methods for user data access
type UserRepository interface {
	Create(user *model.User) error
	// GetByID(id int64) (*model.User, error)
	// GetAll() ([]*model.User, error)
	// Update(id int64, user *model.User) error
	// Delete(id int64) error
}

// userRepository implements UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create inserts a new user into the database
func (r *userRepository) Create(user *model.User) error {
	result := r.db.Create(&user)
	return result.Error
}

// GetByID retrieves a user by ID
// func (r *userRepository) GetByID(id int64) (*model.User, error) {
// 	query := `
// 		SELECT id, username, email, phone, created_at, updated_at
// 		FROM users
// 		WHERE id = $1`

// 	var user model.User
// 	err := r.db.Get(&user, query, id)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get user: %w", err)
// 	}

// 	return &user, nil
// }

// GetAll retrieves all users
// func (r *userRepository) GetAll() ([]*model.User, error) {
// 	query := `
// 		SELECT id, username, email, phone, created_at, updated_at
// 		FROM users
// 		ORDER BY id`

// 	var users []*model.User
// 	err := r.db.Select(&users, query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get users: %w", err)
// 	}

// 	return users, nil
// }

// Update updates an existing user
// func (r *userRepository) Update(id int64, user *model.User) error {
// 	query := `
// 		UPDATE users
// 		SET username = $1, email = $2, phone = $3, updated_at = $4
// 		WHERE id = $5`

// 	user.UpdatedAt = time.Now()

// 	result, err := r.db.Exec(
// 		query,
// 		user.Username,
// 		user.Email,
// 		user.Phone,
// 		user.UpdatedAt,
// 		id,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to update user: %w", err)
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get rows affected: %w", err)
// 	}

// 	if rowsAffected == 0 {
// 		return errors.New("user not found")
// 	}

// 	user.ID = id
// 	return nil
// }

// Delete removes a user by ID
// func (r *userRepository) Delete(id int64) error {
// 	query := `DELETE FROM users WHERE id = $1`

// 	result, err := r.db.Exec(query, id)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete user: %w", err)
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get rows affected: %w", err)
// 	}

// 	if rowsAffected == 0 {
// 		return errors.New("user not found")
// 	}

// 	return nil
// }

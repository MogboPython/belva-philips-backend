package repository

import (
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"gorm.io/gorm"
)

// UserRepository interface defines methods for user data access
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetAll(page, limit int) ([]*model.User, error)
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
	err := r.db.Create(&user).Error
	// if errors.Is(err, gorm.ErrDuplicatedKey) {
	// 	return fmt.Errorf("user with this email exists: %w", err)
	// }
	return err
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id string) (*model.User, error) {
	var user model.User

	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User

	if err := r.db.Where(&model.User{Email: email}).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// func getUserByEmail(e string) (*model.User, error) {
// 	db := database.DB
// 	var user model.User
// 	if err := db.Where(&model.User{Email: e}).First(&user).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &user, nil
// }

// GetAll retrieves all users
func (r *userRepository) GetAll(offset, limit int) ([]*model.User, error) {
	var users []*model.User

	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

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

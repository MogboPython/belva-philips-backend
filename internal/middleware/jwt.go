package middleware

import (
	"github.com/MogboPython/belvaphilips_backend/internal/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.Config("JWT_SECRET"))},
		ErrorHandler: jwtError,
	})
}

func AdminProtected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.Config("ADMIN_JWT_SECRET"))},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

// import (
// 	"github.com/google/uuid"
// 	"github.com/lestrrat-go/jwx/jwa"
// 	"github.com/lestrrat-go/jwx/jwt"
// 	"github.com/uptrace/bun"
// )

// type UserMetadata struct {
// 	StripeCustomerID string `json:"stripe_customer_id"`
// }

// type User struct {
// 	bun.BaseModel
// 	ID        uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
// 	Reference string    `json:"reference"`
// 	Email     string    `json:"email"`

// 	FullName string        `json:"full_name"`
// 	Metadata *UserMetadata `json:"metadata"`

// 	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
// 	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
// 	DeletedAt time.Time `bun:",soft_delete,nullzero"`
// }

// func ParseUserFromJWT() (*User, error) {

// 	token, err := jwt.Parse([]byte(os.Getenv("JWT_SECRET_KEY")), jwt.WithVerify(jwa.HS256, []byte(config.Global().Supabase.JWTKey)))
// 	if err != nil {
// 		return nil, err
// 	}

// 	if token.Expiration().Before(time.Now()) {
// 		return nil, errors.New("token is expired")
// 	}

// 	err = errors.New("not found")

// 	id, exists := token.Get("sub")
// 	if !exists {
// 		return nil, err
// 	}

// 	email, exists := token.Get("email")
// 	if !exists {
// 		return nil, err
// 	}

// 	userMetada, exists := token.Get("user_metadata")
// 	if !exists {
// 		return nil, err
// 	}

// 	data, ok := userMetada.(map[string]interface{})
// 	if !ok {
// 		return nil, errors.New("invalid jwt")
// 	}

// 	return &User{
// 		Reference: id.(string),
// 		FullName:  data["full_name"].(string),
// 		Email:     Email(email.(string)),
// 		Metadata:  &UserMetadata{},
// 	}, nil
// }

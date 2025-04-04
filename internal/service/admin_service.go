package service

// Login get user and password
// func Login(c *fiber.Ctx) error {
// 	type LoginInput struct {
// 		Identity string `json:"identity"`
// 		Password string `json:"password"`
// 	}
// 	type UserData struct {
// 		ID       uint   `json:"id"`
// 		Username string `json:"username"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}
// 	input := new(LoginInput)
// 	var userData UserData

// 	if err := c.BodyParser(&input); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
// 	}

// 	identity := input.Identity
// 	pass := input.Password
// 	userModel, err := new(model.User), *new(error)

// 	if isEmail(identity) {
// 		userModel, err = getUserByEmail(identity)
// 	} else {
// 		userModel, err = getUserByUsername(identity)
// 	}

// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err})
// 	} else if userModel == nil {
// 		CheckPasswordHash(pass, "")
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": err})
// 	} else {
// 		userData = UserData{
// 			ID:       userModel.ID,
// 			Username: userModel.Username,
// 			Email:    userModel.Email,
// 			Password: userModel.Password,
// 		}
// 	}

// 	if !CheckPasswordHash(pass, userData.Password) {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": nil})
// 	}

// 	token := jwt.New(jwt.SigningMethodHS256)

// 	claims := token.Claims.(jwt.MapClaims)
// 	claims["username"] = userData.Username
// 	claims["user_id"] = userData.ID
// 	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

// 	t, err := token.SignedString([]byte(config.Config("SECRET")))
// 	if err != nil {
// 		return c.SendStatus(fiber.StatusInternalServerError)
// 	}

// 	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
// }

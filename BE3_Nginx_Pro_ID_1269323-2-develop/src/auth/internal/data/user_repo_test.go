package data

import (
	"auth/internal/auth"
	"context"
	"errors"
	"log"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func TestRepo(t *testing.T) {
	users := []*auth.User{
		&auth.User{
			Id:           uuid.New(),
			Login:        "test",
			Name:         "test",
			Surname:      "testov",
			Phonenumber:  "+78005553535",
			PasswordHash: make([]byte, 60),
			Email:        "test@example.com",
			Role:         "admin",
		},
		&auth.User{
			Id:           uuid.New(),
			Login:        "тест",
			Name:         "test",
			Surname:      "testov",
			Phonenumber:  "+78005553535",
			PasswordHash: make([]byte, 60),
			Email:        "test@example.com",
			Role:         "admin",
		},
	}

	pool, err := pgxpool.New(context.Background(), "postgresql://postgres:postgres@localhost:5432/shop_auth")
	defer pool.Close()
	if err != nil {
		log.Fatal(err)
	}
	if err = pool.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}
	repo := NewUserRepo(pool)

	t.Run("SaveUser", func(t *testing.T) {
		correct := []*auth.User{users[0], users[1]}

		incorrect := []*auth.User{
			// Nil UUID
			&auth.User{
				Id:           uuid.Nil,
				Login:        "test1",
				Name:         "test",
				Surname:      "testov",
				Phonenumber:  "+78005553535",
				PasswordHash: make([]byte, 60),
				Email:        "test@example.com",
				Role:         "admin",
			},

			// Empty login
			&auth.User{
				Id:           uuid.New(),
				Login:        "",
				Name:         "test2",
				Surname:      "testov",
				Phonenumber:  "+78005553535",
				PasswordHash: make([]byte, 60),
				Email:        "test@example.com",
				Role:         "admin",
			},

			// Existent login
			&auth.User{
				Id:           uuid.New(),
				Login:        "test",
				Name:         "test",
				Surname:      "testov",
				Phonenumber:  "+78005553535",
				PasswordHash: make([]byte, 60),
				Email:        "test@example.com",
				Role:         "admin",
			},
			// Too large login
			&auth.User{
				Id:           uuid.New(),
				Login:        "dfkjsdlkjsklfjjfklsdjklsdjfklsdjf cdfksdjlsdfksdjflskdjfsdlkfjsdlkfjsdlkfjsdklfjkfjfjfjdkjfskldjfklsdhjfjfkjflskdfjlskdjflfjskldfjklsdjf",
				Name:         "test3",
				Surname:      "testov",
				Phonenumber:  "+78005553535",
				PasswordHash: make([]byte, 60),
				Email:        "test@example.com",
				Role:         "admin",
			},

			// Empty email
			&auth.User{
				Id:           uuid.New(),
				Login:        "test4",
				Name:         "test",
				Surname:      "testov",
				Phonenumber:  "+78005553535",
				PasswordHash: make([]byte, 60),
				Email:        "",
				Role:         "admin",
			},

			// Too short email
			&auth.User{
				Id:           uuid.New(),
				Login:        "test5",
				Name:         "test",
				Surname:      "testov",
				Phonenumber:  "+78005553535",
				PasswordHash: make([]byte, 60),
				Email:        "a",
				Role:         "admin",
			},

			// Too long email
			&auth.User{
				Id:           uuid.New(),
				Login:        "test6",
				Name:         "test",
				Surname:      "testov",
				Phonenumber:  "+78005553535",
				PasswordHash: make([]byte, 60),
				Email:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Role:         "admin",
			},

			// Nil password hash
			&auth.User{
				Id:           uuid.New(),
				Login:        "test7",
				Name:         "test",
				Surname:      "testov",
				Phonenumber:  "+78005553535",
				PasswordHash: nil,
				Email:        "test@example.com",
				Role:         "admin",
			},

			// Empty role
			&auth.User{
				Id:           uuid.New(),
				Login:        "test8",
				Name:         "test",
				Surname:      "testov",
				Phonenumber:  "+78005553535",
				PasswordHash: make([]byte, 60),
				Email:        "test@example.com",
				Role:         "",
			},

			// Nonexistent role
			&auth.User{
				Id:           uuid.New(),
				Login:        "test9",
				Name:         "test",
				Surname:      "testov",
				Phonenumber:  "+78005553535",
				PasswordHash: make([]byte, 60),
				Email:        "test@example.com",
				Role:         "god",
			},

			// incorrect role
			&auth.User{
				Id:           uuid.New(),
				Login:        "test10",
				Name:         "test",
				Surname:      "testov",
				Phonenumber:  "+78005553535",
				PasswordHash: make([]byte, 60),
				Email:        "test@example.com",
				Role:         "aDmIn",
			},
		}

		for _, test := range correct {
			if err = repo.SaveUser(context.Background(), test); err != nil {
				log.Fatalf("SaveUser: User: %+v, err: %v", test, err)
			}
		}

		for _, test := range incorrect {
			err = repo.SaveUser(context.Background(), test)
			if err == nil {
				log.Fatalf("SaveUser: User: %+v, missing error", test)
			}
		}
	})

	t.Run("GetUserById", func(t *testing.T) {
		unexistentId := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

		for _, user := range users {
			res, err := repo.GetUserById(context.Background(), user.Id)
			if err != nil {
				log.Fatalf("GetUserById: Id: %v, err: %v", user.Id, err)
			}
			if !compareUsers(user, res) {
				log.Fatalf("GetUserById: Want: %+v, get: %+v", user, res)
			}
		}

		for _, id := range unexistentId {
			res, err := repo.GetUserById(context.Background(), id)
			if !errors.Is(err, auth.ErrUserNotFound) {
				log.Fatalf("GetUserById: Missing ErrUserNotFound error, get result: %v, get error: %v", res, err)
			}
		}
	})

	t.Run("GetUserByLogin", func(t *testing.T) {
		unexistent := []string{"1", "sdfkjsdlkf", "aaaaa"}

		for _, user := range users {
			res, err := repo.GetUserByLogin(context.Background(), user.Login)
			if err != nil {
				log.Fatalf("GetUserByLogin: Login: %v, err: %v", user.Login, err)
			}
			if !compareUsers(user, res) {
				log.Fatalf("GetUserByLogin: Want: %+v, get: %+v", user, res)
			}
		}

		for _, login := range unexistent {
			res, err := repo.GetUserByLogin(context.Background(), login)
			if !errors.Is(err, auth.ErrUserNotFound) {
				log.Fatalf("GetUserByLogin: Missing ErrUserNotFound error, get result: %v, get error: %v", res, err)
			}
		}
	})

	t.Run("ChangeUserPassword", func(t *testing.T) {
		unexistentId := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

		for _, user := range users {
			new, _ := bcrypt.GenerateFromPassword([]byte("qwerty123"), 3)

			err := repo.ChangeUserPassword(context.Background(), user.Id, new)
			if err != nil {
				log.Fatalf("ChangeUserPassword: Id: %v, err: %v", user.Id, err)
			}
		}

		for _, id := range unexistentId {
			new, _ := bcrypt.GenerateFromPassword([]byte("qwerty123"), 3)

			err := repo.ChangeUserPassword(context.Background(), id, new)
			if !errors.Is(err, auth.ErrUserNotFound) {
				log.Fatalf("ChangeUserPassword: Want: ErrUserNotFound, get: %v", err)
			}
		}
	})

	for _, user := range users {
		repo.Delete(context.Background(), user.Id)
	}
}

func compareUsers(user1 *auth.User, user2 *auth.User) bool {
	return user1.Id == user2.Id &&
		user1.Login == user2.Login &&
		user1.Name == user2.Name &&
		user1.Surname == user2.Surname &&
		user1.Phonenumber == user2.Phonenumber &&
		slices.Equal(user1.PasswordHash, user2.PasswordHash) &&
		user1.Email == user2.Email &&
		user1.Role == user2.Role
}

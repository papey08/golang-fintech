package user_repo

import (
	"context"
	"homework8/internal/app"
	"homework8/internal/model/errs"
	"homework8/internal/model/users"
	"sync"
)

type UserRepository struct {
	data map[int64]users.User
	mu   sync.Mutex
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (users.User, error) {
	select {
	case <-ctx.Done():
		return users.User{}, errs.UserRepositoryError
	default:
		r.mu.Lock()
		defer r.mu.Unlock()
		if user, ok := r.data[id]; !ok {
			return users.User{}, errs.UserNotExist
		} else {
			return user, nil
		}
	}
}

func (r *UserRepository) AddUser(ctx context.Context, user users.User) (users.User, error) {
	select {
	case <-ctx.Done():
		return users.User{}, errs.UserRepositoryError
	default:
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.data[user.ID]; ok {
			return users.User{}, errs.UserAlreadyExists
		} else {
			r.data[user.ID] = user
			return user, nil
		}
	}
}

func (r *UserRepository) UpdateUserFields(ctx context.Context, idToUpdate int64, newUser users.User) (users.User, error) {
	select {
	case <-ctx.Done():
		return users.User{}, errs.UserRepositoryError
	default:
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.data[idToUpdate]; !ok {
			return users.User{}, errs.UserNotExist
		} else {
			user := users.User{
				ID:       idToUpdate,
				Nickname: newUser.Nickname,
				Email:    newUser.Email,
			}
			r.data[idToUpdate] = user
			return user, nil
		}
	}
}

func New() app.UserRepository {
	return &UserRepository{
		data: make(map[int64]users.User),
	}
}

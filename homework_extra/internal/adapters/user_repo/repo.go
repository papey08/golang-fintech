package user_repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"homework_extra/internal/app"
	"homework_extra/internal/model/errs"
	"homework_extra/internal/model/users"
	"sync"
)

type UserRepository struct {
	*mongo.Collection
	nextID int64
	mu     *sync.Mutex
}

type UsersField struct {
	ID       int64  `bson:"id"`
	Nickname string `bson:"nickname"`
	Email    string `bson:"email"`
}

func userToUsersField(u users.User) UsersField {
	return UsersField{
		ID:       u.ID,
		Nickname: u.Nickname,
		Email:    u.Email,
	}
}

func usersFieldToUser(u UsersField) users.User {
	return users.User{
		ID:       u.ID,
		Nickname: u.Nickname,
		Email:    u.Email,
	}
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (users.User, error) {
	var res UsersField
	filter := bson.M{"id": id}
	err := r.FindOne(ctx, filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return users.User{}, errs.UserNotExist
	} else if err != nil {
		return users.User{}, errs.UserRepositoryError
	} else {
		return usersFieldToUser(res), nil
	}
}

func (r *UserRepository) AddUser(ctx context.Context, user users.User) (users.User, error) {
	r.mu.Lock()
	r.nextID++
	r.mu.Unlock()
	u := userToUsersField(user)
	u.ID = r.nextID
	_, err := r.InsertOne(ctx, u)
	if err != nil {
		return users.User{}, errs.UserRepositoryError
	}
	return usersFieldToUser(u), nil
}

func (r *UserRepository) UpdateUserFields(ctx context.Context, idToUpdate int64, newUser users.User) (users.User, error) {
	filter := bson.M{"id": idToUpdate}
	update := bson.M{"$set": bson.M{
		"nickname": newUser.Nickname,
		"email":    newUser.Email,
	}}
	updateResult, err := r.UpdateOne(ctx, filter, update)
	if err != nil {
		return users.User{}, errs.UserRepositoryError
	} else if updateResult.ModifiedCount == 0 {
		return users.User{}, errs.UserNotExist
	} else {
		return users.User{
			ID:       idToUpdate,
			Nickname: newUser.Nickname,
			Email:    newUser.Email,
		}, nil
	}
}

func (r *UserRepository) DeleteUser(ctx context.Context, idToDelete int64) error {
	filter := bson.M{"id": idToDelete}
	deleteResult, err := r.DeleteOne(ctx, filter)
	if err != nil {
		return errs.UserRepositoryError
	} else if deleteResult.DeletedCount == 0 {
		return errs.UserNotExist
	} else {
		return nil
	}
}

func New(coll *mongo.Collection, currentID int64) app.UserRepository {
	return &UserRepository{
		Collection: coll,
		nextID:     currentID,
		mu:         new(sync.Mutex),
	}
}

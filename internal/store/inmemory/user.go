package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/mwdev22/CarRental/internal/types"
)

type UserRepo struct {
	users  map[int]types.User
	mu     sync.RWMutex
	nextID int
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users:  make(map[int]types.User),
		nextID: 1,
	}
}

func (r *UserRepo) Create(ctx context.Context, u *types.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u.ID = r.nextID
	r.nextID++

	r.users[u.ID] = *u
	return nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*types.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user with username %s not found", username)
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*types.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	return &user, nil
}

func (r *UserRepo) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user with id %d not found", id)
	}

	delete(r.users, id)
	return nil
}

func (r *UserRepo) Update(ctx context.Context, u *types.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[u.ID]; !exists {
		return fmt.Errorf("user with id %d not found", u.ID)
	}

	r.users[u.ID] = *u
	return nil
}

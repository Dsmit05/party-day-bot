package cache

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"sync"

	"github.com/Dsmit05/party-day-bot/internal/consts"
	"github.com/Dsmit05/party-day-bot/internal/models"
)

type repositoriesI interface {
	CreateUser(ctx context.Context, user models.User) error
	UpdateRole(ctx context.Context, id int64, role string) error
	ReadUsers(ctx context.Context) (map[int64]models.User, error)
}

type configI interface {
	IsDBUse() bool
}

type UserCache struct {
	users map[int64]models.User
	rw    *sync.RWMutex

	cfg configI
	rep repositoriesI
}

func NewUserCache(rep repositoriesI, cfg configI) (*UserCache, error) {
	cache := UserCache{users: make(map[int64]models.User), rw: &sync.RWMutex{}, rep: rep, cfg: cfg}

	if !cfg.IsDBUse() {
		return &cache, nil
	}

	dbUsers, err := rep.ReadUsers(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "rep.ReadUsers()")
	}

	if dbUsers == nil {
		return nil, errors.New("dbUsers == nil")
	}

	cache.users = dbUsers

	return &cache, nil
}

func (u *UserCache) GetAdmins(ctx context.Context) []int64 {
	u.rw.RLock()
	defer u.rw.RUnlock()

	userIDs := make([]int64, 0)

	for k, val := range u.users {
		if val.Role == consts.RoleAdmin {
			userIDs = append(userIDs, k)
		}
	}

	return userIDs
}

// GetUsers return map[userID]UserName
func (u *UserCache) GetUsers(ctx context.Context) map[int64]string {
	u.rw.RLock()
	defer u.rw.RUnlock()

	userIDs := make(map[int64]string)

	for k, val := range u.users {
		if val.Role != consts.RoleAdmin {
			userIDs[k] = val.UserName
		}
	}

	return userIDs
}

func (u *UserCache) CheckAccess(ctx context.Context, userID int64) bool {
	u.rw.RLock()
	defer u.rw.RUnlock()

	if user, ok := u.users[userID]; ok {
		if user.Role == consts.RoleAdmin {
			return true
		}
	}

	return false
}

// AddUser add user in cache and DB, not check is exist
func (u *UserCache) AddUser(ctx context.Context, user models.User) error {
	u.rw.Lock()
	defer u.rw.Unlock()

	u.users[user.TgID] = user

	if !u.cfg.IsDBUse() {
		return nil
	}

	if err := u.rep.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *UserCache) UpdateRole(ctx context.Context, userID int64, role string) error {
	u.rw.Lock()
	defer u.rw.Unlock()

	user, ok := u.users[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}

	user.Role = role
	u.users[userID] = user

	if !u.cfg.IsDBUse() {
		return nil
	}

	if err := u.rep.UpdateRole(ctx, userID, role); err != nil {
		return err
	}

	return nil
}

func (u *UserCache) CheckUser(_ context.Context, userID int64) bool {
	u.rw.RLock()
	_, ok := u.users[userID]
	u.rw.RUnlock()

	return ok
}

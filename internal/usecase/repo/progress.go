package repo

import (
	"strconv"

	"github.com/patrickmn/go-cache"

	"github.com/redstar01/sixtysec/internal/entity"
	"github.com/redstar01/sixtysec/internal/usecase"
)

type progressRepo struct {
	c *cache.Cache
}

func NewProgressRepo(c *cache.Cache) usecase.ProgressRepo {
	return &progressRepo{c: c}
}

func (pr *progressRepo) Get(p usecase.Player) (entity.GameProgress, error) {
	c, found := pr.c.Get(strconv.Itoa(int(p)))
	if !found {
		return entity.GameProgress{}, entity.ErrProgressNotFound
	}

	gp, ok := c.(entity.GameProgress)
	if !ok {
		return entity.GameProgress{}, entity.ErrProgressCachedTypeMismatch
	}

	return gp, nil
}

func (pr *progressRepo) Set(p usecase.Player, gp entity.GameProgress) error {
	pr.c.Set(strconv.Itoa(int(p)), gp, cache.DefaultExpiration)

	return nil
}

func (pr *progressRepo) Delete(p usecase.Player) error {
	pr.c.Delete(strconv.Itoa(int(p)))

	return nil
}

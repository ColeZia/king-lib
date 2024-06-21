package data

import (
	"context"

	"gl.king.im/king-lib/framework/test/skeleton/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

var _ biz.SkeletonRepo = (*skeletonRepo)(nil)

type skeletonRepo struct {
	data *Data
	log  *log.Helper
}

func NewSkeletonRepo(data *Data, logger log.Logger) biz.SkeletonRepo {
	r := &skeletonRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/skeletonRepo")),
	}
	return r
}

func (r *skeletonRepo) Get(ctx context.Context) (val string, err error) {

	return "", nil
}

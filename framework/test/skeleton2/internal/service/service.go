package service

import (
	admv1 "gl.king.im/king-lib/framework/test/skeleton2/internal/service/admin/v1"
	srvV1 "gl.king.im/king-lib/framework/test/skeleton2/internal/service/service/v1"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(srvV1.NewSkeletonService, srvV1.NewSkeletonGinService, admv1.NewSkeletonAdminService)

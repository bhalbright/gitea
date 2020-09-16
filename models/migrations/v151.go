// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package migrations

import (
	"xorm.io/xorm"
)

func addHookTaskPurge(x *xorm.Engine) error {
	type Repository struct {
		ID                              int64 `xorm:"pk autoincr"`
		OverridePruneHookTaskEnabled    bool
		OverrideWebhookDeliveriesToKeep int64
	}

	if err := x.Sync2(new(Repository)); err != nil {
		return err
	}

	return err
}

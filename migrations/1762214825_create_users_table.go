package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1762214825_create_users_table",
		Up: func() error {
			return olympian.Table("users").Create(func() {
				olympian.Uuid("id").Primary()
				olympian.String("name")
				olympian.String("email").Unique()
				olympian.String("password")
				olympian.Timestamp("email_verified_at").Nullable()
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("users").Drop()
		},
	})
}

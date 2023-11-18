package concurrency

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Concurrency struct {
	DB *gorm.DB
	mu sync.Mutex
}

func NewConcurrency(DB *gorm.DB) *Concurrency {
	return &Concurrency{
		DB: DB,
	}

}
func (un *Concurrency) Concurrency() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			un.mu.Lock()
			if err := un.DB.Exec(`
			UPDATE admins
			SET is_blocked = false
			WHERE id IN (
				SELECT admin_id
				FROM admin_infos
				WHERE admin_infos.admin_id = admins.id
				AND admin_infos.block_until < NOW()
			);
			
			DELETE FROM admin_infos
			WHERE admin_id IN (
				SELECT id
				FROM admins
				WHERE admins.id = admin_infos.admin_id
				AND admins.is_blocked = false
			);
            `).Error; err != nil {
				fmt.Println(err.Error())
			}
			if err := un.DB.Exec(`
			UPDATE users
            SET is_blocked = false
            WHERE id IN (
            SELECT users_id
            FROM user_infos
            WHERE user_infos.users_id = users.id
            AND user_infos.block_until < NOW()
             );

            DELETE FROM user_infos
			WHERE users_id IN (
    		SELECT id
    		FROM users
    		WHERE users.id = user_infos.users_id
    		AND users.is_blocked = false
			);

			`).Error; err != nil {
				fmt.Println(err.Error())
			}
			if err := un.DB.Exec(`
			 UPDATE orders SET order_status_id=2 WHERE order_date + INTERVAL '5 minutes'< NOW()
			 AND order_status_id=1
			 `).Error; err != nil {
				fmt.Println(err)
			}
			un.mu.Unlock()

			fmt.Println("worked")
		}
	}()
}

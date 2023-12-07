package concurrency

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
	"main.go/internal/web/middleware"
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
			if err := un.DB.Exec(`
			UPDATE coupons SET is_disabled=true WHERE coupons.created_at+INTERVAL '2 weeks' < NOW() 
			`).Error; err != nil {
				fmt.Println(err)
			}
			if err := un.DB.Exec(`
			DELETE FROM discounts WHERE expiry_date<NOW()
			`).Error; err != nil {
				fmt.Println(err)
			}
			// Check for users with 10 or more completed orders
			var usersIds []int
			err := un.DB.Raw(`
			SELECT users.id
			FROM users
			JOIN orders ON users.id = orders.user_id
			WHERE orders.order_status_id = 4
			GROUP BY users.id
			HAVING COUNT(orders.id) >= 10
			
			`).Scan(&usersIds).Error
			if err != nil {
				fmt.Println(err)
				un.mu.Unlock()
				continue
			}

			// Iterate through the result set
			for _, userId := range usersIds {

				// Check if the user has already received a coupon
				var couponCount int64
				err := un.DB.Raw(`SELECT COUNT(*) FROM user_reward_coupons WHERE users_id=$1 `, userId).Scan(&couponCount).Error
				if err != nil {
					fmt.Println(err)
				}
				if couponCount == 0 {
					var email string
					err := un.DB.Raw(`SELECT email FROM users WHERE id=?`, userId).Scan(&email).Error
					if err != nil {
						fmt.Println(err)
					}
					// Save the coupon in the database to track that the user has received it
					un.DB.Exec(`INSERT INTO user_reward_coupons (users_id,coupon_id) VALUES ($1,9)`, userId)
					// Generate and send the coupon code
					couponCode := "us-reward-10"
					middleware.SendCouponEmail(email, couponCode)

				}
			}
			un.mu.Unlock()

			fmt.Println("worked")
		}
	}()
}

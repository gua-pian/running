package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func CountAll(c *gin.Context) {
	all_users, _ := redisClient.SMembers("joined_people").Result()
	count := 0
	for _, user := range all_users {
		flag := 0
		count += 1
		slice_user_activities, _ := redisClient.HGetAll(user).Result()
		for k, _ := range slice_user_activities {
			if strings.Contains(k, "steps") {
				flag = 1
				count += 1
			}
		}
		if flag == 1 {
			count -= 1
		}
	}
	str_count := strconv.Itoa(count)
	c.String(http.StatusOK, str_count)
}

func Infos(c *gin.Context) {
	all_users, _ := redisClient.SMembers("joined_people").Result()
	// count := 0
	for _, user := range all_users {
		user_info, _ := redisClient.HGetAll(user).Result()
		if user_info["name"] != "" {
			fmt.Println(user_info["name"], "	", user_info["Sex"], "	", user_info["email"])
		}
	}
	c.String(http.StatusOK, "ok")
}

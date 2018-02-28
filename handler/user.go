package handler

import (
	"fmt"
	"strconv"
)

func NewUser(OpenId, HeadImageURL string, Sex int) {
	// Set steps and kilos into wxids.
	_, _ = redisClient.HSet(OpenId, "Sex", Sex).Result()
	_, _ = redisClient.HSet(OpenId, "HeadImageURL", HeadImageURL).Result()
}

func Count(OpenId string) (int, int, string, int, string) {
	total_users, _ := redisClient.SCard("joined_people").Result()
	total_steps, _ := redisClient.Get("total_steps").Result()
	total_kilos, _ := redisClient.Get("total_kilos").Float64()
	str_total_kilos := fmt.Sprintf("%.3f", total_kilos)

	int_total_steps, _ := strconv.Atoi(total_steps)
	// int_total_kilos, _ := strconv.Atoi(total_kilos)

	personal_total_steps, _ := redisClient.Get(OpenId + ":total_steps").Result()
	personal_total_kilos, _ := redisClient.Get(OpenId + ":total_kilos").Float64()
	str_personal_total_kilos := fmt.Sprintf("%.3f", personal_total_kilos)
	int_personal_total_steps, _ := strconv.Atoi(personal_total_steps)
	// int_personal_total_kilos, _ := strconv.Atoi(personal_total_kilos)

	return int(total_users), int_total_steps, str_total_kilos, int_personal_total_steps, str_personal_total_kilos
}

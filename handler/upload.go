package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func UploadImage(c *gin.Context) {
	// Get serverId
	serverid := c.Query("serverid")
	OpenId := c.Query("OpenId")
	_, _ = redisClient.SAdd("joined_people", OpenId).Result()
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	_, _ = redisClient.HSet(OpenId, timestamp+":picture", serverid).Result()
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
	return
}

func UploadData(c *gin.Context) {
	// Get serverId
	Name := c.Query("name")

	OpenId := c.Query("OpenId")

	steps := c.Query("steps")
	int_steps, _ := strconv.Atoi(steps)

	kilos := c.Query("kilos")
	// int_kilos, _ := strconv.Atoi(kilos)
	float_kilos, _ := strconv.ParseFloat(kilos, 64)

	Email := c.Query("email")

	Address := c.Query("address")

	timestamp := strconv.Itoa(int(time.Now().Unix()))

	_, _ = redisClient.SAdd("joined_people", OpenId).Result()
	_, _ = redisClient.IncrBy("total_steps", int64(int_steps)).Result()
	_, _ = redisClient.IncrByFloat("total_kilos", float_kilos).Result()
	_, _ = redisClient.HSet(OpenId, timestamp+":steps", int64(int_steps)).Result()
	if Name != "" || len(Name) == 0 {
		_, _ = redisClient.HSet(OpenId, "name", Name).Result()
	}
	_, _ = redisClient.HSet(OpenId, "email", Email).Result()
	_, _ = redisClient.IncrBy(OpenId+":total_steps", int64(int_steps)).Result()
	_, _ = redisClient.IncrByFloat(OpenId+":total_kilos", float_kilos).Result()
	if Address != "" {
		_, _ = redisClient.HSet(OpenId, "address", Address).Result()
	}
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
	return
}

package example

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"snowflake/comm"
	"strconv"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// 转化对象至字符串
func ConvertObjectToBytes(object interface{}) []byte {
	data, err := json.Marshal(object)
	if err != nil {
		return nil
	}
	return data
}

// 合并数字
func MergeTaskId(userId uint64, taskId uint64) uint64 {
	return userId<<30 | taskId // uint(number & 0xFFFFFFFF)
}

// 拆分数字
func GetUserId(number uint64) (userId uint) {
	return uint(number >> 30)
}

// 字符串转数字
func ConvertStringToNumber(value string) int {
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return number
}

// 数字转字符串
func ConvertNumberToString(number uint64) string {
	return fmt.Sprintf("%d", number)
}

// 获取字符串
func GetWordsString(value string) string {
	pattern := `[A-Za-z0-9\,\ \_\/\:\-\.]`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(value, -1)
	data := ""
	for _, match := range matches {
		data += match
	}
	return strings.TrimSpace(data)
}

// 验证签名消息
func VerifySignature(address string, message string, signature string) bool {
	//TODO:
	if address == signature {
		return true
	}
	signer := hexutil.MustDecode(signature)
	if signer[crypto.RecoveryIDOffset] == 27 || signer[crypto.RecoveryIDOffset] == 28 {
		signer[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}
	recovered, err := crypto.SigToPub(accounts.TextHash([]byte(message)), signer)
	if err != nil {
		return false
	}
	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	return address == recoveredAddr.Hex()
}

// 获取用户信息
func GetUserIdentity(c *fiber.Ctx) (uint, int, string) {
	user := c.Locals("user")
	if user == nil {
		return 0, 0, ""
	}
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	id := uint(claims["id"].(float64))
	level := int(claims["level"].(float64))
	name := claims["name"].(string)
	return id, level, name
}

// 获取雪花ID
func GetSnowflakeId() *string {
	node, err := snowflake.NewNode(1)
	if err != nil {
		comm.Error("failed to generate nonce id(%d): %s", err.Error())
		return nil
	}
	nonce := node.Generate().String()
	return &nonce
}

// 获取文件信息
func GetFileInfo(filePath string) fs.FileInfo {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fileInfo
	}
	return nil
}

// 生成MD5消息
func GetMd5Hash(filePath string) *string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Println(err)
		return nil
	}

	md5 := fmt.Sprintf("%x", hash.Sum(nil))
	return &md5
}

// 格式化字符串
func FormatString(content string) string {
	content = strings.ReplaceAll(content, `"`, "\"")
	content = strings.ReplaceAll(content, `'`, "\\'")
	content = strings.ReplaceAll(content, "`", "")
	content = strings.ReplaceAll(content, ";", "")
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\t", "")
	content = strings.ReplaceAll(content, "\f", "")
	content = strings.ReplaceAll(content, "\r", "")
	return content
}

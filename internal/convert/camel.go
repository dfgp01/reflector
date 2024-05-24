package convert

import (
	"regexp"
	"strings"
)

// camelToSnakeCase 将驼峰式字符串转换为小写下划线形式
func CamelToSnakeCase(camelCase string) string {
	// 先将字符串的首字母转为小写
	snakeCase := strings.ToLower(string(camelCase[0])) + camelCase[1:]

	// 使用正则表达式找到所有大写字母并替换为下划线加小写形式
	re := regexp.MustCompile("([a-z])([A-Z])")
	snakeCase = re.ReplaceAllString(snakeCase, "${1}_${2}")

	return strings.ToLower(snakeCase)
}

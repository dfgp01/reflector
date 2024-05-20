# reflector
golang reflect helper

#   预期序列化目标

转换目标|转换目标|支持情况|说明
:---:|:---:|:---:|:---:
string|number|支持|使用strconv
string|[]number|支持|用逗号隔开
|||
json|struct、proto|encoding/json
json|[]struct、[]proto|encoding/json
json|map[string]interface{}|encoding/json
json|[]map[string]interface{}|encoding/json
json|[]map[string][]interface{}|不支持|interface{}已含[]interface{}
json|map[number]interface{}|不支持|json没有非字符串的key
json|map[number][]interface{}|不支持|同上
|||
bytes|proto|google.golang.org/protobuf/proto|pb3规范
bytes|[]proto|支持|用自己的打包算法
bytes|map[string]proto|支持
bytes|map[number]proto|支持
bytes|map[string][]proto|支持
bytes|map[number][]proto|支持
|||
struct|map[string]interface{}|不支持|可使用github.com/mitchellh/mapstructure
struct|map[number]interface|不支持|struct没有非字符串的key
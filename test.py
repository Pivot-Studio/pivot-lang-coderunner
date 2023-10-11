import requests
import json

url = "http://localhost:8080/coderunner"
code = """use std::io;
pub fn main() i64 {
    io::print_s("Hello, world!");
    return 0;
}"""

data = {"code": code}

headers = {"Content-Type": "application/json"}  # 设置请求头为 JSON 类型

response = requests.post(url, data=json.dumps(data), headers=headers)  # 将数据编码为 JSON 格式
result = response.json()

print(result)

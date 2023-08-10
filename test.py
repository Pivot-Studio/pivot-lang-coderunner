import requests
import json

url = "http://localhost:8080/coderunner"
code = """use std::io;
pub fn main() i64 {
    let result = getFibonacci(10);
    println!(result);
    return 0;
}

pub fn getFibonacci(n: i64) i64 {
    let pre = 0;
    let nxt = 0;
    let result = 1;
    for let i = 0; i < n; i = i + 1 {
        result = result + pre;
        pre = nxt;
        nxt = result;
    }
    return result;
}"""

data = {"code": code}

headers = {"Content-Type": "application/json"}  # 设置请求头为 JSON 类型

response = requests.post(url, data=json.dumps(data), headers=headers)  # 将数据编码为 JSON 格式
result = response.json()

print("Result:", result["result"])

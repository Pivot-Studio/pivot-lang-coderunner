import aiohttp
import asyncio
import json

async def send_request(session, url, data):
    async with session.post(url, data=data) as response:
        result = await response.json()
        return result

async def main():
    url = "http://localhost:8080/coderunner"
    code = """use std::io;
pub fn main() i64 {
    io::print_s("Hello, world!\\n");
    return 0;
}"""
    data = {"code": code}
    headers = {"Content-Type": "application/json"}

    async with aiohttp.ClientSession(headers=headers) as session:
        tasks = [send_request(session, url, json.dumps(data)) for _ in range(11)]
        results = await asyncio.gather(*tasks)

        for result in results:
            print(result)

if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main())

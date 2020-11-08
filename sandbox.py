import requests
import time
import json

url = "http://localhost:8003"
endpoint = "vehicles/1234/engine"

body = {"action": "START"}

start_time = time.time()
print(start_time)

print(url)
print("%s/%s" % (url, endpoint))

r=requests.post("%s/%s" % (url, endpoint), data=json.dumps(body))
print(r.content)
print(r.status_code)
elapsed_time = time.time() - start_time
print(elapsed_time)
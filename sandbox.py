import requests
import time

url = "http://localhost:8003"
endpoint = "vehicles/1235"

start_time = time.time()
print(start_time)

print(url)
print("%s/%s" % (url, endpoint))

r=requests.get("%s/%s" % (url, endpoint))
print(r.content)
print(r.status_code)
elapsed_time = time.time() - start_time
print(elapsed_time)
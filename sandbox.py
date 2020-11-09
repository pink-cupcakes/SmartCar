import requests
import time
import json

url = "http://localhost:8003"
requestList = [
    {"endpoint": "vehicles/1234", "method": "GET"},
    {"endpoint": "vehicles/1234/engine", "method": "POST", "body": json.dumps({"action": "START"})},
    {"endpoint": "vehicles/1234/doors", "method": "GET"},
    {"endpoint": "vehicles/1234/fuel", "method": "GET"},
    {"endpoint": "vehicles/1235/battery", "method": "GET"},

    # Failure vehicle doesn't exist
    {"endpoint": "vehicles/1236/battery", "method": "GET"},

    # Failure invalid action
    {"endpoint": "vehicles/1234/engine", "method": "POST", "body": json.dumps({"action": "FOOBAR"})},

    # Null battery
    {"endpoint": "vehicles/1234/battery", "method": "GET"},
]

resDict = {}

for request in requestList:
    body = ""
    if(request.has_key("body")):
        body = request["body"]

    start_time = time.time()
    r = requests.request(request["method"], "%s/%s" % (url, request["endpoint"]), data=body)
    elapsed_time = time.time() - start_time
    if(resDict.has_key(request["endpoint"])):
        resDict[request["endpoint"]].append({
            "status_code": r.status_code,
            "content": r.content,
            "request_time": elapsed_time
        })
        continue

    resDict[request["endpoint"]] = [{
        "status_code": r.status_code,
        "content": r.content,
        "request_time": elapsed_time
    }]

print(json.dumps(resDict, sort_keys=True, indent=4))
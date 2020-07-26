import json

import requests


with open("test.tf") as f:
    resp = requests.post("http://localhost:9000/process", json={"code": f.read()})

    print(json.dumps(resp.json(), indent=4))
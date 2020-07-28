import json
import os
import subprocess
from typing import Dict, Any

def process_object(item):
    if isinstance(item, dict):
        for key, value in item.items():
            if key == "provisioners" and isinstance(value, list) and len(value) != 0:
                raise ValueError("Usage of provisioners is not supported")
            else:
                process_object(value)
    elif isinstance(item, list):
        for list_item in item:
            process_object(list_item)

with open("main.tf", "w", encoding="utf8") as f:
    f.write(os.environ["TERRAFORM_SOURCE"])

subprocess.run(
    ["terraform", "init", "-no-color"], stdout=subprocess.DEVNULL
).check_returncode()

subprocess.run(
    ["terraform", "plan", "-out", "plan.tfplan" "-no-color"], stdout=subprocess.DEVNULL
).check_returncode()

terraform_detail_process = subprocess.run(
    ["terraform", "show", "-json", "plan.tfplan"], stdout=subprocess.PIPE
)
terraform_detail_process.check_returncode()

detail_json = json.loads(terraform_detail_process.stdout.decode("utf8"))

process_object(detail_json)

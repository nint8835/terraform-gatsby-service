import json
import os
import subprocess
import sys


def process_object(item):
    if isinstance(item, dict):
        for key, value in item.items():
            if key == "provisioners" and isinstance(value, list) and len(value) != 0:
                sys.exit("Usage of provisioners is not supported")
            else:
                process_object(value)
    elif isinstance(item, list):
        for list_item in item:
            process_object(list_item)


with open("main.tf", "w", encoding="utf8") as f:
    f.write(os.environ["TERRAFORM_SOURCE"])

try:
    subprocess.run(
        ["terraform", "init", "-no-color"], stdout=subprocess.DEVNULL
    ).check_returncode()

    subprocess.run(
        ["terraform", "plan", "-no-color", "-out", "plan.tfplan"],
        stdout=subprocess.DEVNULL,
    ).check_returncode()

    terraform_detail_process = subprocess.run(
        ["terraform", "show", "-no-color", "-json", "plan.tfplan"],
        stdout=subprocess.PIPE,
    )
    terraform_detail_process.check_returncode()

    detail_json = json.loads(terraform_detail_process.stdout.decode("utf8"))

    process_object(detail_json)

    subprocess.run(
        ["terraform", "apply", "-no-color", "-auto-approve", "plan.tfplan"],
        stdout=subprocess.DEVNULL,
    ).check_returncode()

    subprocess.run(["terraform", "output", "contents", "-no-color"]).check_returncode()
except subprocess.CalledProcessError:
    sys.exit(1)

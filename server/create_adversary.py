import requests
import json
import sys

def payload(arg):
    return {
        "atomic_ordering":arg[2:],
        "name":arg[0],
        "description":arg[1],
        "objective":"495a9828-cab1-44dd-a0ca-66e58177d8cc", # default objective
        "plugin":"stockpile", # default stockpile
        "tags":[],
    }

def exec_api(payload):
    headers = {
        'KEY': 'ADMIN123',
        'accept': 'application/json',
        'Content-Type' : 'application/json; charset=utf-8'
    }

    res = requests.post('http://pdxf.malhyuk.info:8888/api/v2/adversaries', data=json.dumps(payload), headers=headers)
    # print(str(res.status_code))
    print(json.loads(res.text)["adversary_id"])



if __name__ == '__main__':
    exec_api(payload(sys.argv[1:]))

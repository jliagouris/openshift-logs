import json
import requests

def fetch_loki_logs(auth_token, x_scoppe_org_id, query, start, end, limit = 5000):
    url = "https://loki-frontend-opf-observatorium.apps.smaug.na.operate-first.cloud/loki/api/v1/query_range?"
    url += "query=" + query
    url += "&start=" + start + "&end=" + end + "&limit=" + limit

    payload={}
    headers = {
        'X-Scope-OrgID': x_scoppe_org_id,
        'Authorization': auth_token
    }

    response = requests.request("GET", url, headers=headers, data=payload)
    print("response:")
    print(response)
    data = response.json()['data']['result']
    return data
from bottle import route, run, request
import json
import copy

with open("sample.json") as f:
    doc = json.load(f)
counter = {"value": 1}


@route('/', 'POST')
def index():
    global counter
    global doc
    doc_copy = copy.deepcopy(doc.copy())
    nodes = doc_copy["data"]["project"]["runs"]["nodes"]
    nodes=nodes[(len(nodes)-counter["value"]):len(nodes)]
    print(len(nodes))
    doc_copy["data"]["project"]["runs"]["nodes"] = nodes
    doc_copy["data"]["project"]["runs"]["totalCount"] = 5000+counter["value"]
    counter = {"value": counter["value"]+1}
    print(counter)
    return doc_copy


run(host='0.0.0.0', port=8888)
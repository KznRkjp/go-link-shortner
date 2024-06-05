import requests
from threading import Thread
from queue import Queue
from requests.sessions import Session
import random
import string
import time
from pprint import pprint

THREAD_COUNT = 10
URL_STREAM_COUNT = 10
URL = "http://localhost:8080/"

def callback(i, result):
    try:
        for k in range(i, i + URL_STREAM_COUNT ):
            result[k] = post_url()
    except KeyboardInterrupt:
        return


def post_url():
    result = {}
    url = URL
    data = url_gen()
    x = requests.post(url, data=data)
    # print(x)
    # print(x.cookies, x.json, x.url)
    result['full_url'] = data
    result['cookie'] = x.cookies.get("JWT")
    result['short_url'] = x.text
    result['status'] = x.status_code
    # print(data)
    return result

def delete_url(list):
    print("Deleting")
    url = URL + "api/user/urls"
    for item in list:
        data = []
        data.append(item['short_url'].split("/")[-1])
        cookies = dict(JWT=str.lstrip(item['cookie']))
        # print(cookies)
        x = requests.delete(url, json=data)
        if x.status_code != 401:
            print("Error, expected 401, got", x.status_code)
        x = requests.delete(url, json=data, cookies=cookies)
        if x.status_code != 202:
            print("Error, expected 202, got", x.status_code)
        # print(x.status_code)
        


def url_gen():
    url = "https://"
    url += ''.join(random.choices(string.ascii_lowercase + string.digits, k=random.randrange(3,10,1)))
    url += "."
    url += ''.join(random.choices(string.ascii_lowercase + string.digits, k=random.randrange(1,3,1)))
    url += "/"
    url += ''.join(random.choices(string.ascii_lowercase + string.digits, k=random.randrange(1,10,1)))
    return url

def populate_db():
    threads = [None] * THREAD_COUNT
    results = [None] * THREAD_COUNT * URL_STREAM_COUNT
    for i in range(THREAD_COUNT):
        threads[i] = Thread(target=callback, args=(i*URL_STREAM_COUNT, results))
        threads[i].start()
        # t.start()

    for i in range(THREAD_COUNT):
        threads[i].join
    time.sleep(1+0.1 * THREAD_COUNT + 0.1 * URL_STREAM_COUNT)
    return (results)

def check_populate_results(list):
    ok = 0
    not_ok = 0
    for item in list:
        if item["status"] == 201:
            ok +=1
        else:
            not_ok += 1
    print("*"*58)
    print("Populate check:")
    print("OK - ", ok)
    print("Failed -", not_ok)

def check_redirect(list):
    ok = 0
    not_ok = 0
    for item in list:
        # pprint(item)
        x = requests.get(item['short_url'], allow_redirects=False)
        if x.status_code == 307 and x.headers['Location'] == item['full_url']:
            ok += 1
        else:
            not_ok += 1
            print(x.status_code, x.headers['Location'], item['full_url'])
    print("*"*58)
    print("Redirect check:")
    print("OK - ", ok)
    print("Failed -", not_ok)

def check_after_delete(list):
    ok = 0
    not_ok = 0
    for item in list:
        # pprint(item)
        x = requests.get(item['short_url'], allow_redirects=False)
        if x.status_code == 410:
            ok += 1
        else:
            not_ok += 1
            print(x.status_code)
    print("*"*58)
    print("Response after delete check:")
    print("OK - ", ok)
    print("Failed -", not_ok)

def main():
    results = populate_db()
    check_populate_results(results)
    check_redirect(results)
    delete_url(results)
    check_after_delete(results)
    # pprint(results)
    # threads = [None] * THREAD_COUNT
    # results = [None] * THREAD_COUNT * URL_STREAM_COUNT
    # for i in range(THREAD_COUNT):
    #     threads[i] = Thread(target=callback, args=(i*URL_STREAM_COUNT, results))
    #     threads[i].start()
    #     # t.start()

    # for i in range(THREAD_COUNT):
    #     threads[i].join
    # time.sleep(1)
    # pprint (results)
    # print (" ".join(results))
    print("*"*58)

if __name__ == "__main__":
    q = Queue(maxsize=0) 

    main()
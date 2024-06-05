import requests
from threading import Thread
from queue import Queue
from requests.sessions import Session
import random
import string
import time
from pprint import pprint

THREAD_COUNT = 10
URL_STREAM_COUNT = 20
URL = "http://localhost:8080/"

def callback(i, result):
    try:
        for k in range(i, i + URL_STREAM_COUNT ):
            result[k] = post_url()
    except KeyboardInterrupt:
        return

def callback_multy(i, result):
    try:
        for k in range(i, i + URL_STREAM_COUNT ):
            # print(k)
            result[k] = post_url_multy()
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
    # print(x.cookies.get("JWT"), "post url")
    return result

def post_url_multy():
    result = {}
    url = URL + "api/shorten/batch"
    data = []
    for i in range(random.randint(2,4)):
        url_data ={
            "correlation_id":cor_id_gen(),
            "original_url":url_gen()
        }
        data.append(url_data)
    # print(data)
    x = requests.post(url, json=data)
    # print(x)
    # print(x.cookies, x.json, x.url)
    result['urls_list'] = data
    result['cookie'] = x.cookies.get("JWT")
    result['url_list'] = x.json()
    result['status'] = x.status_code
    # print(x.cookies.get("JWT"), "post url multy")
    if x.status_code != 201:
        print("Error, expected 201, got", x.status_code)
    return result

def delete_url(list):
    print("Deleting")
    url = URL + "api/user/urls"
    for item in list:
        data = []
        data.append(item['short_url'].split("/")[-1])
        cookies = dict(JWT=str.lstrip(item['cookie']))
        # print(cookies)
        x = requests.delete(url, json=data, cookies=cookies)
        # if x.status_code != 401:
        #     print("Error, expected 401, got", x.status_code)
        # x = requests.delete(url, json=data, cookies=cookies)
        if x.status_code != 202:
            print("Error, expected 202, got", x.status_code)
        # print(x.cookies.items(), "post delete url")
        # print(x.status_code)

def delete_m_url(list):
    print("Deleting multiple")
    url = URL + "api/user/urls"
    for item in list:
        data = []
        for url_m in item['url_list']:
            data.append(url_m['short_url'].split("/")[-1])
        cookies = dict(JWT=str.lstrip(item['cookie']))
        # print(cookies)
        x = requests.delete(url, json=data, cookies=cookies)
        # x = requests.delete(url, json=data)
        # if x.status_code != 401:
        #     print("Error, expected 401, got", x.status_code)
        # x = requests.delete(url, json=data, cookies=cookies)
        if x.status_code != 202:
            print("Error, expected 202, got", x.status_code)


def url_gen():
    url = "https://"
    url += ''.join(random.choices(string.ascii_lowercase + string.digits, k=random.randrange(3,10,1)))
    url += "."
    url += ''.join(random.choices(string.ascii_lowercase + string.digits, k=random.randrange(1,3,1)))
    url += "/"
    url += ''.join(random.choices(string.ascii_lowercase + string.digits, k=random.randrange(1,10,1)))
    return url

def cor_id_gen():
    return ''.join(random.choices(string.ascii_lowercase + string.digits, k=6))+ "-" + ''.join(random.choices(string.ascii_lowercase + string.digits, k=6)) +"-"+''.join(random.choices(string.ascii_lowercase + string.digits, k=6))

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

def populate_db_multirecords():
    threads = [None] * THREAD_COUNT
    results = [None] * THREAD_COUNT * URL_STREAM_COUNT
    for i in range(THREAD_COUNT):
        threads[i] = Thread(target=callback_multy, args=(i*URL_STREAM_COUNT, results))
        threads[i].start()
        # t.start()
    for i in range(THREAD_COUNT):
        threads[i].join
    time.sleep(3+0.1 * THREAD_COUNT + 0.1 * URL_STREAM_COUNT)
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

def check_m_redirects(list):
    ok = 0
    not_ok = 0
    for item in list:
        # pprint(item)
        for url in item['url_list']:
            x = requests.get(url['short_url'], allow_redirects=False)
            if x.status_code == 307:
                ok += 1
            else:
                not_ok += 1
                print(x.status_code, x.headers['Location'])
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
    print(10*"*","Single record check",10*"*")
    results = populate_db()
    check_populate_results(results)
    check_redirect(results)
    delete_url(results)
    check_after_delete(results)
    print("*"*58)
    
    print(10*"*","Multiple record check",10*"*")
    time.sleep(2)
    m_results = populate_db_multirecords()
    time.sleep(2)
    # pprint(m_results)
    check_populate_results(m_results)
    check_m_redirects(m_results)
    delete_m_url(m_results)

    print("*"*58)

if __name__ == "__main__":
    # q = Queue(maxsize=0) 

    main()
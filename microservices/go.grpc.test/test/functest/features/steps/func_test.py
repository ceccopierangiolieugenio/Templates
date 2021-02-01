from behave import *
import requests
from pythonping import ping

@given('the backend is running')
def step_impl(context):
    assert ping('backend').success()

@given('the frontend is running')
def step_impl(context):
    assert ping('web_server').success()

@given('the nginx is running')
def step_impl(context):
    assert ping('nginx').success()


@when('I request {name}')
def step_impl(context, name):
    r =requests.get(f'http://nginx/{name}')
    print("status code",r.status_code)
    assert r.status_code is 200
    print("message",r.text)
    context.response = r.text

@then('the {message} is returned')
def step_impl(context, message):
    print("response",context.response)
    print("message",message)
    assert context.response == message
import requests
import os
from dotenv import load_dotenv

def exchange_code():
  
    load_dotenv()

    API_ENDPOINT = 'https://discord.com/api/v10'
    CLIENT_ID = os.getenv("DISCORD_CLIENT_ID")
    CLIENT_SECRET = os.getenv("DISCORD_CLIENT_SECRET")
    REDIRECT_URI = "http://127.0.0.1:8000"
    DISCORD_CODE = "ixCfi4e6LVXReOp3va5Y9iiCF6tCm0"

    print("CLIENT_ID: ", CLIENT_ID)
    print("CLIENT SECRET: ", CLIENT_SECRET)
    print("CODE: ", DISCORD_CODE)


    data = {
    'grant_type': 'authorization_code',
    'code': DISCORD_CODE,
    'redirect_uri': REDIRECT_URI
    }
    headers = {
    'Content-Type': 'application/x-www-form-urlencoded'
    }
    r = requests.post('%s/oauth2/token' % API_ENDPOINT, data=data, headers=headers, auth=(CLIENT_ID, CLIENT_SECRET))
    r.raise_for_status()
    return r.json()

print(exchange_code())
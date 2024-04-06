import redis
from flask import Flask, Response
from flask_cors import CORS
import os
import json
import re

import google.generativeai as genai

from constants import PROMPT

import textwrap

from IPython.display import display
from IPython.display import Markdown


def to_markdown(text):
    text = text.replace('•', '  *')
    return Markdown(textwrap.indent(text, '> ', predicate=lambda _: True))


GOOGLE_API_KEY = os.getenv('GOOGLE_API_KEY')

genai.configure(api_key=GOOGLE_API_KEY)

model = genai.GenerativeModel('gemini-1.0-pro-latest')


redis_pass = os.getenv('REDIS_PASSWORD')
host = os.getenv('REDIS_HOST')
port = int(os.getenv('REDIS_PORT'))


app = Flask(__name__)
CORS(app, resources={r"/stream/*": {"origins": "http://localhost:3000"}})
cache = redis.Redis(host=host, port=port)


def preprocess_json_string(json_string):
    json_string = re.sub(
        r'([{,]\s*)([A-Za-z_][A-Za-z0-9_]*)\s*:', r'\1"\2":', json_string)
    return json_string


def parse_json_from_string(string):
    json_match = re.search(r'```json(.+?)```', string, re.DOTALL)
    if json_match:
        json_string = json_match.group(1).strip()

        json_string = preprocess_json_string(json_string)

        try:
            parsed_json = json.loads(json_string)
            return parsed_json
        except json.JSONDecodeError as e:
            print(f"Error parsing JSON: {e}")
            return None
    else:
        print("No JSON part found in the string.")
        return None


@app.route('/stream/<user_id>')
def stream(user_id):
    channel_name = f"user:{user_id}"
    pubsub = cache.pubsub()
    pubsub.subscribe(channel_name)
    return Response(event_stream(pubsub), mimetype="text/event-stream")


def get_content(message: str, selectedLayer):
    prompt = PROMPT()
    messages = prompt.scratch_generation(message=message, selectedLayer=selectedLayer)
    response = model.generate_content(messages)

    try:
        text = response.parts[0].text
        return text
    except:
        return response


def event_stream(pubsub):
    try:
        for message in pubsub.listen():
            if message['type'] == 'message':
                payload = json.loads(message['data'])
                selectedLayer = json.loads(payload['selectedLayer'])
                chunk = get_content(payload['prompt'],selectedLayer=selectedLayer)
                resp = {
                    "content": parse_json_from_string(chunk),
                }
                yield f"data: {json.dumps(resp)}\n\n"

    except GeneratorExit:
        pass


if __name__ == '__main__':
    app.run(host="0.0.0.0", debug=True)

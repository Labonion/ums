import redis
from flask import Flask, Response
from flask_cors import CORS
import os
import json
import re

import google.generativeai as genai

import textwrap

from IPython.display import display
from IPython.display import Markdown


def to_markdown(text):
  text = text.replace('•', '  *')
  return Markdown(textwrap.indent(text, '> ', predicate=lambda _: True))

GOOGLE_API_KEY=os.getenv('GOOGLE_API_KEY')

genai.configure(api_key=GOOGLE_API_KEY)

model = genai.GenerativeModel('gemini-pro')


redis_pass = os.getenv('REDIS_PASSWORD')
host = os.getenv('REDIS_HOST')
port = int(os.getenv('REDIS_PORT'))

app = Flask(__name__)
CORS(app, resources={r"/stream/*": {"origins": "http://localhost:3000"}})  # Allow requests from localhost:3000
cache = redis.Redis(host='redis', port=6379)

def preprocess_json_string(json_string):
    # Ensure keys are enclosed in double quotes
    json_string = re.sub(r'([{,]\s*)([A-Za-z_][A-Za-z0-9_]*)\s*:', r'\1"\2":', json_string)
    return json_string

def parse_json_from_string(string):
    # Extract JSON part using regular expression
    json_match = re.search(r'```json(.+?)```', string, re.DOTALL)
    if json_match:
        json_string = json_match.group(1).strip()

        # Preprocess JSON string
        json_string = preprocess_json_string(json_string)

        # Load JSON
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

def get_content(message: str):
    OUTPUT = """
    ```json
        <result>
    ```
    """
    DOM = """
            ```js
                export interface DOM {
                    type: string; // html element
                    props: {
                        style?: CSSProperties;
                        children: Array<DOM | string>; // Array of DOM | strings, children is a part of the props object, add empty array if no children is present
                        onClick?: (event: React.MouseEvent<HTMLElement>) => void; // onClick function for the current node
                        onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void; // onChange function for the current node
                        className?: string; // Tailwind classnames
                        type?: string; // html prop type
                        placeholder?: string; // placeholder if any
                        href?: string; // href for links
                        src?: string;
                        alt?: string;
                        [x: string]: any;
                    };
                    id: string  // id for the node
                }
                ```
            """
    EXAMPLE = """
    ```json
          {
    type: "div",
    props: {
      className: "w-full rounded-lg shadow bg-gray-50",
      children: [
        {
          type: "img",
          props: {
            src: "image url from unsplash",
            alt: "",
            children: [],
          },
          id: "id",
        },
        {
          type: "div",
          props: {
            className: "p-4",
            children: [
              {
                type: "h2",
                props: {
                  className: "font-semibold text-lg",
                  children: ["Card Title"],
                },
                id: "id1",
              },
              {
                type: "p",
                props: {
                  className: "text-gray-500",
                  children: [
                    "Lorem ipsum dolor sit amet",
                  ],
                },
                id: "id2",
              },
            ],
          },
          id: "id3",
        },
      ],
    },
    id: "main",
  };
    ```
    """
    messages = [
        {'role':'user',
         'parts': [
             "You are an expert React and Tailwind developer",
             f"Your goal is to generate a DOM json for the query i provide. The DOM and PROPS interface are as follows {DOM}",
             "Follow the comments given in the interfaces for each key",
             f"Use the following example {EXAMPLE}",
            "Use Tailwind classes only",
            "Use valid dummy images from unsplash",
            f"Generate the result in the format {OUTPUT}"
            ],
        },
        {'role':'model', 'parts':['sure, please provide me your query']},
        {'role': 'user', 'parts': message}
        ]

    response = model.generate_content(messages)

    try:
        text = response.parts[0].text
        return text
    except:
        return "Could Not Generate, Try again"

def event_stream(pubsub):
    try:
        for message in pubsub.listen():
            if message['type'] == 'message':
                payload = json.loads(message['data'])
                chunk = get_content(payload['prompt'])
                resp = {
                        "content": parse_json_from_string(chunk)
                    }
                yield f"data: {json.dumps(resp)}\n\n"

    except GeneratorExit:
        pass

if __name__ == '__main__':
    app.run(debug=True)
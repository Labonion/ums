from dataclasses import dataclass


@dataclass
class PROMPT:
    _EXAMPLE: str = """
```json
{
    "type": "div",
    "props": {
        "className": "w-full rounded-lg shadow bg-gray-50",
        "children": [
            {
                "type": "img",
                "props": {
                    "src": "image url from unsplash",
                    "alt": "",
                    "children": []
                },
                "id": "id"
            },
            {
                "type": "div",
                "props": {
                    "className": "p-4",
                    "children": [
                        {
                            "type": "h2",
                            "props": {
                                "className": "font-semibold text-lg",
                                "children": ["Card Title"]
                            },
                            "id": "id1"
                        },
                        {
                            "type": "p",
                            "props": {
                                "className": "text-gray-500",
                                "children": ["Lorem ipsum dolor sit amet"]
                            },
                            "id": "id2"
                        }
                    ]
                },
                "id": "id3"
            }
        ]
    },
    "id": "main"
}
```
"""
    _DOM: str = """
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
    _OUTPUT: str = """
    ```json
        <result>
    ```
    """

    _OUTPUT_TASKS: str = """
    ```json
      {
        "instructions" : [instruction1, instruction2, ....]
      }
    ```
    """

    _EXAMPLE_TASK: str = """
      ## Instruction set for creating a Dashboard
      ```json
      {
        "instructions" : ["Create the main div that fits the screen",
        "Divide the main div into 4 pars, header, footer, sidebar and body",
        "Add a fixed header to the main div",
        "Add a fixed footer to the main div",
        "Add a side bar in the main div",
        "Utilise the remaining space for the content of the dashboard"
        ]
      }
      ```
    """

    def task_division(self, message):
        return [
            {'role': 'user',
             'parts': [
                 "You are an expert UI Designer with a lot of creativity and knowledge in react and tailwind",
                 f"Your goal is to provide an array of descriptive instruction to the developer in order to develope the query",
                 f"Here's an example instruction set {self._EXAMPLE_TASK}",
                 f"Theme Output must be in the format {self._OUTPUT_TASKS}"
             ],
             },
            {'role': 'model', 'parts': [
                'sure, please provide me your query']},
            {'role': 'user', 'parts': message}
        ]

    def scratch_generation(self, message, selectedLayer):
        return [
            {'role': 'user',
             'parts': [
                 "You are an expert React and Tailwind developer with a lot of creativity",
                 f"Your goal is to generate a DOM json for the query i provide with respect to the className {selectedLayer['className']}. The DOM and PROPS interface are as follows {self._DOM}",
                 f"children & type should be present in each node, if not with an empty array",
                 "Always wrap functions like onClick, onChange etc around with double quotes",
                 "Follow the comments given in the interfaces for each key",
                 "Use Tailwind classes only",
                 "Use valid dummy images from unsplash",
                 "Use icons from font awesome",
                 f"Use the following example {self._EXAMPLE}",
                 f"Generate the result in the format {self._OUTPUT}"
             ],
             },
            {'role': 'model', 'parts': ['sure, please provide me your query']},
            {'role': 'user', 'parts': message}
        ]

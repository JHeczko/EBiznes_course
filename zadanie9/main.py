from fastapi import FastAPI

from pydantic import BaseModel

from ollama import chat
from ollama import ChatResponse

app = FastAPI()


class InputText(BaseModel):
    question: str

@app.post("/ask")
async def ask_gpt(req: InputText):

    response: ChatResponse = chat(model='gemma3', messages=[
        {
            'role': 'user',
            'content': f'{req.question}',
        },
    ])

    response_text = response.message.content

    return {"response": response_text}

if __name__ == '__main__':
    response: ChatResponse = chat(model='gemma3', messages=[
        {
            'role': 'user',
            'content': 'Why is the sky blue?',
        },
    ])
    print(response['message']['content'])
    # or access fields directly from the response object
    print(response.message.content)

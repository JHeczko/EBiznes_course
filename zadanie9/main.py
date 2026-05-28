import os

import uvicorn
from fastapi import FastAPI

from pydantic import BaseModel

from ollama import chat
from ollama import ChatResponse, Client

app = FastAPI()


class InputText(BaseModel):
    question: str

@app.post("/ask")
async def ask_gpt(req: InputText):
    print(f"--> Odebrano pytanie: {req.question}")
    try:
        ollama_url = os.getenv("OLLAMA_URL", "http://host.docker.internal:11434")
        print(f"--> Łączenie z Ollama pod: {ollama_url}")

        client = Client(host=ollama_url)

        response = client.chat(model='gemma3', messages=[
            {
                'role': 'user',
                'content': req.question,
            },
        ])

        print(f"--> Surowa odpowiedź z Ollamy: {response}")

        # Bezpieczne wyciąganie tekstu (obsługuje słownik i obiekt)
        if isinstance(response, dict):
            response_text = response.get('message', {}).get('content', '')
        else:
            response_text = getattr(getattr(response, 'message', None), 'content', str(response))

        print(f"--> Wyciągnięty tekst: {response_text}")
        return {"response": response_text}

    except Exception as e:
        import traceback
        error_trace = traceback.format_exc()
        print(f"!!! BŁĄD W TRY: \n{error_trace}")
        return {"error": str(e), "traceback": error_trace}

if __name__ == '__main__':
    print(os.environ["OLLAMA_URL"])
    uvicorn.run(app, host='0.0.0.0', port=8001)

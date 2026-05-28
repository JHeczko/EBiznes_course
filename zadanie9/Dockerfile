FROM python:3.12-alpine

WORKDIR /server

RUN  apk update && \
     apk add bash && \
     pip install fastapi "uvicorn[standard]" ollama

RUN addgroup server-group && \
    adduser -S server-user -G server-group

COPY ./main.py .

RUN chmod +x ./main.py && \
    chown server-user:server-group main.py

USER server-user

EXPOSE 8001

CMD ["python", "main.py"]
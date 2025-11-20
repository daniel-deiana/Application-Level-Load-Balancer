FROM python:3.9 AS build
COPY ./backend_tests /app/backend_tests
WORKDIR /app/backend_tests
RUN pip install -r requirements.txt
ENTRYPOINT ["python3","test_backend_basic.py"]
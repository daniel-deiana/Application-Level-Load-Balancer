FROM python:3.9 AS build

COPY ./backend_tests /app/backend_tests
COPY ./rr /app/rr

RUN pip install -r /app/backend_tests/requirements.txt
WORKDIR /app/rr
ENTRYPOINT ["python3","test_rr_logic.py"]
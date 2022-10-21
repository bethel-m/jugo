FROM python 
WORKDIR /data_processor
ENV PIPENV_VENV_IN_PROJECT=1
RUN pip install pipenv
COPY Pipfile ./
COPY Pipfile.lock ./
RUN pipenv sync
COPY data_processor.py ./
CMD ["pipenv","run","python3","data_processor.py"]
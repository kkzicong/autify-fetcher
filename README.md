# Usage

### 1. (Optional) Update config

Set the number of concurrent workers in `config.yml` to match the number of CPU cores in your environment to maximize CPU utilization.

### 2. Build the application docker container

```
docker build -t autify-fetch .
```

### 2. Start the container

Optional `-metadata` flag to display metadata

```
docker run -v ${PWD}/html:/app/html -it --rm --name running-autify-fetch autify-fetch /app/fetcher -metadata <url 1> <url 2> ...
```

HTML files will be downloaded in the `html` directory.

Example:

```
docker run -v ${PWD}/html:/app/html -it --rm --name running-autify-fetch autify-fetch /app/fetcher -metadata https://www.google.com https://autify.com

ls html
```
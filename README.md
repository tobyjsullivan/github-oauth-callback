# Github OAuth Callback

A simple go service which directs a client through the Github
OAuth flow and produces the final access token.

## Running

### Docker

Build the Docker image:

```
docker build -t github-oauth-callback .
```

Run the docker image with your github client credentials

```
docker run -ti -p 3000:3000 -e CALLBACK_URL=http://localhost:3000/callback -e GITHUB_CLIENT_ID=[YOUR_CLIENT_ID] -e GITHUB_CLIENT_SECRET=[YOUR_CLIENT_SECRET] github-oauth-callback
```


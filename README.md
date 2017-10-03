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

## Using with your SPA

To use this service with your single-page app:

1. Configure a callback URI which is a path on your SPA (not this app is hosted)
    * Be sure to set CALLBACK_URL accurately for this service.
2. Have your SPA direct users to this services `/begin` route.
    * At this point, the user will be directed to Github's auth flow.
3. Users will be directed back to your SPA's route as configured in step 1.
4. Capture the `code` and `state` parameter from the query string
5. Make an AJAX call to this service's `/callback` route with the same `code` and `state` param
    * Example: `GET http://localhost:3000/callback?code=...&state=...`
    * Note: Presently, you need to pass the parameters in the query string.
    There is still an outstanding task to convert this endpoint to a proper POST form method for this use case.
6. The Service should respond with a valid `access_token` for the user.

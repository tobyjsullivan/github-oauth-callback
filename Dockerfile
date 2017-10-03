FROM golang
ADD . /go/src/github.com/tobyjsullivan/github-oauth-callback
RUN  go install github.com/tobyjsullivan/github-oauth-callback

EXPOSE 3000

CMD /go/bin/github-oauth-callback


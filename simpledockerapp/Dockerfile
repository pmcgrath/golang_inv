FROM golang:1.3.1-onbuild

# We override the default cmd which is CMD ["go-wrapper", "run"] which will run "app" with no args. The override allows us to pass app args
CMD ["app", "-workers", "5", "-sleep", "2"]

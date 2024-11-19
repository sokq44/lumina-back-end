# Documentation

## Overview

Subjectively speaking, big comments don't look so good inside the code on which people are constantly working. Hence the best idea would be to hold a copy of the code somewhere locally, comment it there and then generate a static html documentation. If you don't yet possess the documented project copy which's been alreay put in place, reach out to the person who's made the last documentation commit.

## How to generate?

Since godoc tool doesn't support generating documentation but rather enforces the approach of hosting it, we have go around it a bit.

1. Host your documentation

```bash
godoc -http=:8080 # or any port you want
```

2. Grab the files 'manually'

```bash
wget -m -k -q -erobots=off --no-host-directories --no-use-server-timestamps http://localhost:8080 # or any port you've specified
```

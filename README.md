# sam-loves-beans

I was supposed to go visit family over the Thanksgiving holiday but got sick and
had to cancel, so I decided to play around and learn a little Go.

I have a collegue who LOVES beans. So I built this site dedicated to him and
beans.

Cheers 🍻, Sam

Built on:

- [OpenAI API](https://beta.openai.com/docs/introduction) 🧠
- [Gin](https://github.com/gin-gonic/gin) 🍸
- [fly.io](https://fly.io/) 🚀

This is my first project in Go and as such it is super hacky, so I'm sure there
is a lot of room for improvement...

I landed on fly.io for running the app since it was SUPER easy to setup.

## Notes

Apparently, building this into a binary doesn't include the templates or static
assets 🤷🏻‍♂️... so those need to be manually copied to a persistent storage
location. Make sure to include:

- ./data
- ./static
- ./templates

### OpenAI 🧠

To enable the OpenAI APIs you will need an `API_KEY` from OpenAI and you'll need
to set the env var `ENABLE_OPENAI=true`

### fly.io

Create the app, `API_KEY` and volume:

```bash
$ fly launch
$ flyctl secrets set API_KEY=sk-
$ fly volumes create beans -r lax -s 1
```

Fly offers persistent volumes (so we can store the image and text responses) -
I'm sure I could hook these into some caching layer, but that's overkill. Then
just add the volume to the toml:

```toml
[mounts]
  source="beans"
  destination="/app/data"
```

Fly also offers custom domain support - and uses LE under the hood. Sweet!

```bash
$ fly ips list
$ fly certs create sam-loves-beans.com
```

## Deploying

Deploying is easy. I could hook this into a GitHub Action... but I'm not gonna.

```shell
$ fly auth login
$ fly deploy
```

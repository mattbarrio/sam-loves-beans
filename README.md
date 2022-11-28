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
assets... so those need to be manually copied to a persistent storage location.
Make sure to include:

- ./data
- ./static
- ./templates

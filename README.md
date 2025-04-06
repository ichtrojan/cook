# Cook

> Keep Building </br>
> <hr>Sir Cook alot (2025)

## What's included

- [MySQL](https://www.mysql.com)
- [Chi Router](https://go-chi.io)
- [Host Router](https://github.com/go-chi/hostrouter)
- [Redis](https://redis.io)
- [GORM](https://gorm.io)
- [GOTP](https://github.com/ichtrojan/gotp)
- [Go Validator](https://github.com/thedevsaddam/govalidator)
- [Laravel Migration](https://laravel.com/docs/12.x/migrations) (I know this is crazy)
- [Asynq](https://github.com/hibiken/asynq)
- [Asynqmon](https://github.com/hibiken/asynqmon)
- [AWS SES](https://aws.amazon.com/ses)

## Usage

Create `.env` from template 

```bash
cp .env.example .env
```

Populate `.env` with appropriate values and run the server

```bash
go run server.go
```

Cook would be available on port `6666`

## Why does this exist?

I build a lot of random projects—most of them never make it past localhost. Eventually, I felt the need to pull out the commonly used pieces and turn them into a more structured template.

This setup probably won’t work for everyone—or maybe even most people writing Go—but I’m putting it out there anyway. It’s mainly for my own sanity, to help me spin up new projects faster.

## Who should use this?

- Lazy people 
- Non-vibe coders

## Time for some devsplaining

I would use this section to explain some not so obvious items.

### Go Validator

You can find more details in the original documentation on the package’s GitHub page. That said, I’ve added a few custom validation rules of my own:

- `array`: checks that the JSON attribute is a valid array
- `base64`: checks that it’s a valid base64-encoded string
- `uuid_array`: checks that it’s an array of valid UUIDs
- `array_unique_uuid`: checks that it’s an array of unique UUIDs
- `valid_future_timestamp`: checks that it’s a valid timestamp in the future

### Why do we have a laravel project in the `misc` folder?

Very good question. In my (not-so-long) stint writing Go, I haven’t found a migration library I actually like. I’ve tried a bunch—most of them expect you to write raw SQL like a caveman. Thing is, I came to Go after years of being pampered by Laravel's migration system.

So, in the misc folder, you’ll find a super stripped-down Laravel project whose only job is to manage migrations. It’s totally optional—you can delete it and use whatever tool you're more comfortable with. It's not critical to running the template at all.

### Asynq/Asynqmon

Asynq is available by default when you run the server. However, the Asynqmon monitoring dashboard isn’t enabled out of the box. To turn it on, set the `ASYNQMON_SERVICE` key to true in your .env file. After that, restart the server, and Asynqmon will be available at port `6660`.

## Closing remarks

This README was written while listening to LOCKED IN by Maison2500 ([Apple Music](https://music.apple.com/gb/album/locked-in/1780407404?i=1780407718), [Spotify](https://open.spotify.com/track/16mBGLL2q29af8swK1pnRW?si=0f28c04348a040bd)) not like it matters. There are a bunch of things I don't have the time to explain but I'm sure you'd figure out as you attempt to use this template.

<hr>

Before Nothing, there was me </br>
Trojan 𓂀 


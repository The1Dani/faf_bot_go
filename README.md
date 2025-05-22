## FAF24X Bot

This is the repository where our beloved bot is living. This project is open source but never the less this project was made for the fun of our colleagues.

The project contain a bot written in GO witch was rewritten from another bot thats in python [the original bot](https://github.com/Strelkoveg/BetterUserOfTheDay). The bot messages were modified to suit our internal pricoale (jokes) so if you like to run the project yourself change to suit your own.

### The main differences from the original project :
- The bot is written in GO.
- The database has been changed to postgres.
- The jokes were modified.
- Some of functionalities were removed.
- The function run times were reduced from one day to custom amount of times per day witch can be changed in the messages package.

### Current Functionalities:
- /pidor - The bad person of the day
- /run - The good person of the day
- /stats - The good stats
- /pidorstats - The bad pidorstats
- /percentstats - The percent stats
- /carmicdices - Switch for Carmic Dices Feature
- /reg - Registers you to the bot
- /unreg - Unregisters you from the bot

### How To Build?

You can build the project using the bash scripts provided in the scripts using the makefile. The scripts are develepment use only. In order to use the scripts you will need to create a **.env** file in the /scripts directory. With your API-TOKEN inside. 

```bash
#/scripts/.env
TELEGRAM_APITOKEN="YOUR-TELEGRAM-TOKEN"
```

To build the GO binary run

```bash
make
```

To build and run the executable while creating a new docker container for postgres if you don't already have one named some-db run:

```bash
make run
```

For configuration see [set-env.sh](./scripts/set-env.sh) 

To stop the already started DB use 

```bash
make stop-db
```

If you want to host this bot see [docker-compose](docker-compose.yaml)
# To BE CONTINUED
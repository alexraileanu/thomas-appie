### thomas

runs on a schedule (cronjob format) and makes requests to the ah's api to fetch bonus info about products. the products are saved in the db and generally need 4 fields:
- `api_name`: must match the name of the product from the appie website (example: for [this product](https://www.ah.nl/producten/product/wi1487/perla-huisblends-aroma-snelfiltermaling), the `api_name` would be `Perla Huisblends Aroma snelfiltermaling`.
- `friendly_name`: this is just a shorter version of the appie name and is only used for the notifications
- `referer_url`: the product page from the appie website
- `appie_id`: this is the code extracted from the `referer_url`. for example in this url: `https://www.ah.nl/producten/product/wi1487/perla-huisblends-aroma-snelfiltermaling` the `appie_id` would be `1487`. all products follow a similar pattern and the id can be extracted this way in some form or the other.

the schedule is configured in the `config.toml` file, which has a structure like so:

```toml
[thomas]
cron = "30 10 * * MON" # every monday at 10:30, standard cron format

[appie]
client_name = "abc"
client_version = "def"
user_agent = "ghi"
bonus_day = 1 # 0=sunday, 1=monday, 2=tuesday, etc.
client_platform_type = "Web"
```

`client_name`, `client_version` and `user_agent` are taken from the requests in the browser and the values i found work for me are: 

```toml
client_name = "ah-products"
client_version = "6.609.32"
user_agent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:134.0) Gecko/20100101 Firefox/134.0"
client_platform_type = "Web"
```

the `client_name` and `client_version` sometimes change and the server stops returning data, so that needs updating every once in a while. 

i personally deploy thomas in a docker container with:

```yaml
services:
  app:
    container_name: thomas.appie
    image: arcscloud/thomas-appie
    restart: unless-stopped
    env_file:
      - .env
    ports:
      - "7002:7008"
    networks:
      - mariadb
    volumes:
      - ./config:/app/config
      - ./data:/app/data

networks:
  mariadb:
    name: mariadb
    external: true
```

with a mariadb container defined elsewhere like:

```yaml
services:
  db:
    container_name: "mariadb"
    image: "mariadb:lts"
    restart: unless-stopped
    ports:
      - "3306:3306"
    env_file:
      - .env
    volumes:
      - ./data:/var/lib/mysql
    networks:
      - mariadb

networks:
  mariadb:
    name: mariadb
```

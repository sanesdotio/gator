# gator

An RSS feed aggregator and reader for the command line, built with Go.

## Installation & Setup

You need Go and Postgres installed on your system.

```bash
go install github.com/sanesdotio/gator@latest
```

Create file in your machine's home directory named '.gatorconfig.json' with the following content:

```json
{
  "db_url": "your postgres connection string",
  "current_user_name": "leave empty, gator will set it automatically"
}
```

Create a database by running the following command in the Postgres shell:

```bash
CREATE DATABASE gator
```

Connect to the database by running:

```bash
\c gator
```

From here, you can connect to the database using your preferred Postgres client or use the command line interface provided by gator:

```bash
psql gator
```

Finally, to create the necessary tables by either of the following methods:

- Manually run the SQL commands in 'sql/schema' folder.
- Install the 'goose' tool and run the migrations until the latest version:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose postgres <connection string> up
```

## Usage

gator is a command line tool that allows you to manage your RSS feeds. You can register users, add, remove, and browse feeds, as well as read articles. Each user can also manage a follow list of their favorite feeds.

gator will automatically scrape the registered feeds, and store the articles in the database. You can then browse the articles and read them directly from the command line.

### Commands

Commands can be executed by running `gator <command> <args>`. Here are the available commands:

#### Register a new user with the specified username:

```
gator register <username>
```

#### Log in with the specified username:

```
gator login <username>
```

#### Print a list of all registered users:

```
gator users
```

#### Add a new feed to the database with the specified URL:

```
gator addfeed <feed_name> <feed_url>
```

#### Follow a feed:

```
gator follow <feed_url>
```

#### Unfollow a feed:

```
gator unfollow <feed_url>
```

#### List all followed feeds by the current user:

```
gator following
```

#### List all feeds in the database:

```
gator feeds
```

#### Browse articles from followed feeds with an optional limit:

```
gator browse <limit>
```

#### Scrape feeds for new articles:

```
gator agg
```

#### Reset all data in the database:

```
gator reset
```

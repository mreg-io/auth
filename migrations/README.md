## Migrations

We use [Flyway](https://documentation.red-gate.com/fd/welcome-to-flyway-184127914.html) for database schema migration.

### File structure

All sql migration scripts are located in `./sql` directory.
In flyway, there are several types of migration, and the most common type is [Versioned migration](https://documentation.red-gate.com/fd/migrations-184127470.html#Migrations-VersionedMigrations).
Its filename should be like `V[version]__[description].sql`.

There are some special types of sql scripts like [callbacks](https://documentation.red-gate.com/fd/callback-concept-184127466.html).
Its name starts with a callback event and then followed by a description.
One example is `afterClean__type.sql`.

#### Version format

In order to make the migration version consistent, please run the following command to generate the `[version]` segment of migration filename.

```shell
date "+%Y.%m.%d.%H.%M.%S"
```

If your current time is `2024-09-05 22:17:20`, then the output of the above command will be `2024.09.05.22.17.20`.
The migration filename will be `V2024.09.05.22.17.20__[description].sql`.

### Usage

Thanks to the `flyway.toml` configuration, we're able to use the same migration script across development and production environment.
However, there are a few notes we need to follow to ensure smooth migration process.

1. After a migration script is merged into `main`, never modify its filename and content.
2. Create a new migration script if you work on a new issue and plan to modify the database schema.
3. Test the new migration script to find any underlying errors.
4. Multiline transaction is not allowed in accordance with [CockroachDB declarative schema changer](https://www.cockroachlabs.com/docs/stable/online-schema-changes#declarative-schema-changer).

#### Run migration locally

Running schema migration with Flyway is easy.
All you need to do is setting appropriate environment variables and firing the command.
Flyway will start the migration process and execute all pending migration scripts.

Set the following environment variables:
```shell
export FLYWAY_SCHEMAS=<database-schema>
export FLYWAY_USER=<database-user>
export FLYWAY_PASSWORD=<database-password>
```

Then, change the working directory to `migrations` and execute the following command to start migration:
```shell
flyway migrate
```

#### Clean development schema

During development, you might need to reset the database schema and re-apply the migration.
In order to achieve this, set the environment variables just as above and run the following command:
```shell
flyway clean
```

After the command finish, you have reset the schema.
Then, you can apply migration from the ground up.

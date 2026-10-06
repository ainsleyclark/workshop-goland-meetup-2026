# Infra

The plumbing of the app: The config, the SQLite driver, sqlc's generated output and `dbtest`. The
things the application needs in order to boot, rather than the things it calls out to.

The SQL you write isn't here. Migrations and queries live at the top of the app, in `migrations/`
and `queries/` beside `sqlc.yaml`, since nothing imports a `.sql` file and sqlc, goose and people
all find them more easily there. What sqlc generates from them lands back in `db/sqlc`.

The line against `clients/` is ownership: we run this, somebody else runs that.
Plumbing and generated code live here, decisions don't.

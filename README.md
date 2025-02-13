# CTF API Emulator

This project serves a few purposes:

1. I wanted a project to practice and learn Go
2. I wanted a project to build a REST API server
3. I wanted a back end that I could test some potential CTF tools against

The service is intended to emulate the Attack / Defend CTF API from the 2025 Wild West Hackin' Fest at Mile High (or as best as I could recall it from notes and tools that I built that day). It leverages a simple SQLite database that you'll have to populate with your own data. I included a sample [populate_table.sql](./src/populate_table.sql) file to facilitate this. Yes, it includes some example API keys, and yes the admin one is `test`.

This tool is **not** intended to emulate all functionality of the CTF back end. Flags won't automatically be added to the database, there's no scoring, etc.

For more information about how the database is configured, see [database.md](./database.md). I'm not a database engineer, so please don't judge my work there (though I do welcome feedback).

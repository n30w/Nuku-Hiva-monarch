# And Then Some

Reddit only caches 1000 saved posts+comments. This is a problem if you keep saving posts and want to view your oldest ones.

The name for this project is named "And Then Some", because it has the ability to retrieve all your Reddit posts, and then some, because it can retrieve any updates as well.

Pushes or merges to master branch will initiate github actions.

## Docker

Docker build command:

```docker build $(pwd) --tag n30w/andthensome:latest --label n30w/andthensome:latest```

Docker run command:

```docker run --env-file=env_vars --publish 4000:4000 andthensome```

## Notes

- Use ```make bin``` to create a binary, or ```make run``` to run it.

### Migration from MySQL

As of March 29, 2024, I have migrated from Planetscale's MySQL database to NeonDB's Postgresql database. Planetscale had discontinued their free tier on April 8, 2024. I used this guide to migrate to Neon: [Migrate from MySQL to Neon](https://neon.tech/docs/import/migrate-mysql). The guide uses ```pgloader```.

There's also [this online tool](https://www.sqlines.com/online) that converts MySQL statements to Postgresql ones. Pretty sure one could use ChatGPT as well.

## TODO

- [x] Store reddit credentials in .env file
- [x] SQL Comparison function
- [x] Setup webserver
- [x] Use a docker container for API calls
- [x] Deletion function for SQL db
- [ ] Implement Go Routines and channels for faster performance (blazingly fast you could say)
- [ ] Switch to Taskfile.dev
- [ ] Benchmark Tests for receiving data

## Links

- [Multiarchitecture Build With GitHub Actions Docker](https://dev.to/cloudx/multi-arch-docker-images-the-easy-way-with-github-actions-4k54)
- [Docker Building Multi-platform Images](https://docs.docker.com/build/building/multi-platform/)
- [How to get more JSON Results](https://old.reddit.com/r/redditdev/comments/d7egb/how_to_get_more_json_results_i_get_only_30/)
- [SQL Go Create Table Insert Row](https://golangbot.com/mysql-create-table-insert-row/)
- [Go Docker Tutorial](https://tutorialedge.net/golang/go-docker-tutorial/)
- [SQL Getting Unique Names when IDS are different](https://dba.stackexchange.com/questions/19511/getting-unique-names-when-the-ids-are-different-distinct)
- [Learn Git Basics](https://nulab.com/learn/software-development/git-tutorial/git-basics/#undoing-changes)

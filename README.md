# Code review slack bot

Code review slack bot is Dockerized Golang App, that collects and check (by command right now) github pull request statuses.
Automatically delete PRs that have 2 or more approves.

Each chat has its own separated database and may contain different github credentials and admins. 

Currently, bot has file based database and doesn't need something better.

This bot uses `slack socket connection` for inbound commands and `slack api` for message delivery.

#Setup

1. Clone project
   ```shell
    git clone git@github.com:DmitriyDev/cr_checker_bot.git 
    ```

2. Create environment
    ```shell
   cp .env.dist .env
    ```
3. Set credentials in `.env`   
   - SLACK_APP_TOKEN - Slack token for app. Begins with `xapp-`
   - SLACK_BOT_TOKEN - Slack token for bot. Begins with `xoxb-`
   - GITHUB_TOKEN - Github access token with read permissions
   - SILENT_MODE - Enable/Disable silence mode (Bot will not highlight users in messages)
   - ADMIN_USERS - List of Slack users ids with admin privileges(Delimiter - coma (`,`))
   - GITHUB_COMMUNICATION_MODE - Type of communication (`ASYNC` or `SYNC`). Default: `SYNC`

4. Run the service
    ```shell
    docker-compose up -d
    ```
5. Slack bot running and waiting for requests
6. Add your app to slack channel, where you want to use it

# Usage
#### Every message should contain mention of your app 
Example : `@CodeReviewApp !command`

1. Add PR to checker
   
   _(Github tocken should provide access to all PRs you sent)_
   
   Send message with all PRs you want to track.

   Example:
   ```text
   @CodeReviewApp https://github.com/some/project/pull/1
   https://github.com/some/anotherproject/pull/2
   https://github.com/some/project/pull/3/commits/234a5a495a5cdbd52e95323e2b28f0b4e3252bdd
   ```
   
   In reply tread you will receive information about every PR that bot found and gather information about.


2. Get information about all PR in bot DB 

   Example:
   ```text 
   @CodeReviewApp !stats
   ```

   In reply tread you will receive information about all PRs from DB
   (Without requests to github)

3. Update information about all PR

   Example:
   ```text 
   @CodeReviewApp !update
   ```
   In reply tread you will receive information about all PRs from DB
   (Before that bot will update info about every PR in DB)

4. (Admin only) Invalidate cached credentials

   _This app store some credentials data on DB, but cached it in memory. If you update DB manually, you will need an ability to reread credentials from DB_

   Example:
   ```text 
   @CodeReviewApp !invalidateCache
   ```


# Available commands

#### All commands started with mention app `@YourApp`

`!help` - Show available commands

`!add2team <slack user> <github login>` - (Admin only) add user to team

`!rm_user <slack user>` - (Admin only) remove user from team

`!team` - Show all users from team

`!stats` - Show all tracked PRs

`!stats_me` - Show all tracked PRs by user, who asked

`!update` - Update PR statuses and show new stats

`!invalidateCache` - (Admin only)Invalidate cached configs to reread from DB

`<PR url> [...<PR url>]` - Add PR to track
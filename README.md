# party-day-bot - golang service for u best day

You can't invite a photographer to every event, so why not ask everyone present
at the event to send their photos? or congratulations only to you? But why clog up a personal account, and then how to sort
out a huge number of messages?

party-day-but will help to optimize this procedure. Any user can send photos and congratulations to the bot,
band in a convenient format, the culprits of the celebration will always be able to view them in the chat of this bot.
Everything is centralized.

All the data that users will send to the bot can optionally be saved in the database.

Before launching, set up a secret command so that you can become an administrator. The bottom
will send various user events (photos, congratulations) to each administrator.

Since the administrator is a strong being, he can send a message to all the guests with one command.

### List of main commands
| commands      | access | description                                                   |
|---------------|--------|---------------------------------------------------------------|
| /whoAmI       | all    | gives you ID                                                  |
| /help         | all    | list of guest teams                                           |
| /helpA        | Admin  | extended list of commands for admin                           |
| /list         | Admin  | guest list                                                    |
| /root id      | Admin  | grants rights to the user                                     |
| /sendAll text | Admin  | send a message to all guests                                  |
| /secretCMD    | all    | Special configurable command to get super rights (see config) |

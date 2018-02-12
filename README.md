# tracr-bot
The executable which creates, starts, stops, and destroys bots. 

## Usage

To create a bot one must do the following

### Create a strategy
A strategy file looks something like this. The top level object contains a `name`, `pair`, `exchange` and `strategies` 
property. Each time the bot is run it executes __one__ of the `strategies`. The strategy that's picked to execute 
corresponds with the current trading position of the bot e.g. if the bot is in a _closed_ position it will execute the 
associated _closed_ strategy.

    {
      "name": "aBotNamedBill",
      "pair": "USD-BTC",
      "exchange": "kraken",
      "props": {
        "position" : "closed",
        "successfulRuns" : 4,
        "unsuccessfulRuns" : 2,
        "lastTimeRun": 1518067482
      },
      "actionHandler": "actionHandlerFunction",
      "strategies": [
        {
          "position": "closed",
          "trees": [
            {
              "name": "christmas tree",
              "root": {
                "isRoot": true,
                "condition": "TrueFunction",
                "action": null,
                "children": [
                  {
                    "isRoot": false,
                    "condition": "TrueFunction",
                    "action": "ShortPositionAction",
                    "children": []
                  }
                ]
              }
            }
          ]
        }
      ]
    }

Each `strategy` can have many decision `trees`. Each `tree` has many `signals`. `Signals` make up the branches and leaves 
of the `tree`. When a `strategy` is executed each `tree` starts with its root `signal`. The `signal` executes the function 
whos name is defined as `condition`. If the function returns *true* than its child `signals` are executed so and so forth 
down the tree. If the function returns *false* than its child `signals` are not executed. If a leaf `signal` 
i.e. the bottom most `signal` is reached and its condition function returns true than its `action` is added to the queue.


### Define `condition` functions
`Condition` functions test if a condition is true or false. Each `condition` function is passed several  arguments
- a struct representing the `vars` object defined in the strategy document. These are read-only (see action consumers)
- a reference to the `tracr-cache` package to retrieve information from the cache
- a reference to the `tracr-store` package to retrieve information from the store


### Define `action` functions
`Action` functions return an `Action` struct. `Actions` have many parts

- Intent - what the action intends to do e.g. open a short position, close its position etc.

- Order type - in the case where a position is opened this defines the type of order e.g. market, limit, etc.

- Consumer - who is going to handle this action. The action may be handled internally by the bot e.g. to update a the 
bot's `vars` property. Internal actions are passed to the `actionHandler` function.
External actions are handled by the Executor module. This module sends these actions 
to an exchange and handles the response. The bot doesn't finish execution until all `actions` have executed successfully
and the errors have been handled.

- Data - this map contains misc information pertaining to the action. e.g. the volume of a order, the leverage used, the 
currency pair, exchange name, etc. For internal actions this could contain the information for updating the bot's `vars` 
property.

---

### Operating your bots

To create a bot

    tracr create <botStrategyFilePath>
    
After the bot is created it can be started

    tracr [options...] start <botName>
    
The `-s` option can be used to start the bot immediately. The `-i` option defines an interval (in minutes) that the bot
executes on. e.g. if the bot should execute on the ten min mark (3:00, 3:10, 3:20, etc) than `-i 10` should be specified.
A default interval of 5 min is used.

To stop a bot

    tracr stop <botName>
    
To delete a bot

    tracr destroy <botName>
    
To list created bots

    tracr list

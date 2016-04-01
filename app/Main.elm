module Main (..) where

import Task
import Effects exposing (Never)
import StartApp
import RouteHash
import Signal
import Pinlist.Actions exposing (..)
import Pinlist.Model exposing (Model)
import Pinlist.View exposing (view)
import Pinlist.Update exposing (..)
import Pinlist exposing (init)
import Html exposing (Html)


messages : Signal.Mailbox Action
messages =
  Signal.mailbox NoOp


app : StartApp.App Model
app =
  StartApp.start
    { init = init
    , update = update
    , view = view
    , inputs = [ messages.signal ]
    }


port routeTasks : Signal (Task.Task () ())
port routeTasks =
  RouteHash.start
    { prefix = RouteHash.defaultPrefix
    , address = messages.address
    , models = app.model
    , delta2update = delta2update
    , location2action = location2action
    }


port tasks : Signal (Task.Task Never ())
port tasks =
  app.tasks


main : Signal Html.Html
main =
  app.html

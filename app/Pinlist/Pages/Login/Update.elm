module Pinlist.Pages.Login.Update (..) where

import Config exposing (serverUrl)
import Pinlist.App.Action as App
import Pinlist.Pages.Login.Action exposing (..)
import Pinlist.User exposing (UserAndToken)
import Pinlist.Pages.Login.Model exposing (..)
import Pinlist.User exposing (userAndTokenDecoder)
import Pinlist.Utils exposing (justModel)
import Http.Extra exposing (..)
import Json.Encode
import Effects exposing (Effects)
import Task exposing (Task)
import Result exposing (..)


update : Action -> Model -> ( Model, Effects App.Action )
update action model =
  case action of
    UpdateLogin v ->
      justModel { model | login = v }

    UpdatePass v ->
      justModel { model | pass = v }

    Submit ->
      ( { model | status = Loading }, login model.login model.pass )

    Login result ->
      case result of
        Ok resp ->
          ( initialModel
          , Task.succeed (App.SetUser ( resp.data.user, resp.data.token ))
              |> Effects.task
          )

        Err _ ->
          justModel
            { model
              | status = Ready
              , error = Maybe.Just "Invalid credentials"
            }


login : String -> String -> Effects App.Action
login username password =
  post (serverUrl ++ "account/login")
    |> withHeader "Content-Type" "application/json"
    |> withJsonBody
        (Json.Encode.object
          [ ( "login", Json.Encode.string username )
          , ( "password", Json.Encode.string password )
          ]
        )
    |> send (jsonReader userAndTokenDecoder) stringReader
    |> Task.toResult
    |> Task.map (\a -> App.LoginAction (Login a))
    |> Effects.task

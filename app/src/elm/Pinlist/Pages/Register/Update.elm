module Pinlist.Pages.Register.Update (..) where

import Config exposing (serverUrl)
import Pinlist.App.Action as App
import Pinlist.Pages.Register.Action exposing (..)
import Pinlist.User exposing (UserAndToken)
import Pinlist.Pages.Register.Model exposing (..)
import Pinlist.User exposing (userAndTokenDecoder)
import Pinlist.Utils exposing (justModel)
import Http.Extra exposing (..)
import Json.Encode
import Effects exposing (Effects)
import Task exposing (Task)
import Result exposing (..)
import Regex
import String


update : Action -> Model -> ( Model, Effects App.Action )
update action model =
  case action of
    UpdateUsername v ->
      justModel { model | username = v }

    UpdatePass v ->
      justModel { model | pass = v }

    UpdateEmail v ->
      justModel { model | email = v }

    Submit ->
      if not (Regex.contains (Regex.regex "^([a-zA-Z0-9]{2,60})$") model.username) then
        justModel { model | status = Ready, error = Just InvalidUsername }
      else if not (Regex.contains (Regex.regex "^.+@.+\\..+$") model.email) then
        justModel { model | status = Ready, error = Just InvalidEmail }
      else if String.length model.pass < 8 then
        justModel { model | status = Ready, error = Just InvalidPassword }
      else
        ( { model | status = Loading }, signup model )

    Register result ->
      case result of
        Ok resp ->
          ( initialModel
          , Task.succeed (App.SetUser ( resp.data.user, resp.data.token ))
              |> Effects.task
          )

        Err _ ->
          justModel { model | status = Ready, error = Maybe.Just DataTaken }


signup : Model -> Effects App.Action
signup model =
  post (serverUrl ++ "account/create")
    |> withHeader "Content-Type" "application/json"
    |> withJsonBody
        (Json.Encode.object
          [ ( "username", Json.Encode.string model.username )
          , ( "email", Json.Encode.string model.email )
          , ( "password", Json.Encode.string model.pass )
          ]
        )
    |> send (jsonReader userAndTokenDecoder) stringReader
    |> Task.toResult
    |> Task.map (\r -> App.RegisterAction (Register r))
    |> Effects.task

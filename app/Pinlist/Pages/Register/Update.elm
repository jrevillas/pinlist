module Pinlist.Pages.Register.Update (..) where

import Config exposing (serverUrl)
import Pinlist.App.Action as App
import Pinlist.Pages.Register.Action exposing (..)
import Pinlist.User exposing (UserAndToken)
import Pinlist.Pages.Register.Model exposing (..)
import Pinlist.User exposing (userAndTokenDecoder)
import Http.Extra exposing (..)
import Json.Decode.Extra exposing (..)
import Json.Decode exposing (..)
import Json.Encode
import Effects exposing (Effects)
import Task exposing (Task)
import Result exposing (..)


justModel : Model -> ( Model, Effects App.Action )
justModel model =
  ( model, Effects.none )


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
      ( { model | status = Loading }, signup model )

    Register result ->
      case result of
        Ok resp ->
          ( initialModel, Effects.none )

        Err _ ->
          justModel { model | status = Ready, error = Maybe.Just "Invalid fields" }


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

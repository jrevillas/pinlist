module Pinlist.Account.Effects (..) where

import Config exposing (serverUrl)
import Pinlist.App.Action as App
import Pinlist.User exposing (User, Token, encodeUser, encodeToken, userAndTokenDecoder)
import Http.Extra exposing (..)
import Pinlist.App.Action exposing (Action)
import Effects exposing (Effects)
import LocalStorage
import Task


saveUserAndToken : User -> Token -> Effects Action
saveUserAndToken user token =
  let
    _ =
      LocalStorage.set "user" (encodeUser user)

    _ =
      LocalStorage.set "token" (encodeToken token)
  in
    Effects.none


checkAuth : Token -> Effects Action
checkAuth token =
  post (serverUrl ++ "account/auth")
    |> withHeader "Content-Type" "application/json"
    |> withHeader "Authorization" ("bearer " ++ token.hash)
    |> send (jsonReader userAndTokenDecoder) stringReader
    |> Task.toResult
    |> Task.map App.Auth
    |> Effects.task

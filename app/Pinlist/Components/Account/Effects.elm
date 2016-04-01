module Pinlist.Components.Account.Effects (..) where

import Pinlist.Entities exposing (User, Token, UserAndToken)
import Pinlist.Actions exposing (..)
import Pinlist.Components.Account.Actions exposing (..)
import Config exposing (serverUrl)
import Http.Extra as HttpExtra exposing (..)
import Json.Decode.Extra exposing (..)
import Json.Decode exposing (..)
import Json.Encode
import Effects exposing (Effects)
import Task exposing (Task)


userDecoder : Decoder User
userDecoder =
  succeed User
    |: ("id" := int)
    |: ("email" := string)
    |: ("username" := string)
    |: ("status" := int)


tokenDecoder : Decoder Token
tokenDecoder =
  succeed Token
    |: ("hash" := string)
    |: ("until" := date)


userAndTokenDecoder : Decoder UserAndToken
userAndTokenDecoder =
  succeed UserAndToken
    |: ("user" := userDecoder)
    |: ("token" := tokenDecoder)


authUser : String -> Effects Action
authUser token =
  HttpExtra.post (serverUrl ++ "account/auth")
    |> withHeader "Content-Type" "application/json"
    |> withHeader "Authorization" ("bearer " ++ token)
    |> send (jsonReader userAndTokenDecoder) stringReader
    |> Task.toResult
    |> Task.map (\r -> Account (AuthUser r))
    |> Effects.task


login : String -> String -> Effects Action
login username password =
  HttpExtra.post (serverUrl ++ "account/login")
    |> withHeader "Content-Type" "application/json"
    |> withJsonBody
        (Json.Encode.object
          [ ( "login", Json.Encode.string username )
          , ( "password", Json.Encode.string password )
          ]
        )
    |> send (jsonReader userAndTokenDecoder) stringReader
    |> Task.toResult
    |> Task.map (\r -> Account (LoginUser r))
    |> Effects.task


signup : String -> String -> String -> Effects Action
signup username email password =
  HttpExtra.post (serverUrl ++ "account/create")
    |> withHeader "Content-Type" "application/json"
    |> withJsonBody
        (Json.Encode.object
          [ ( "username", Json.Encode.string username )
          , ( "email", Json.Encode.string email )
          , ( "password", Json.Encode.string password )
          ]
        )
    |> send (jsonReader userAndTokenDecoder) stringReader
    |> Task.toResult
    |> Task.map (\r -> Account (RegisterUser r))
    |> Effects.task

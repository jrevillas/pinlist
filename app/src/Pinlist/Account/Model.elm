module Pinlist.Account.Model (Model, initialModel, fromUserAndToken, emptyModel) where

import Pinlist.User exposing (User, Token, UserAndToken, userDecoder, tokenDecoder)
import Json.Decode
import Maybe exposing (..)
import LocalStorage


type alias Model =
  { user : Maybe User
  , token : Maybe Token
  , validated : Bool
  }


decodeToken : String -> Maybe Token
decodeToken str =
  case Json.Decode.decodeString tokenDecoder str of
    Ok token ->
      Just token

    Err _ ->
      Nothing


decodeUser : String -> Maybe User
decodeUser str =
  case Json.Decode.decodeString userDecoder str of
    Ok user ->
      Just user

    Err _ ->
      Nothing


initialModel : Model
initialModel =
  let
    token =
      case LocalStorage.get "token" of
        Just t ->
          decodeToken t

        Nothing ->
          Nothing

    user =
      case LocalStorage.get "user" of
        Just u ->
          decodeUser u

        Nothing ->
          Nothing
  in
    Model user token False


emptyModel : Model
emptyModel =
  Model Nothing Nothing False


fromUserAndToken : ( User, Token ) -> Model
fromUserAndToken userAndToken =
  Model
    (Just (fst userAndToken))
    (Just (snd userAndToken))
    True

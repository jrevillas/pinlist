module Pinlist.Account.Model (Model, initialModel) where

import Pinlist.User exposing (User, Token, UserAndToken)
import Maybe exposing (..)


type alias Model =
  { user : Maybe User
  , token : Maybe Token
  }


initialModel : Model
initialModel =
  Model Nothing Nothing


fromUserAndToken : UserAndToken -> Model
fromUserAndToken userAndToken =
  Model
    (Just userAndToken.user)
    (Just userAndToken.token)

module Pinlist.Pages.Login.Model (..) where

import Maybe exposing (..)


type Status
  = Ready
  | Loading


type ErrorMessage
  = EmptyField
  | InvalidCredentials


type alias Model =
  { login : String
  , pass : String
  , error : Maybe ErrorMessage
  , status : Status
  }


initialModel : Model
initialModel =
  Model "" "" Nothing Ready

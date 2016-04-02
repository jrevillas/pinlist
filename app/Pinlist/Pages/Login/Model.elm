module Pinlist.Pages.Login.Model (..) where

import Maybe exposing (..)


type Status
  = Ready
  | Loading


type alias Model =
  { login : String
  , pass : String
  , error : Maybe String
  , status : Status
  }


initialModel : Model
initialModel =
  Model "" "" Nothing Ready

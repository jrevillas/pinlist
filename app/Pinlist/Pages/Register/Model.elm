module Pinlist.Pages.Register.Model (..) where

import Maybe exposing (..)


type Status
  = Ready
  | Loading


type alias Model =
  { username : String
  , email : String
  , pass : String
  , error : Maybe String
  , status : Status
  }


initialModel : Model
initialModel =
  Model "" "" "" Nothing Ready

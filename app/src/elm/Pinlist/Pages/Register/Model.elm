module Pinlist.Pages.Register.Model (..) where

import Maybe exposing (..)


type Status
  = Ready
  | Loading


type ErrorMessage
  = InvalidUsername
  | InvalidEmail
  | InvalidPassword
  | DataTaken


type alias Model =
  { username : String
  , email : String
  , pass : String
  , error : Maybe ErrorMessage
  , status : Status
  }


initialModel : Model
initialModel =
  Model "" "" "" Nothing Ready

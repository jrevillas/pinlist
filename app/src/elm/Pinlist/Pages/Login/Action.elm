module Pinlist.Pages.Login.Action (..) where

import Pinlist.User exposing (UserAndToken)
import Http.Extra


type Action
  = UpdateLogin String
  | UpdatePass String
  | Submit
  | Login (Result (Http.Extra.Error String) (Http.Extra.Response UserAndToken))

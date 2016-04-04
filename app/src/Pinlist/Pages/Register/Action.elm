module Pinlist.Pages.Register.Action (..) where

import Pinlist.User exposing (UserAndToken)
import Http.Extra


type Action
  = UpdateUsername String
  | UpdatePass String
  | UpdateEmail String
  | Submit
  | Register (Result (Http.Extra.Error String) (Http.Extra.Response UserAndToken))

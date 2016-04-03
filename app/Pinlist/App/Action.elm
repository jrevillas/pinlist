module Pinlist.App.Action (..) where

import Pinlist.App.Model exposing (Page)
import Pinlist.User exposing (User, Token, UserAndToken)
import Pinlist.Pages.Login.Action as Login
import Pinlist.Pages.Register.Action as Register
import Http.Extra
import Result exposing (Result)


type Action
  = LoginAction Login.Action
  | RegisterAction Register.Action
  | SetActive Page
  | SetUser ( User, Token )
  | Auth (Result (Http.Extra.Error String) (Http.Extra.Response UserAndToken))
  | NoOp

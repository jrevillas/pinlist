module Pinlist.App.Action (..) where

import Pinlist.App.Model exposing (Page)
import Pinlist.Pages.Login.Action as Login
import Pinlist.Pages.Register.Action as Register
import Effects exposing (Effects)


type Action
  = LoginAction Login.Action
  | RegisterAction Register.Action
  | SetActive Page
  | NoOp

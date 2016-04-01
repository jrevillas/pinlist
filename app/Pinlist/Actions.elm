module Pinlist.Actions (..) where

import Pinlist.Components.Account.Actions exposing (AccountAction)
import Pinlist.Routes exposing (Page)


type Action
  = NoOp
  | SetActive (Page)
  | Account AccountAction

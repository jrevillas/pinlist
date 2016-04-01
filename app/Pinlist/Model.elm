module Pinlist.Model (..) where

import Pinlist.Components.Account.Model exposing (..)


{- Page model -}


type PageModel
  = Login LoginModel
  | Register RegisterModel
  | Empty


type alias Model =
  { account : AccountModel
  , pageModel : PageModel
  }


initialModel : Model
initialModel =
  Model
    initialAccountModel
    Empty

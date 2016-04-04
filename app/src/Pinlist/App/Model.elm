module Pinlist.App.Model (..) where

import Pinlist.Account.Model as Account
import Pinlist.Pages.Login.Model as Login
import Pinlist.Pages.Register.Model as Register
import Maybe exposing (..)


type Page
  = Login
  | Register
  | Home
  | Loading


type alias Model =
  { account : Account.Model
  , login : Login.Model
  , register : Register.Model
  , activePage : Page
  , nextPage : Maybe Page
  }


initialModel : Model
initialModel =
  Model
    Account.initialModel
    Login.initialModel
    Register.initialModel
    Loading
    Nothing

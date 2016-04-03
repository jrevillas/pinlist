module Pinlist.App.View (..) where

import Pinlist.App.Model as App exposing (Model)
import Pinlist.Pages.Login.View as Login
import Pinlist.Pages.Register.View as Register
import Pinlist.App.Action exposing (Action)
import Html exposing (..)


view : Signal.Address Action -> Model -> Html.Html
view address model =
  case model.activePage of
    App.Login ->
      Login.view address model.login

    App.Register ->
      Register.view address model.register

    _ ->
      div [] []

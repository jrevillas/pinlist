module Pinlist.View (..) where

import Pinlist.Model exposing (..)
import Pinlist.Components.Account.View exposing (loginView, registerView)
import Pinlist.Actions exposing (Action)
import Html exposing (..)


view : Signal.Address Action -> Model -> Html.Html
view address model =
  case model.pageModel of
    Login model' ->
      loginView address model'

    Register model' ->
      registerView address model'

    {- TODO: Remove this -}
    _ ->
      div [] []

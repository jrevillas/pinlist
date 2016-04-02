module Pinlist.App.Update (..) where

import Pinlist.Utils exposing (justModel)
import Pinlist.App.Model as App exposing (Model)
import Pinlist.Pages.Login.Update as Login
import Pinlist.Pages.Register.Update as Register
import Pinlist.App.Action exposing (..)
import Effects exposing (Effects)


update : Action -> Model -> ( Model, Effects Action )
update action model =
  case action of
    LoginAction action' ->
      let
        ( model', effects' ) =
          Login.update action' model.login
      in
        ( { model | login = model' }, effects' )

    RegisterAction action' ->
      let
        ( model', effects' ) =
          Register.update action' model.register
      in
        ( { model | register = model' }, effects' )

    SetActive page ->
      let
        model' =
          { model | activePage = page }
      in
        case page of
          App.Login ->
            justModel model'

          App.Register ->
            justModel model'

          App.Loading ->
            justModel model'

    NoOp ->
      justModel model

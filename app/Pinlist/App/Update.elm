module Pinlist.App.Update (..) where

import Pinlist.Utils exposing (justModel)
import Pinlist.App.Model as App exposing (Model)
import Pinlist.Account.Model as Account
import Pinlist.Account.Effects exposing (saveUserAndToken, checkAuth)
import Pinlist.Pages.Login.Update as Login
import Pinlist.Pages.Register.Update as Register
import Pinlist.App.Action exposing (..)
import Effects exposing (Effects)
import Debug


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

    SetUser userToken ->
      ( { model | account = Account.fromUserAndToken userToken, activePage = App.Home }
      , saveUserAndToken
          (fst userToken)
          (snd userToken)
      )

    Auth result ->
      case result of
        Ok resp ->
          ( { model | activePage = App.Home, account = Account.fromUserAndToken ( resp.data.user, resp.data.token ) }
          , saveUserAndToken
              resp.data.user
              resp.data.token
          )

        Err err ->
          justModel { model | activePage = App.Login, account = Account.emptyModel }

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
            case model.account.token of
              Nothing ->
                justModel { model | activePage = App.Login }

              Just token ->
                ( model', checkAuth token )

          App.Home ->
            justModel model'

    NoOp ->
      justModel model

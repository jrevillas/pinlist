module Pinlist.Components.Account (..) where

import Pinlist.Actions exposing (..)
import Pinlist.Components.Account.Model exposing (..)
import Pinlist.Components.Account.Actions exposing (..)
import Pinlist.Components.Account.Effects exposing (login)
import Pinlist.Utils exposing (justModel)
import Effects exposing (Effects)


updateLogin : AccountAction -> LoginModel -> ( LoginModel, Effects Action )
updateLogin action model =
  case action of
    ChangeLoginForm field val ->
      case field of
        LoginUsernameField ->
          justModel { model | username = val }

        LoginPasswordField ->
          justModel { model | password = val }

    SubmitLogin ->
      ( { model | writable = False }, login model.username model.password )

    _ ->
      justModel model


updateRegister : AccountAction -> RegisterModel -> ( RegisterModel, Effects Action )
updateRegister action model =
  case action of
    _ ->
      justModel model

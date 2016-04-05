module Pinlist.Pages.Login.View (..) where

import Pinlist.Pages.Login.Model exposing (Model)
import Pinlist.App.Action as App
import Pinlist.Pages.Login.Action exposing (..)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Pinlist.Utils exposing (preventDefault)


view : Signal.Address App.Action -> Model -> Html.Html
view address model =
  let
    sendMessage =
      \a -> Signal.message address (App.LoginAction a)
  in
    div
      [ class "login" ]
      [ h2 [] [ text "Sign in" ]
      , Html.form
          [ preventDefault address "submit" ]
          [ div
              [ class "form__field" ]
              [ label [ for "username" ] [ text "Username or email" ]
              , input
                  [ name "username"
                  , value model.login
                  , placeholder "Username or email"
                  , type' "text"
                  , on
                      "input"
                      targetValue
                      (\v -> sendMessage (UpdateLogin v))
                  ]
                  []
              ]
          , div
              [ class "form__field" ]
              [ label [ for "password" ] [ text "Password" ]
              , input
                  [ name "password"
                  , value model.pass
                  , placeholder "Password"
                  , type' "password"
                  , on
                      "input"
                      targetValue
                      (\v -> sendMessage (UpdatePass v))
                  ]
                  []
              ]
          , button
              [ type' "submit"
              , onClick address (App.LoginAction Submit)
              ]
              [ text "Login" ]
          ]
      ]

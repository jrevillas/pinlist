module Pinlist.Pages.Register.View (..) where

import Pinlist.Pages.Register.Model exposing (Model)
import Pinlist.App.Action as App
import Pinlist.Pages.Register.Action exposing (..)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Pinlist.Utils exposing (preventDefault)


view : Signal.Address App.Action -> Model -> Html.Html
view address model =
  let
    sendMessage =
      \a -> Signal.message address (App.RegisterAction a)
  in
    div
      [ class "register" ]
      [ h2 [] [ text "Sign up" ]
      , Html.form
          [ preventDefault address "submit" ]
          [ div
              [ class "form__field" ]
              [ label [ for "username" ] [ text "Username" ]
              , input
                  [ name "username"
                  , value model.username
                  , placeholder "Username"
                  , type' "text"
                  , on
                      "input"
                      targetValue
                      (\v -> sendMessage (UpdateUsername v))
                  ]
                  []
              ]
          , div
              [ class "form__field" ]
              [ label [ for "email" ] [ text "Email address" ]
              , input
                  [ name "email"
                  , value model.email
                  , placeholder "Email address"
                  , type' "email"
                  , on
                      "input"
                      targetValue
                      (\v -> sendMessage (UpdateEmail v))
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
              , onClick address (App.RegisterAction Submit)
              ]
              [ text "Register" ]
          ]
      ]

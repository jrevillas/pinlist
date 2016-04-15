module Pinlist.Pages.Register.View (..) where

import Pinlist.Pages.Register.Model exposing (..)
import Pinlist.App.Model as AppModel
import Pinlist.App.Action as App
import Pinlist.Pages.Register.Action exposing (..)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Pinlist.Utils exposing (preventDefault, link)


errorView : Maybe ErrorMessage -> Html.Html
errorView error =
  case error of
    Just err ->
      span
        [ class "form__error" ]
        [ text
            (case err of
              InvalidUsername ->
                "Invalid username. Must be alphanumeric between 2 and 60 characters."

              InvalidEmail ->
                "Invalid email."

              InvalidPassword ->
                "The password is not valid. A minimum length of 8 characters is required."

              DataTaken ->
                "Username or email already in use."
            )
        ]

    Nothing ->
      span [] []


view : Signal.Address App.Action -> Model -> Html.Html
view address model =
  let
    sendMessage =
      \a -> Signal.message address (App.RegisterAction a)
  in
    div
      [ class "form form--register" ]
      [ h2 [] [ text "Sign up" ]
      , errorView model.error
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
                  , disabled (model.status == Loading)
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
                  , disabled (model.status == Loading)
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
                  , disabled (model.status == Loading)
                  , on
                      "input"
                      targetValue
                      (\v -> sendMessage (UpdatePass v))
                  ]
                  []
              ]
          , button
              [ type' "submit"
              , disabled (model.status == Loading)
              , onClick address (App.RegisterAction Submit)
              ]
              [ text "Register" ]
          ]
      , div
          [ class "form--link" ]
          [ link AppModel.Login "I already have an account" address ]
      ]

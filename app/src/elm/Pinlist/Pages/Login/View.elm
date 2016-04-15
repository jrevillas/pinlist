module Pinlist.Pages.Login.View (..) where

import Pinlist.Pages.Login.Model exposing (..)
import Pinlist.App.Action as App
import Pinlist.App.Model as AppModel
import Pinlist.Pages.Login.Action exposing (..)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Pinlist.Utils exposing (preventDefault, link)
import Maybe exposing (..)


errorView : Maybe ErrorMessage -> Html.Html
errorView error =
  case error of
    Just err ->
      span
        [ class "form__error" ]
        [ text
            (case err of
              EmptyField ->
                "You must fill all the fields before submit."

              InvalidCredentials ->
                "Invalid username or password."
            )
        ]

    Nothing ->
      span [] []


view : Signal.Address App.Action -> Model -> Html.Html
view address model =
  let
    sendMessage =
      \a -> Signal.message address (App.LoginAction a)
  in
    div
      [ class "form form--login" ]
      [ errorView model.error
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
                  , disabled (model.status == Loading)
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
              , onClick address (App.LoginAction Submit)
              ]
              [ text "Login" ]
          ]
      , div
          [ class "form--link" ]
          [ link AppModel.Register "I don't have an account yet" address ]
      ]

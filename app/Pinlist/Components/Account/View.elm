module Pinlist.Components.Account.View (loginView, registerView) where

import Pinlist.Actions exposing (..)
import Pinlist.Utils exposing (preventDefault)
import Pinlist.Components.Account.Model exposing (..)
import Pinlist.Components.Account.Actions exposing (..)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)


loginView : Signal.Address Action -> LoginModel -> Html.Html
loginView address model =
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
                , value model.username
                , placeholder "Username or email"
                , type' "text"
                , on
                    "input"
                    targetValue
                    (\v ->
                      Signal.message address (Account (ChangeLoginForm LoginUsernameField v))
                    )
                ]
                []
            ]
        , div
            [ class "form__field" ]
            [ label [ for "password" ] [ text "Password" ]
            , input
                [ name "password"
                , value model.password
                , placeholder "Password"
                , type' "password"
                , on
                    "input"
                    targetValue
                    (\v ->
                      Signal.message address (Account (ChangeLoginForm LoginPasswordField v))
                    )
                ]
                []
            ]
        , button
            [ type' "submit"
            , onClick address (Account SubmitLogin)
            ]
            [ text "Login" ]
        ]
    ]


registerView : Signal.Address Action -> RegisterModel -> Html.Html
registerView address model =
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
                    (\v ->
                      Signal.message address (Account (ChangeRegisterForm RegisterUsernameField v))
                    )
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
                    (\v ->
                      Signal.message address (Account (ChangeRegisterForm RegisterEmailField v))
                    )
                ]
                []
            ]
        , div
            [ class "form__field" ]
            [ label [ for "password" ] [ text "Password" ]
            , input
                [ name "password"
                , value model.password
                , placeholder "Password"
                , type' "password"
                , on
                    "input"
                    targetValue
                    (\v ->
                      Signal.message address (Account (ChangeRegisterForm RegisterPasswordField v))
                    )
                ]
                []
            ]
        , button
            [ type' "submit"
            , onClick address (Account SubmitRegister)
            ]
            [ text "Register" ]
        ]
    ]

module Pinlist.Update (update, delta2update, location2action) where

import Pinlist.Actions exposing (..)
import Pinlist.Model exposing (..)
import Pinlist.Routes exposing (..)
import Pinlist.Components.Account exposing (..)
import Pinlist.Components.Account.Model exposing (..)
import Pinlist.Utils exposing (justModel)
import Effects exposing (Effects)
import RouteHash
import Maybe exposing (..)


type alias Updater a b =
  a -> b -> ( b, Effects Action )


type alias Converter a =
  a -> PageModel



{-
delegate passes the responsibility of transforming the
current model to another page component
-}


delegate : Model -> a -> b -> Updater b a -> Converter a -> ( Model, Effects Action )
delegate model model' action' fn type' =
  let
    r =
      fn action' model'
  in
    ( { model | pageModel = type' (fst r) }, snd r )


update : Action -> Model -> ( Model, Effects Action )
update action model =
  case action of
    SetActive page ->
      case page of
        LoginPage ->
          justModel { model | pageModel = Login initialLoginModel }

        RegisterPage ->
          justModel { model | pageModel = Register initialRegisterModel }

        {- TODO: remove this -}
        _ ->
          justModel model

    Account action' ->
      case model.pageModel of
        Login m ->
          delegate model m action' updateLogin Login

        Register m ->
          delegate model m action' updateRegister Register

        _ ->
          justModel { model | pageModel = Login initialLoginModel }

    NoOp ->
      justModel model


delta2update : Model -> Model -> Maybe RouteHash.HashUpdate
delta2update prev next =
  case next.account.user of
    Just user ->
      case next.pageModel of
        Login _ ->
          Just <| RouteHash.set [ "login" ]

        Register _ ->
          Just <| RouteHash.set [ "register" ]

        Empty ->
          Just <| RouteHash.set [ "login" ]

    Nothing ->
      {- If the user is not logged in we want to redirect
      always to login unless is the register page
      or the list page (TODO)
      -}
      case next.pageModel of
        Register _ ->
          Just <| RouteHash.set [ "register" ]

        _ ->
          Just <| RouteHash.set [ "login" ]


location2action : List String -> List Action
location2action args =
  case args of
    first :: rest ->
      case first of
        "login" ->
          [ SetActive LoginPage ]

        "register" ->
          [ SetActive RegisterPage ]

        "" ->
          [ SetActive LoginPage ]

        _ ->
          [ SetActive DashboardPage ]

    _ ->
      [ SetActive LoginPage ]

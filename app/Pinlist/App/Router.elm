module Pinlist.App.Router (..) where

import Pinlist.App.Model as App exposing (Model)
import Pinlist.App.Action exposing (..)
import Pinlist.Pages.Login.Model as Login
import Pinlist.Pages.Register.Model as Register
import RouteHash
import Maybe exposing (..)


delta2update : Model -> Model -> Maybe RouteHash.HashUpdate
delta2update prev next =
  case next.account.user of
    Just user ->
      case next.activePage of
        App.Login ->
          Just <| RouteHash.set [ "login" ]

        App.Register ->
          Just <| RouteHash.set [ "register" ]

        App.Empty ->
          Just <| RouteHash.set [ "login" ]

    Nothing ->
      {- If the user is not logged in we want to redirect
      always to login unless is the register page
      or the list page (TODO)
      -}
      case next.activePage of
        App.Register ->
          Just <| RouteHash.set [ "register" ]

        _ ->
          Just <| RouteHash.set [ "login" ]


location2action : List String -> List Action
location2action args =
  case args of
    first :: rest ->
      case first of
        "login" ->
          [ SetActive App.Login ]

        "register" ->
          [ SetActive App.Login ]

        _ ->
          [ SetActive App.Login ]

    _ ->
      [ SetActive App.Login ]

module Pinlist.App.Router (..) where

import Pinlist.App.Model as App exposing (Model)
import Pinlist.App.Action exposing (..)
import RouteHash
import Maybe exposing (..)


delta2update : Model -> Model -> Maybe RouteHash.HashUpdate
delta2update prev next =
  case next.account.user of
    Just user ->
      case next.activePage of
        App.Login ->
          Just <| RouteHash.set [ "" ]

        App.Register ->
          Just <| RouteHash.set [ "" ]

        App.Loading ->
          Just <| RouteHash.set [ "" ]

        App.Home ->
          Just <| RouteHash.set [ "home" ]

    Nothing ->
      case next.activePage of
        App.Register ->
          Just <| RouteHash.set [ "register" ]

        _ ->
          Just <| RouteHash.set [ "login" ]


location2action : List String -> List Action
location2action args =
  case args of
    [ "login" ] ->
      [ SetActive App.Login ]

    [ "register" ] ->
      [ SetActive App.Register ]

    "home" :: rest ->
      [ SetActive App.Home ]

    [ "" ] ->
      [ SetActive App.Loading ]

    _ ->
      -- TODO: show 404
      [ SetActive App.Loading ]

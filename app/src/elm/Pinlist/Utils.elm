module Pinlist.Utils (justModel, preventDefault, link) where

{- TODO: Split utils in types -}

import Pinlist.App.Action exposing (Action)
import Pinlist.App.Model as AppModel
import Effects exposing (Effects)
import Html.Events exposing (onWithOptions, targetValue, Options)
import Json.Decode as Json
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)


justModel : a -> ( a, Effects Action )
justModel model =
  ( model, Effects.none )


evtWithOpts : Signal.Address Action -> String -> Options -> Html.Attribute
evtWithOpts address evt opts =
  onWithOptions evt opts Json.value (\_ -> Signal.message address Pinlist.App.Action.NoOp)


preventDefault : Signal.Address Action -> String -> Html.Attribute
preventDefault address evt =
  evtWithOpts address evt { preventDefault = True, stopPropagation = False }


link : AppModel.Page -> String -> Signal.Address Action -> Html.Html
link page msg address =
  a
    [ href "javascript:;"
    , preventDefault address "click"
    , onClick address (Pinlist.App.Action.SetActive page)
    ]
    [ text msg ]

module Pinlist.Utils (justModel, preventDefault) where

{- TODO: Split utils in types -}

import Pinlist.App.Action exposing (Action)
import Effects exposing (Effects)
import Html.Events exposing (onWithOptions, targetValue, Options)
import Json.Decode as Json
import Html


justModel : a -> ( a, Effects Action )
justModel model =
  ( model, Effects.none )


evtWithOpts : Signal.Address Action -> String -> Options -> Html.Attribute
evtWithOpts address evt opts =
  onWithOptions evt opts Json.value (\_ -> Signal.message address Pinlist.App.Action.NoOp)


preventDefault : Signal.Address Action -> String -> Html.Attribute
preventDefault address evt =
  evtWithOpts address evt (Options False True)

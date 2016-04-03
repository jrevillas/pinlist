module Pinlist (..) where

import Effects exposing (Effects)
import Pinlist.App.Model exposing (initialModel, Model)
import Pinlist.App.Action exposing (Action)


init : ( Model, Effects Action )
init =
  ( initialModel
  , Effects.none
  )

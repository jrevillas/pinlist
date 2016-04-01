module Pinlist (..) where

import Effects exposing (Effects)
import Pinlist.Model exposing (initialModel, Model)
import Pinlist.Actions exposing (Action)


init : ( Model, Effects Action )
init =
  ( initialModel
  , Effects.none
  )

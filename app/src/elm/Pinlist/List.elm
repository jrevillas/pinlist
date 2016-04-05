module Pinlist.List (..) where

import Pinlist.User exposing (User)


type alias ListUser =
  { role : Int
  , user : User
  }


type alias PinList =
  { id : Int
  , name : String
  , description : String
  , public : Bool
  , pins : Int
  , users : List ListUser
  }

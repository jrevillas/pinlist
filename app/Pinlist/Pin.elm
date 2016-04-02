module Pinist.Pin (..) where

import Pinlist.User exposing (User)
import Pinlist.Tag exposing (Tag)
import Date exposing (Date)


type alias Pin =
  { id : Int
  , url : String
  , title : String
  , creator : User
  , createdAt : Date
  , tags : List Tag
  }

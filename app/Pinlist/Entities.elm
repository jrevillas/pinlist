module Pinlist.Entities (User, Token, UserAndToken, Tag, Pin, PinList) where

import Date exposing (..)


{- Entity Models -}


type alias User =
  { id : Int
  , email : String
  , username : String
  , status : Int
  }


type alias Token =
  { hash : String
  , until : Date
  }


type alias UserAndToken =
  { user : User
  , token : Token
  }


type alias Tag =
  { name : String
  , count : Int
  , id : Int
  }


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


type alias Pin =
  { id : Int
  , url : String
  , title : String
  , creator : User
  , createdAt : Date
  , tags : List Tag
  }

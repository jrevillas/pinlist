module Pinlist.User (..) where

import Json.Decode.Extra exposing ((|:), date)
import Json.Decode exposing ((:=), string, int, succeed, Decoder)
import Date exposing (Date)


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


userDecoder : Decoder User
userDecoder =
  succeed User
    |: ("id" := int)
    |: ("email" := string)
    |: ("username" := string)
    |: ("status" := int)


tokenDecoder : Decoder Token
tokenDecoder =
  succeed Token
    |: ("hash" := string)
    |: ("until" := date)


userAndTokenDecoder : Decoder UserAndToken
userAndTokenDecoder =
  succeed UserAndToken
    |: ("user" := userDecoder)
    |: ("token" := tokenDecoder)

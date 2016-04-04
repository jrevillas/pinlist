module Pinlist.User (..) where

import Json.Decode.Extra exposing ((|:), date)
import Json.Decode exposing ((:=), string, int, succeed, Decoder)
import Json.Encode
import Date exposing (Date)
import Date.Format


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


encodeUser : User -> Json.Encode.Value
encodeUser user =
  Json.Encode.object
    [ ( "id", Json.Encode.int user.id )
    , ( "email", Json.Encode.string user.email )
    , ( "username", Json.Encode.string user.username )
    , ( "status", Json.Encode.int user.status )
    ]


encodeToken : Token -> Json.Encode.Value
encodeToken token =
  Json.Encode.object
    [ ( "hash", Json.Encode.string token.hash )
    , ( "until", Json.Encode.string (Date.Format.formatISO8601 token.until) )
    ]

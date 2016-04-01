module Pinlist.Components.Account.Model (..) where

import Pinlist.Entities exposing (User, Token)
import Maybe exposing (..)


type alias RegisterModel =
  { username : String
  , email : String
  , password : String
  , error : Maybe String
  , writable : Bool
  }


initialRegisterModel : RegisterModel
initialRegisterModel =
  RegisterModel "" "" "" Nothing True


type alias LoginModel =
  { username : String
  , password : String
  , error : Maybe String
  , writable : Bool
  }


initialLoginModel : LoginModel
initialLoginModel =
  LoginModel "" "" Nothing True


type alias AccountModel =
  { user : Maybe User
  , token : Maybe Token
  }


initialAccountModel : AccountModel
initialAccountModel =
  AccountModel
    Nothing
    Nothing

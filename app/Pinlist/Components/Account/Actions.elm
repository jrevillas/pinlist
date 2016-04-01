module Pinlist.Components.Account.Actions (..) where

import Result exposing (Result)
import Pinlist.Entities exposing (UserAndToken)
import Http.Extra as HttpExtra


type AccountAction
  = ChangeRegisterForm RegisterFormField String
  | SubmitRegister
  | ChangeLoginForm LoginFormField String
  | SubmitLogin
  | AuthUser (Result (HttpExtra.Error String) (HttpExtra.Response UserAndToken))
  | RegisterUser (Result (HttpExtra.Error String) (HttpExtra.Response UserAndToken))
  | LoginUser (Result (HttpExtra.Error String) (HttpExtra.Response UserAndToken))


type LoginFormField
  = LoginUsernameField
  | LoginPasswordField


type RegisterFormField
  = RegisterUsernameField
  | RegisterEmailField
  | RegisterPasswordField

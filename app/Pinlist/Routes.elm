module Pinlist.Routes (..) where


type Page
  = LoginPage
  | DashboardPage
  | ListPage (Int)
  | SettingsPage
  | RegisterPage
  | Loading

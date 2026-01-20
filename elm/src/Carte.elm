port module Carte exposing
    ( Model
    , init
    , view
    , initMap
    , drawSquare
    )

import Html exposing (Html, div)
import Html.Attributes exposing (id)
import Draw_square


port initMap :
    { lat : Float, lon : Float, zoom : Int }
    -> Cmd msg


port drawSquare :
    Draw_square.Bounds
    -> Cmd msg


type alias Model =
    ()


init : ( Model, Cmd msg )
init =
    ( ()
    , initMap
        { lat = 46.603354
        , lon = 1.888334
        , zoom = 6
        }
    )


view : Model -> Html msg
view _ =
    div [ id "map" ] []

port module Carte exposing
    ( Model
    , Coord
    , Msg(..)
    , init
    , update
    , view
    , subscriptions
    , initMap
    , drawSquare
    , autoView
    , clearSquares
    )

import Html exposing (Html, div)
import Html.Attributes exposing (id)
import Json.Decode as Decode
import Draw_square



-- PORTS

port initMap :
    { lat : Float, lon : Float, zoom : Int }
    -> Cmd msg


port drawSquare :
    Draw_square.Bounds
    -> Cmd msg


port autoView : () -> Cmd msg


port clearSquares : () -> Cmd msg


port click_coord : (Decode.Value -> msg) -> Sub msg



-- TYPES

type alias Coord =
    { lat : Float
    , lon : Float
    }


type alias Model =
    { clicked : Maybe Coord }


type Msg
    = Click Decode.Value



-- INIT

init : ( Model, Cmd Msg )
init =
    ( { clicked = Nothing }
    , initMap
        { lat = 46.603354
        , lon = 1.888334
        , zoom = 6
        }
    )



-- UPDATE

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Click value ->
            case Decode.decodeValue coordDecoder value of
                Ok coord ->
                    ( { model | clicked = Just coord }, Cmd.none )

                Err _ ->
                    ( model, Cmd.none )



-- SUBSCRIPTIONS (JS → ELM)

subscriptions : Model -> Sub Msg
subscriptions _ =
    click_coord Click



-- VIEW

view : Model -> Html msg
view _ =
    div [ id "map" ] []



-- DECODER

coordDecoder : Decode.Decoder Coord
coordDecoder =
    Decode.map2 Coord
        (Decode.field "lat" Decode.float)
        (Decode.field "long" Decode.float)
